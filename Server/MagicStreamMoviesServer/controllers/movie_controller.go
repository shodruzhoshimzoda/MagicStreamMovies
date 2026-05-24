package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	models "github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second * 100)
			defer cancel()

			movieID := c.Param("imdb_id") // Берём идентификатор фильма из тело-запроса

			if movieID == "" {
				c.JSON(http.StatusNotFound, gin.H{"error":"movie ID required"})
				return 
			}

			var movie models.Movie

			// Находим фильм с таким идентификатором и декодируем его
			err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error":"movie not found"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"movie":movie})

			
	}}