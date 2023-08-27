package usecase

import (
	"fmt"

	"experiment.io/internal/entity"
)

type UserRepo interface {
	NewUser(user entity.User) (int, error)
}

type UserUsecase struct {
	r UserRepo
}

func NewUserUsecase(r UserRepo) *UserUsecase {
	return &UserUsecase{r}
}

func (uc *UserUsecase) NewUser(user entity.User) (int, error) {
	op := "usecase.user.New"
	id, err := uc.r.NewUser(user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
