package usecase

import (
	"testing"
	"time"

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

func TestRemoveUserSegments(t *testing.T) {
	r := new(mocks.UserRepo)
	uc := NewUserUsecase(r)

	testCase := []struct {
		name        string
		userID      int
		removed     []string
		repoErr     error
		expectedErr error
	}{
		{
			name: "Existent  segments, existent user",
			removed: []string{
				"ExistSegment", "ExistSegment2",
			},
			userID:      1,
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Non-existent segments, existent user",
			removed: []string{
				"UniqSegment", "UniqSegment2",
			},
			userID:      1,
			repoErr:     entity.ErrSegmentNotFound,
			expectedErr: entity.ErrSegmentNotFound,
		},
		{
			name: "Existent segments, non-existent user",
			removed: []string{
				"ExistSegment", "ExistSegment2",
			},
			userID:      0,
			repoErr:     entity.ErrSegmentNotFound,
			expectedErr: entity.ErrSegmentNotFound,
		},
		{
			name: "Non-existent segments, non-existent user",
			removed: []string{
				"UniqSegment", "UniqSegment2",
			},
			userID:      0,
			repoErr:     entity.ErrUserNotFound,
			expectedErr: entity.ErrUserNotFound,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("RemoveUserSegments", tc.userID, tc.removed).Return(tc.expectedErr)

			err := uc.RemoveUserSegments(tc.userID, tc.removed)
			require.ErrorIs(t, err, tc.expectedErr)

			mockCall.Unset()
		})
	}

}

func TestAddUserSegments(t *testing.T) {
	r := new(mocks.UserRepo)
	uc := NewUserUsecase(r)

	testCase := []struct {
		name        string
		userID      int
		added       []entity.SlugWithExpiredDate
		repoErr     error
		expectedErr error
	}{
		{
			name: "Existent user",
			added: []entity.SlugWithExpiredDate{
				{Slug: "NewSegment1", ExpiredDate: time.Now().Add(time.Hour)},
				{Slug: "NewSegment2", ExpiredDate: time.Now().Add(2 * time.Hour)},
			},
			userID:      1,
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Non-existent user",
			added: []entity.SlugWithExpiredDate{
				{Slug: "NewSegment1", ExpiredDate: time.Now().Add(time.Hour)},
				{Slug: "NewSegment2", ExpiredDate: time.Now().Add(2 * time.Hour)},
			},
			userID:      0,
			repoErr:     entity.ErrUserNotFound,
			expectedErr: entity.ErrUserNotFound,
		},
		{
			name: "Same expiration time as now",
			added: []entity.SlugWithExpiredDate{
				{Slug: "CurrentSegment", ExpiredDate: time.Now()},
			},
			userID:      1,
			repoErr:     nil,
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("AddUserSegments", tc.userID, tc.added).Return(tc.repoErr)

			err := uc.AddUserSegments(tc.userID, tc.added)
			require.ErrorIs(t, err, tc.expectedErr)

			mockCall.Unset()
		})
	}
}


func TestUserSegments(t *testing.T) {
	r := new(mocks.UserRepo)
	uc := NewUserUsecase(r)

	testCase := []struct {
		name        string
		userID      int
		repoSegments []entity.SlugWithExpiredDate
		repoErr     error
		expectedErr error
	}{
		{
			name: "Get segments for existent user",
			userID: 1,
			repoSegments: []entity.SlugWithExpiredDate{
				{Slug: "Segment1", ExpiredDate: time.Now().Add(time.Hour)},
				{Slug: "Segment2", ExpiredDate: time.Now().Add(2 * time.Hour)},
			},
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Get segments for non-existent user",
			userID: 0,
			repoSegments: nil,
			repoErr:     entity.ErrUserNotFound,
			expectedErr: entity.ErrUserNotFound,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("UserSegments", tc.userID).Return(tc.repoSegments, tc.repoErr)

			segments, err := uc.UserSegments(tc.userID)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Nil(t, segments)
			} else {
				require.NoError(t, err)
				require.NotNil(t, segments)
				require.Equal(t, tc.repoSegments, segments)
			}

			mockCall.Unset()
		})
	}
}