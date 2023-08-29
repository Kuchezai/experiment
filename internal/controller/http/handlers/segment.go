package handlers

import (
	"errors"
	"net/http"

	"experiment.io/internal/entity"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
)

type segmentHandler struct {
	uc SegmentUsecase
	l  *logger.Logger
}
type SegmentUsecase interface {
	NewSegment(seg entity.Segment) error
	DeleteSegment(slug string) error
}

func NewSegmentHandler(route *gin.RouterGroup, l *logger.Logger, uc SegmentUsecase) {
	h := &segmentHandler{uc, l}

	{
		route.DELETE("/segments/:slug", h.deleteSegment)
		route.POST("/segments", h.newSegment)
	}
}

type requestNewSegment struct {
	Slug string `json:"slug" binding:"required,max=100"`
}

func (h *segmentHandler) newSegment(c *gin.Context) {
	var req requestNewSegment
	if err := c.BindJSON(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	if err := h.uc.NewSegment(entity.Segment{
		Slug: req.Slug,
	}); err != nil {
		h.l.Error(err)
		if errors.Is(err, entity.ErrSegmentAlreadyExist) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg:": entity.ErrSegmentAlreadyExist.Error()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *segmentHandler) deleteSegment(c *gin.Context) {
	slug := c.Param("slug")

	if err := h.uc.DeleteSegment(slug); err != nil {
		h.l.Error(err)
		if errors.Is(err, entity.ErrSegmentNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *segmentHandler) segmentBySlug(c *gin.Context) {

}
