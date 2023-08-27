package usecase

import (
	"testing"

	"experiment.io/internal/entity"
	"experiment.io/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestNewSegment(t *testing.T) {
	r := new(mocks.SegmentRepo)
	uc := NewSegmentUsecase(r)

	testCases := []struct {
		name        string
		segment     entity.Segment
		repoErr     error
		expectedErr error
	}{
		{
			name: "Uniq slug",
			segment: entity.Segment{
				Slug: "slug",
			},
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Duplicate slug",
			segment: entity.Segment{
				Slug: "slug",
			},
			repoErr:     entity.ErrSegmentAlreadyExist,
			expectedErr: entity.ErrSegmentAlreadyExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("NewSegment", tc.segment).Return(tc.expectedErr)

			err := uc.NewSegment(tc.segment)
			require.ErrorIs(t, err, tc.expectedErr)

			mockCall.Unset()
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	r := new(mocks.SegmentRepo)
	uc := NewSegmentUsecase(r)

	testCases := []struct {
		name        string
		slug        string
		repoErr     error
		expectedErr error
	}{
		{
			name:        "Uniq slug",
			slug:        "slug",
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "Duplicate slug",
			slug:        "slug",
			repoErr:     entity.ErrSegmentNotFound,
			expectedErr: entity.ErrSegmentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCall := r.On("DeleteSegment", tc.slug).Return(tc.repoErr)
			err := uc.DeleteSegment(tc.slug)

			require.ErrorIs(t, err, tc.expectedErr)

			mockCall.Unset()
		})
	}
}
