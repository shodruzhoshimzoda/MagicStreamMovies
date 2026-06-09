package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupProtectedRoutes настраивает маршруты, не требующие авторизации.

func SetupUnProtectedRoutes(router *gin.Engine, client  *mongo.Client) {

	router.GET("/movies", controllers.GetMovies(client))
	router.POST("/register", controllers.RegisterUser(client))
	router.POST("/login", controllers.LogginUser(client))

}
