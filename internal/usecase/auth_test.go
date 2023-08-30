package usecase

import (
	"testing"

	"experiment.io/internal/entity"
	"experiment.io/internal/mocks"
	"experiment.io/pkg/hasher"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	r := new(mocks.AuthRepo)
	hasher := hasher.New()
	secretKey := "secret"
	uc := NewAuthUsecase(r, hasher, secretKey)

	testCase := []struct {
		name        string
		user        entity.User
		repoVal     int
		repoErr     error
		expectedVal int
		expectedErr error
	}{
		{
			name: "Success",
			user: entity.User{
				Name:     "name",
				Password: "pass",
			},
			repoVal:     1,
			repoErr:     nil,
			expectedVal: 1,
			expectedErr: nil,
		},
		{
			name: "Already registered",
			user: entity.User{
				Name:     "name",
				Password: "pass",
			},
			repoVal:     0,
			repoErr:     entity.ErrUserAlreadyExist,
			expectedVal: 0,
			expectedErr: entity.ErrUserAlreadyExist,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("NewUser", mock.Anything).Return(tc.repoVal, tc.expectedErr)

			actual, err := uc.Registration(tc.user)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, actual, tc.expectedVal)

			mockCall.Unset()
		})
	}

}