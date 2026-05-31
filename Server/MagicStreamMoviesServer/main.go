package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
)

func main() {

	router := gin.Default() // Создание сервера

	// Регистрирруем обработчики
	router.GET("/hello", func(c *gin.Context) {c.String(200, "Hello, World!")})
	routes.SetupUnProtectedRoutes(router)		// Защищеённые обработчики	
	routes.SetupProtectedRoutes(router)			// Не защищённые обработчики
	
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to run server: ", err)
	}

	

}
