package usecase

import (
	"fmt"

	"experiment.io/internal/entity"
)

type UserRepo interface {
	NewUser(user entity.User) (int, error)
	UserSegments(userID int) ([]entity.SlugWithExpiredDate, error)
	AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error
	RemoveUserSegments(userID int, removed []string) error
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

func (uc *UserUsecase) RemoveUserSegments(userID int, removed []string) error {
	op := "usecase.user.RemoveUserSegments"

	if err := uc.r.RemoveUserSegments(userID, removed); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *UserUsecase) AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error {
	op := "usecase.user.AddUserSegments"

	if err := uc.r.AddUserSegments(userID, added); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (uc *UserUsecase) UserSegments(userID int) ([]entity.SlugWithExpiredDate, error) {
	op := "usecase.user.AddUserSegments"

	segments, err := uc.r.UserSegments(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return segments, nil
}
