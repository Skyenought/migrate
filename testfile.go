package main

import (
	"github.com/gin-gonic/gin"
)

const _addr = ":8080"

func newGinServer() {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		type Test struct {
			Name        string `json:"name" uri:"name"`
			ContentType string `json:"contentType" header:"Content-Type"`
		}
		var tt Test
		c.ShouldBindJSON(&tt)
		c.ShouldBindHeader(&tt)
		c.ShouldBindUri(&tt)
		c.Request.FormValue("")
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

func testMiddleware() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		c.Next()
	}, nil
}
