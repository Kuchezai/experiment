package handlers

import (
	"errors"
	"net/http"

	"experiment.io/internal/entity"
	"experiment.io/internal/usecase"
	"github.com/gin-gonic/gin"
)

type segmentRoutes struct {
	uc *usecase.SegmentUsecase
}

func NewSegmentRoutes(handler *gin.RouterGroup, uc *usecase.SegmentUsecase) {
	r := &segmentRoutes{uc}

	h := handler.Group("segments")
	{
		h.GET("/:slug", r.segmentBySlug)
		h.DELETE("/:slug", r.deleteSegment)
		h.POST("/", r.newSegment)
	}
}

type requestNewSegment struct {
	Slug string `json:"slug" binding:"required"`
}

func (r *segmentRoutes) newSegment(c *gin.Context) {
	var req requestNewSegment
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := r.uc.NewSegment(entity.Segment{
		Slug: req.Slug,
	}); err != nil {
		if errors.Is(err, entity.ErrSegmentAlreadyExist) {
			c.AbortWithError(http.StatusConflict, err)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (r *segmentRoutes) deleteSegment(c *gin.Context) {
	slug := c.Param("slug")

	if err := r.uc.DeleteSegment(slug); err != nil {
		if errors.Is(err, entity.ErrSegmentNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}


func (r *segmentRoutes) segmentBySlug(c *gin.Context) {

}