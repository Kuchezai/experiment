package http

import (
	"experiment.io/internal/controller/http/handlers"
	"experiment.io/internal/usecase"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, l *logger.Logger, segmentUC *usecase.SegmentUsecase, userUC *usecase.UserUsecase, authUC *usecase.AuthUsecase,
	authMiddleware gin.HandlerFunc) {
	router := g.Group("/api/v1")
	{
		handlers.NewSegmentHandler(router, l, segmentUC)
		handlers.NewUserHandler(router, l, userUC)
		handlers.NewAuthHandler(router, l, authUC)
	}

	static := g.Group("/history", authMiddleware)
	static.Static("", "history")
}
