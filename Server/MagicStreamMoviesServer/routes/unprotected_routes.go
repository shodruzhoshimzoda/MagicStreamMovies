package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
)

// SetupProtectedRoutes настраивает маршруты, не требующие авторизации.

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controllers.GetMovies())
	router.POST("/register", controllers.RegisterUser())
	router.POST("/login", controllers.LogginUser())

}
