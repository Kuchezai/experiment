package http

import (
	"net/http"

	"experiment.io/config"
	"experiment.io/pkg/logger/gin-logger"
	"github.com/gin-gonic/gin"
)

func NewServer(cfg config.HTTP) (*http.Server, error) {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.LogsGinToJSON())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Handler:      r,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	return srv, nil
}
