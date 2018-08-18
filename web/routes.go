package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())

	g.Use(mw...)

	// 404
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "NotFound API route")
	})

	g.GET("/", Index)
	g.GET("/users/:id", User)

	admin := g.Group("/admin")
	{
		admin.GET("", nil)
	}

	return g
}

func checkErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}
