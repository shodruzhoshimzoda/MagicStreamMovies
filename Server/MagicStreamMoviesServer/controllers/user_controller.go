package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	model "github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// получаем коллекцию пользователей

// Фнукция для геренации хеша для паролья пользователья
func HashPassword(password string) (string, error) {
		HashPassword, err  := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return "", err
		}

		return string(HashPassword), nil
}

// Фнукция для регистрация пользователей
func RegisterUser(client  *mongo.Client) gin.HandlerFunc {

	return func(c *gin.Context) {

		var user model.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		validate := validator.New() // Создаём валидатора для валидации структуры по его полям

		if err := validate.Struct(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		 // Хешируем пароль
		 hashedPassword, err := HashPassword(user.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error":"Unable to hash password"})
				return 
			}
			
		ctx, cancel := context.WithTimeout(c, 100 * time.Second)

		defer cancel()
		var userCollection *mongo.Collection = database.OpenCollection("users", client)


		// Данная функция проверяет наличии уникальнсти электронной почты
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to check existing user"})

			return 
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error":"User already exist"})
			return 
		}

		user.UserId = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword		// Заменим пароль на хешированный пароль

		result ,err := userCollection.InsertOne(ctx,user)

		if err != nil {
			
			c.JSON(http.StatusInternalServerError, gin.H {
				"error":"Failed to created user",
			})
		}

		c.JSON(http.StatusCreated, result)

	}
}

// Функция для логирования пользователья
func LogginUser(client  *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLoggin model.UserLoggin			// Создаём объект на основе структуры для логгирования	

		// Декодируем значения полученной из тело запроса и присвавываем его в созданной структуре
		if err := c.ShouldBind(&userLoggin); err != nil {		
			c.JSON(http.StatusBadRequest, gin.H{"error":"invalid input data"})
			return 
		}

		ctx, cancel := context.WithTimeout(c, 100 *time.Second)
		defer cancel()

		var foundUser model.User
		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		// Находим пользователья с таким же email, если нашли то декодируем наше значения в созданной структуре
		err := userCollection.FindOne(ctx, bson.M{"email":userLoggin.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid email or password"})
			return
		}

		// Сравнения паролей
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLoggin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid email or password"})
			return
		}

		// создаём токен доступа
		 token, refreshToken, err  := utils.GenerateAllToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserId)

		 if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":"Failed to generate token",
				})
				return 
			}

		err = utils.UpdateAllTokens(foundUser.UserId, token, refreshToken, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to update http token"})
			return 
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name: "access_token",
			Value: token,
			Path: "/",
			MaxAge: 86400, // 1 hour
			Secure: true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		http.SetCookie(c.Writer, &http.Cookie{
			Name: "refresh_token",
			Value: refreshToken,
			Path: "/",
			MaxAge: 604800, // 1 week
			Secure: true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		c.JSON(http.StatusOK, model.UserResponse{
			UserID: foundUser.UserId,
			Email: foundUser.Email,
			FirstName: foundUser.FirstName,
			LastName: foundUser.LastName,
			Role: foundUser.Role,
			// Token: token,
			// RefeshToken: refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		},
		)

		
		
		
		
		
	}
}
