package main

import (
	"net/http"
	"os"

	"github.com/example/demo-monorepo/libs/go/common"
	"github.com/gin-gonic/gin"
)

// version is injected at build time via -ldflags "-X main.version=X.Y.Z".
var version = "dev"

const serviceName = "service-b"

// newRouter builds the Gin engine. Separated from main so tests can exercise it.
func newRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": serviceName, "version": version})
	})

	// Demonstrates live-at-head use of the shared lib.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": common.Greeting(c.Query("name"))})
	})

	return r
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := newRouter().Run(":" + port); err != nil {
		panic(err)
	}
}
