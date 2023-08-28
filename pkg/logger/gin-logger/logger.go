package logger

import (
	"time"

	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		fields := logrus.Fields{
			"TIME":      startTime.Format(time.RFC3339Nano),
			"LATENCY":   endTime.Sub(startTime).String(),
			"METHOD":    ctx.Request.Method,
			"URI":       ctx.Request.RequestURI,
			"STATUS":    ctx.Writer.Status(),
			"CLIENT_IP": ctx.ClientIP(),
		}

		l.WithLogrusFields(fields)

		ctx.Next()
	}
}
