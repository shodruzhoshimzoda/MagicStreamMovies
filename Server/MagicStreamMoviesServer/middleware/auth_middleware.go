package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shodruzxoshimzoda/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
)


// Мидлварь для авторизации пользоватл
func AuthMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Пытаемся достать токен
		// 
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":err.Error()})
			c.Abort()
			return 
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"No token provided"})
			c.Abort() 

			return 
		}
		// 2. Проверяем подпись и срок годности токена
		claims, err := utils.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid token"})
			c.Abort()
			return 	
		}

		// 3. Записываем данные в контекст для будущих хендлеров
		c.Set("userId",claims.UserId)
		c.Set("role",claims.Role)

		c.Next()
	}
}

