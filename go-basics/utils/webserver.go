package utils

import "github.com/gin-gonic/gin"

func WebServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("localhost:5000") // listen and serve on
}