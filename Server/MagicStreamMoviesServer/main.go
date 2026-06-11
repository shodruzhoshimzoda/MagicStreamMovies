package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {

	router := gin.Default() // Создание сервера

	// Регистрирруем обработчики
	router.GET("/hello", func(c *gin.Context) { c.String(200, "Hello, World!") })

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}
	// Используем CORS для межсетегого соединения с другими портами
	config := cors.Config{}
	config.AllowAllOrigins = true
	// config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	var client *mongo.Client = database.DBConnect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to reach server: ", err.Error())
	}

	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal("Failed to dissconnect from MongoDB: ", err)
		}
	}()

	routes.SetupUnProtectedRoutes(router, client) // Защищеённые обработчики
	routes.SetupProtectedRoutes(router, client)   // Не защищённые обработчики

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to run server: ", err)
	}

}
