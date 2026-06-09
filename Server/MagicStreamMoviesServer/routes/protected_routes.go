package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetupProtectedRoutes настраивает маршруты, требующие авторизации.
func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	// Подключаем мидлварь AuthMiddlware ко всем роутам, объявленным ниже.

	// Каждый запрос к этим эндпоинтам сначала пройдет проверку токена.
	router.Use(middleware.AuthMiddlware())
	
	// Регистрируем GET-маршрут для получения фильма по его IMDB ID.
	router.GET("/movie/:imdb_id", controllers.GetMovie(client))
	
	// Регистрируем POST-маршрут для добавления нового фильма в систему.
	router.POST("/addmovie", controllers.AddMovie(client))

	router.GET("recommendedmovies", controllers.GetRecomendedMovies(client))
	router.PATCH("/updatereview/:imdb_id", controllers.AdminReviewUpdate(client))
	
}