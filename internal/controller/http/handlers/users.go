package handlers

import (
	"errors"
	"net/http"

	"experiment.io/internal/entity"
	"experiment.io/internal/usecase"
	"experiment.io/pkg/hasher"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	uc *usecase.UserUsecase
}

func NewUserRoutes(handler *gin.RouterGroup, uc *usecase.UserUsecase) {
	r := &userRoutes{uc}

	{
		handler.POST("/users", r.newUser)
	}
}

type requestNewUser struct {
	Name string `json:"name" binding:"required,max=100"`
	Pass string `json:"pass" binding:"required,max=50"`
}

type responseNewUser struct {
	Id int `json:"id"`
}

func (r *userRoutes) newUser(c *gin.Context) {
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

	id, err := r.uc.NewUser(entity.User{
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
