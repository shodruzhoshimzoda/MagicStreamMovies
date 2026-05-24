package controllers

import "github.com/gin-gonic/gin"


// Возвращает список фильмов
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
			c.JSON(200, gin.H{"message":"List Of Movies"})
		
	}

}