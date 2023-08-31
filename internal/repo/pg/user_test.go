package pg

import (
	"testing"

	"experiment.io/internal/entity"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	

	testCases := []struct {
		name        string
		user        entity.User
		poolErr     error
		expectedID  int
		expectedErr error
	}{
		{
			name:        "Success",
			user:        entity.User{Name: "testuser", Password: "testpassword"},
			poolErr:     nil,
			expectedID:  123,
			expectedErr: nil,
		},
		{
			name:        "Duplicate User",
			user:        entity.User{Name: "duplicateuser", Password: "duplicatepassword"},
			poolErr:     pgxpool.UniqueViolationError{}, // Simulate duplicate PK error
			expectedID:  0,
			expectedErr: entity.ErrUserAlreadyExist,
		},
		{
			name:        "Other Pool Error",
			user:        entity.User{Name: "othererroruser", Password: "othererrorpassword"},
			poolErr:     errors.New("some error"),
			expectedID:  0,
			expectedErr: errors.New("repo.pg.user.New: some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPool.ExpectQuery("^INSERT INTO users").
				WithArgs(tc.user.Name, tc.user.Password).
				WillReturnError(tc.poolErr)

			id, err := uc.NewUser(tc.user)

			require.Equal(t, tc.expectedID, id)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}