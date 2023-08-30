package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"experiment.io/internal/entity"
	"experiment.io/pkg/hasher"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc UserUsecase
	l  *logger.Logger
}

type UserUsecase interface {
	NewUser(user entity.User) (int, error)
	UserSegments(userID int) ([]entity.SlugWithExpiredDate, error)
	AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error
	RemoveUserSegments(userID int, removed []string) error
	UsersHistoryInCSVByDate(year int, month int) (string, error)
}

func NewUserHandler(route *gin.RouterGroup, l *logger.Logger, uc UserUsecase) {
	h := &userHandler{uc, l}

	{
		route.POST("/users", h.newUser)
		route.GET("/users/segments/history", h.usersHistoryInCSVByDate)
		route.PATCH("/users/:user_id/segments", h.editUserSegments)
		route.GET("/users/:user_id/segments", h.userSegments)
	}
}

type requestNewUser struct {
	Name string `json:"name" binding:"required,max=100"`
	Pass string `json:"pass" binding:"required,max=50"`
}

type responseNewUser struct {
	ID int `json:"id"`
}

func (h *userHandler) newUser(c *gin.Context) {
	var req requestNewUser
	if err := c.BindJSON(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	hashedPass, err := hasher.HashString(req.Pass)
	if err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	id, err := h.uc.NewUser(entity.User{
		Name:     req.Name,
		Password: hashedPass,
	})
	if err != nil {
		h.l.Error(err)
		if errors.Is(err, entity.ErrUserAlreadyExist) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg:": entity.ErrUserAlreadyExist.Error()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responseNewUser{
		ID: id,
	})
}

// added segments will be ignored after ttl expires
type requestEditUserSegments struct {
	AddSegments    []AddSegments `json:"add_segments" binding:"max=100"`
	RemoveSegments []string      `json:"remove_segments" binding:"max=100"`
}

type AddSegments struct {
	Slug string `json:"slug" binding:"required,max=100"`
	TTL  int    `json:"ttl" binding:"min=0,max=100"`
}

func (h *userHandler) editUserSegments(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		h.l.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var req requestEditUserSegments
	if err := c.BindJSON(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	validator := NewValidator()
	if !validator.checkAddedSegmentsIsValid(req.AddSegments) {
		h.l.Error(entity.ErrInvalidAddedSegment)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": entity.ErrInvalidAddedSegment.Error()})
		return
	}
	// check whether added and removed segments intersect
	added := make([]string, len(req.AddSegments))
	for i := range req.AddSegments {
		added[i] = req.AddSegments[i].Slug
	}
	if validator.IsIntersect(added, req.RemoveSegments) {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": entity.ErrSegmentsIntersect.Error()})
		return
	}

	if len(req.RemoveSegments) > 0 {
		if err := h.uc.RemoveUserSegments(id, req.RemoveSegments); err != nil {
			h.l.Error(err)
			if errors.Is(err, entity.ErrUserToSegmentNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg:": entity.ErrUserToSegmentNotFound.Error()})
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	addedSlugWithTTL := make([]entity.SlugWithExpiredDate, len(req.AddSegments))
	for i, reqSeg := range req.AddSegments {
		addedSlugWithTTL[i] = entity.SlugWithExpiredDate{
			Slug:        reqSeg.Slug,
			ExpiredDate: time.Now().Add(time.Duration(reqSeg.TTL) * 24 * time.Hour),
		}
	}

	if len(added) > 0 {
		if err := h.uc.AddUserSegments(id, addedSlugWithTTL); err != nil {
			h.l.Error(err)
			status := http.StatusInternalServerError
			respErr := entity.ErrInternalServer
			switch {
			case errors.Is(err, entity.ErrUserNotFound):
				status = http.StatusNotFound
				respErr = entity.ErrUserNotFound
			case errors.Is(err, entity.ErrSegmentNotFound):
				status = http.StatusUnprocessableEntity
				respErr = entity.ErrSegmentNotFound
			case errors.Is(err, entity.ErrUserAlreadyAssigned):
				status = http.StatusConflict
				respErr = entity.ErrUserAlreadyAssigned
			}
			c.AbortWithStatusJSON(status, gin.H{"msg:": respErr.Error()})
			return
		}
	}

	c.Status(http.StatusOK)

}

type responseUserSegments struct {
	Slug        string    `json:"slug"`
	ExpiredDate time.Time `json:"expired_date"`
}

func (h *userHandler) userSegments(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		h.l.Error(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	segments, err := h.uc.UserSegments(id)
	if err != nil {
		h.l.Error(err)
		if errors.Is(err, entity.ErrUserNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp := make([]responseUserSegments, len(segments))
	for i, seg := range segments {
		resp[i].Slug = seg.Slug
		resp[i].ExpiredDate = seg.ExpiredDate
	}

	c.JSON(http.StatusOK, resp)

}

type requestHistoryInCSVByDate struct {
	Year  int `form:"year" binding:"required,numeric,min=2007,max=2100"`
	Month int `form:"month" binding:"required,numeric,min=1,max=12"`
}

type responseHistoryInCSVByDate struct {
	Link string `json:"link"`
}

func (h *userHandler) usersHistoryInCSVByDate(c *gin.Context) {

	var req requestHistoryInCSVByDate
	if err := c.ShouldBindQuery(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	path, err := h.uc.UsersHistoryInCSVByDate(req.Year, req.Month)
	if err != nil {
		h.l.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responseHistoryInCSVByDate{
		Link: path,
	})

}
