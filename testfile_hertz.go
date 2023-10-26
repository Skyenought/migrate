package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/gin-gonic/gin"
)

const addr = ":8080"

func newHertzServer() {
	engine := server.New(server.WithHostPorts(addr))
	engine.GET("/", func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Request.Method())
		c.JSON(200, gin.H{"message": method})
		c.Next(ctx)
	})
	engine.GET("/echo", echoHandlerForHertz)
	engine.Spin()
}

func echoHandlerForHertz(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, gin.H{"message": "hello world"})
}
