package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	model "github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// получаем коллекцию пользователей
var userCollection *mongo.Collection = database.OpenCollection("users")

// Фнукция для геренации хеша для паролья пользователья
func HashPassword(password string) (string, error) {
		HashPassword, err  := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return "", err
		}

		return string(HashPassword), nil
}

// Фнукция для регистрация пользователей
func RegisterUser() gin.HandlerFunc {

	return func(c *gin.Context) {

		var user model.User

		if err := c.ShouldBind(&user); err != nil {
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
			
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

		defer cancel()


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
