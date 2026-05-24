package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
)

func main() {

	router := gin.Default() // Создание сервера

	// Регистрирруем обработчики
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	router.GET("/movies", controller.GetMovies())
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to run server: ", err)

	}

}
