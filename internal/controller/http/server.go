package http

import (
	"net/http"

	"experiment.io/config"
	"experiment.io/pkg/logger/gin-logger"
	"github.com/gin-gonic/gin"
)

func NewServer(g *gin.Engine, cfg config.HTTP) (*http.Server, error) {
	g.Use(gin.Recovery())
	g.Use(logger.LogsGinToJSON())

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Handler:      g,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	return srv, nil
}
