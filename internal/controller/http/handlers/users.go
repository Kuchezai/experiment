package handlers

import (
	"errors"
	"fmt"
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
}

func NewUserHandler(route *gin.RouterGroup, l *logger.Logger, uc *usecase.UserUsecase) {
	h := &userHandler{uc}

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
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hashedPass, err := hasher.HashString(req.Pass)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := h.uc.NewUser(entity.User{
		Name:     req.Name,
		Password: hashedPass,
	})
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExist) {
			c.AbortWithError(http.StatusConflict, err)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, responseNewUser{
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
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var req requestEditUserSegments
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// check whether added and removed segments intersect
	added := make([]string, len(req.AddSegments))
	for i := range req.AddSegments {
		added[i] = req.AddSegments[i].Slug
	}
	if isIntersect(added, req.RemoveSegments) {
		c.AbortWithError(http.StatusBadRequest, entity.ErrSegmentsIntersect)
		return
	}

	if len(req.RemoveSegments) > 0 {
		if err := h.uc.RemoveUserSegments(id, req.RemoveSegments); err != nil {
			if errors.Is(err, entity.ErrUserToSegmentNotFound) {
				c.AbortWithError(http.StatusNotFound, err)
				return
			}
			c.AbortWithError(http.StatusInternalServerError, err)
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
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, entity.ErrUserNotFound):
				status = http.StatusNotFound
			case errors.Is(err, entity.ErrSegmentNotFound):
				status = http.StatusUnprocessableEntity
			case errors.Is(err, entity.ErrUserAlreadyAssigned):
				status = http.StatusConflict
			}
			c.AbortWithError(status, err)
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
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	segments, err := h.uc.UserSegments(id)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, entity.ErrUserNotFound) {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
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
