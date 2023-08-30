package usecase

import (
	"fmt"

	"experiment.io/internal/entity"
)

type UserRepo interface {
	UserSegments(userID int) ([]entity.SlugWithExpiredDate, error)
	AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error
	RemoveUserSegments(userID int, removed []string) error
	UsersHistoryInByDate(year int, month int) ([]entity.UserSegmentsHistory, error)
	WriteHistoryToCSV(history []entity.UserSegmentsHistory, year int, month int) (string, error)
}

type UserUsecase struct {
	r UserRepo
}

func NewUserUsecase(r UserRepo) *UserUsecase {
	return &UserUsecase{r}
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
	op := "usecase.user.UserSegments"

	segments, err := uc.r.UserSegments(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return segments, nil
}

func (uc *UserUsecase) UsersHistoryInCSVByDate(year int, month int) (string, error) {
	op := "usecase.user.UsersHistoryInCSVByDate"

	history, err := uc.r.UsersHistoryInByDate(year, month)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	pathToCSV, err := uc.r.WriteHistoryToCSV(history, year, month)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return pathToCSV, nil
}
