package usecase

import (
	"testing"

	"experiment.io/internal/entity"
	"experiment.io/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	r := new(mocks.UserRepo)
	uc := NewUserUsecase(r)

	testCase := []struct {
		name        string
		user        entity.User
		repoVal     int
		repoErr     error
		expectedVal int
		expectedErr error
	}{
		{
			name: "Uniq slug",
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
			name: "Duplicate slug",
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
			mockCall := r.On("NewUser", tc.user).Return(tc.repoVal, tc.expectedErr)

			actual, err := uc.NewUser(tc.user)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, actual, tc.expectedVal)

			mockCall.Unset()
		})
	}

}
