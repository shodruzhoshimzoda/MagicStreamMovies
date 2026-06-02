package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	models "github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Берём коллекцию фильмов (таблицу)
var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

//  Создаём структуру валидатора

// Возвращает список фильмов
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
			return
		}

		fmt.Println(movies)
		c.JSON(http.StatusOK, movies)

	}

}

// Получение одного определённого фильма по его уникальному идентификатору
func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		movieID := c.Param("imdb_id") // Берём идентификатор фильма из тело-запроса

		if movieID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie ID required"})
			return
		}

		var movie models.Movie

		// Находим фильм с таким идентификатором и декодируем его
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"movie": movie})

	}
}

// Добавление фильмов
func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		err := c.ShouldBindJSON(&movie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
			return
		}

		// Валидируем структуру по его структураным тегам ( "validate": "required")
		if err = validate.Struct(movie); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		// Добавляем фильм в нашу структуру
		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add movie", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

// Функция для обновления комментария админа н
func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID := c.Param("imdb_id")		// получаем идентификатор фильма из контекста

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}

		// структура для коментария администратора
		var req struct {
			AdminReview string `json:"admin_review"`
		}

		// структура для ответа ИИ
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {		// присваиваем полученный комментарий по фильму в req
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}


		// получаем ответ от ИИ а также его рейтинг
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			// Добавляем поле details, которое покажет точную причину
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error getting review ranking",
				"details": err.Error(),
			})
			return
		}

		// Фильтруем по полю идентификатора фильма  и  обновим значения некоторых полей
		filter := bson.M{"imdb_id": movieID}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update) // Обновления 
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)		// Возвращаем ответ от ИИ 

	}
}

// Функция для получения оценки от ИИ а также его рейтинг
func GetReviewRanking(adminReview string) (string, int, error) {
	rankings, err := GetRankings()	// Получаем все комментарии

	if err != nil {
		return "", 0, err
	}

	sentimentDelimeted := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimeted = sentimentDelimeted + ranking.RankingName + ","
		}
	}

	sentimentDelimeted = strings.Trim(sentimentDelimeted, ",")   // чистим лишные запятые с лев и справа

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Error: .env file not found. ")
	}
	openAiKey := os.Getenv("AI_KEY")		// Получаем API ключ для подключения к AI
	if openAiKey == "" {
		return "", 0, errors.New("could not read  AI_KEY")
	}

	llm, err := openai.New(
		openai.WithToken(openAiKey),
		// Возвращаем базовый URL на OpenRouter
		openai.WithBaseURL("https://openrouter.ai/api/v1"),
		// Ставим железно рабочую бесплатную модель (без тега :free, который ломается)
		openai.WithModel("deepseek/deepseek-r1-distill-llama-70b"),
	)

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMT_TEMPLATE")  // промт для ИИ
	if base_prompt_template == "" {
		return "", 0, errors.New("could not read  BASE_PROMT_TEMPLATE")
	}

	// вместо плейсхолдера {rankings}. Цифра 1 означает, что мы заменяем только первое совпадение.
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimeted, 1)		
	respose, err := llm.Call(context.Background(), base_prompt+adminReview)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == respose {
			rankVal = ranking.RankingValue
			break
		}
	}

	return respose, rankVal, nil
}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Возвращаем все rankings из коллекции rankings
	cursor, err := rankingCollection.Find(ctx, bson.M{}) // bson.M{} - не фильтрует коллекцию и возврщает все ranking

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}




func GetRecomendedMovies() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
