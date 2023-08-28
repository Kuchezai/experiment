package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"experiment.io/internal/entity"
	"experiment.io/internal/usecase"
	"experiment.io/pkg/hasher"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc *usecase.UserUsecase
	l  *logger.Logger
}

func NewUserHandler(route *gin.RouterGroup, l *logger.Logger, uc *usecase.UserUsecase) {
	h := &userHandler{uc, l}

	{
		route.POST("/users", h.newUser)
		route.PATCH("/users/:user_id/segments", h.editUserSegments)
		route.GET("/users/:user_id/segments", h.userSegments)
	}
}

type requestNewUser struct {
	Name string `json:"name" binding:"required,max=100"`
	Pass string `json:"pass" binding:"required,max=50"`
}

type responseNewUser struct {
	Id int `json:"id"`
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
		Id: id,
	})
}

// added segments will be ignored after ttl expires
type requestEditUserSegments struct {
	AddSegments []struct {
		Slug string `json:"slug" binding:"required,max=100"`
		TTL  int    `json:"ttl" binding:"max=100"`
	} `json:"add_segments" binding:"max=100"`

	RemoveSegments []string `json:"remove_segments" binding:"max=100"`
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

	// check whether added and removed segments intersect
	added := make([]string, len(req.AddSegments))
	for i := range req.AddSegments {
		added[i] = req.AddSegments[i].Slug
	}
	if isIntersect(added, req.RemoveSegments) {
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

func isIntersect(a, b []string) bool {
	m := make(map[string]struct{})
	for _, el := range a {
		m[el] = struct{}{}
	}

	for _, el := range b {
		if _, ok := m[el]; ok {
			return true
		}
	}
	return false
}
