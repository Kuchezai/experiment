package usecase

import (
	"fmt"

	"experiment.io/internal/entity"
	"experiment.io/pkg/hasher"
	"github.com/golang-jwt/jwt"
)

type AuthRepo interface {
	NewUser(user entity.User) (int, error)
	Password(username string) (string, error)
}

type AuthUsecase struct {
	r       AuthRepo
	hasher  *hasher.Hasher
	singKey string
}

func NewAuthUsecase(r AuthRepo, hasher *hasher.Hasher, singKey string) *AuthUsecase {
	return &AuthUsecase{r, hasher, singKey}
}

func (uc *AuthUsecase) Registration(user entity.User) (int, error) {
	op := "usecase.auth.Registration"

	hashedPass, err := uc.hasher.HashString(user.Password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, entity.ErrInvalidPassString)
	}

	id, err := uc.r.NewUser(entity.User{
		Name:     user.Name,
		Password: hashedPass,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (uc *AuthUsecase) Login(user entity.User) (string, error) {
	op := "usecase.auth.Login"

	encryptedPass, err := uc.r.Password(user.Name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !uc.hasher.IsHashedPassEquals(user.Password, encryptedPass) {
		return "", fmt.Errorf("%s: %w", op, entity.ErrInvalidNameOrPass)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
	})

	tokenString, err := token.SignedString([]byte(uc.singKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
