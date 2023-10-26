package main

import "github.com/gin-gonic/gin"

const _addr = ":8080"

func newGinServer() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		method := c.Request.Method
		c.JSON(200, gin.H{"message": method})
	})
	engine.GET("/echo", echoHandler)
	engine.Run(_addr)
}

func echoHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello world"})
}
