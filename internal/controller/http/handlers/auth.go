package handlers

import (
	"errors"
	"net/http"

	"experiment.io/internal/entity"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	uc AuthUsecase
	l  *logger.Logger
}

type AuthUsecase interface {
	Registration(user entity.User) (int, error)
	Login(user entity.User) (string, error)
}

func NewAuthHandler(route *gin.RouterGroup, l *logger.Logger, uc AuthUsecase) {
	h := &authHandler{uc, l}

	{
		route.POST("/registration", h.registration)
		route.POST("/login", h.login)
	}
}

type requestUser struct {
	Name string `json:"name" binding:"required,max=100"`
	Pass string `json:"pass" binding:"required,max=50"`
}

type responseRegistration struct {
	ID int `json:"id"`
}

func (h *authHandler) registration(c *gin.Context) {
	var req requestUser
	if err := c.BindJSON(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	id, err := h.uc.Registration(entity.User{
		Name:     req.Name,
		Password: req.Pass,
	})
	if err != nil {
		h.l.Error(err)
		switch {
		case errors.Is(err, entity.ErrUserAlreadyExist):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"msg:": entity.ErrUserAlreadyExist.Error()})
			return
		case errors.Is(err, entity.ErrInvalidPassString):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": entity.ErrInvalidPassString.Error()})
			return
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusCreated, responseRegistration{
		ID: id,
	})
}

type responseLogin struct {
	Token string `json:"token"`
}

func (h *authHandler) login(c *gin.Context) {
	var req requestUser
	if err := c.BindJSON(&req); err != nil {
		h.l.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": err.Error()})
		return
	}

	token, err := h.uc.Login(entity.User{
		Name:     req.Name,
		Password: req.Pass,
	})
	if err != nil {
		h.l.Error(err)
		switch {
		case errors.Is(err, entity.ErrInvalidNameOrPass):
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg:": entity.ErrInvalidNameOrPass.Error()})
			return
		case errors.Is(err, entity.ErrInvalidPassString):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg:": entity.ErrInvalidPassString.Error()})
			return
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusCreated, responseLogin{
		Token: token,
	})
}
