package http

import (
	"experiment.io/internal/controller/http/handlers"
	"experiment.io/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, segmentUC *usecase.SegmentUsecase, userUC *usecase.UserUsecase) {

	router := g.Group("/api/v1")
	{
		handlers.NewSegmentRoutes(router, segmentUC)
		handlers.NewUserRoutes(router, userUC)
	}
}
