package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hertz-contrib/migrate/pkg/app"
	"github.com/hertz-contrib/migrate/pkg/app/server"
	"github.com/hertz-contrib/migrate/pkg/common/utils"
)

const _addr = ":8080"

func newGinServer() {
	engine := gin.New()
	mutils.IsSrvRequestFunc(nil)
	server.H()
	engine.GET("/", func(cx *gin.Context) {
		method := cx.Request.Method
		_ = cx.Request.FormValue("test")
		cx.JSON(200, gin.H{"message": method})
		cx.Next()
	})
	engine.GET("/echo", echoHandler)
	engine.POST("/", func(cc *gin.Context) {
		value := cc.Request.FormValue("test")
		cc.JSON(200, gin.H{"message": value})
	})
	engine.Run(_addr)
}

func echoHandler(c *gin.Context) {
	app.H()
	c.JSON(200, gin.H{"message": "hello world"})
}

func testMiddleware() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		c.Next()
	}, nil
}

func testBinder(c *gin.Context) {
	type Test struct {
		Name        string `json:"name" uri:"name"`
		ContentType string `json:"contentType" header:"Content-Type"`
	}
	var tt Test
	c.ShouldBindJSON(&tt)
	c.ShouldBindHeader(&tt)
	c.ShouldBindUri(&tt)
}
