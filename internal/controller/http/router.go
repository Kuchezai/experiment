package http

import (
	"experiment.io/internal/usecase"
	"experiment.io/internal/controller/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(g *gin.Engine, segmentUC *usecase.SegmentUsecase){
	
	router := g.Group("/v1")
	{
		handlers.NewSegmentRoutes(router, segmentUC)
	}
}