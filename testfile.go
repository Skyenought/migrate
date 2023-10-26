package main

import "github.com/gin-gonic/gin"

const _addr = ":8080"

func newGinServer() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		c.ShouldBindJSON(nil)
		c.ShouldBindHeader(nil)
		c.ShouldBindUri(nil)
		method := c.Request.Method
		c.JSON(200, gin.H{"message": method})
		c.Next()
	})
	engine.GET("/echo", echoHandler)
	engine.Run(_addr)
}

func echoHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello world"})
}
