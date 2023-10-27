package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.POST("/", func(c *gin.Context) {
		value := c.Request.FormValue("test")
		c.JSON(200, gin.H{"message": value})
	})
	engine.Run(":8080")
}
