package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
)

// SetupProtectedRoutes настраивает маршруты, требующие авторизации.
func SetupProtectedRoutes(router *gin.Engine) {
	// Подключаем мидлварь AuthMiddlware ко всем роутам, объявленным ниже.
	// Каждый запрос к этим эндпоинтам сначала пройдет проверку токена.
	router.Use(middleware.AuthMiddlware())
	
	// Регистрируем GET-маршрут для получения фильма по его IMDB ID.
	router.GET("/movie/:imdb_id", controllers.GetMovie())
	
	// Регистрируем POST-маршрут для добавления нового фильма в систему.
	router.POST("/addmovie", controllers.AddMovie())
	
}