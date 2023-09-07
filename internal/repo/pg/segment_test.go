package pg

import (
	"testing"

	"experiment.io/internal/entity"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func TestNewSegment(t *testing.T) {
	testCases := []struct {
		name        string
		segment     entity.Segment
		poolErr     error
		expectedErr error
	}{
		{
			name: "Success",
			segment: entity.Segment{
				Slug: "slug",
			},
			poolErr:     nil,
			expectedErr: nil,
		},
		{
			name: "Segment already exists",
			segment: entity.Segment{
				Slug: "slug",
			},
			poolErr:     &pgconn.PgError{Code: DuplicatePKErrCode},
			expectedErr: entity.ErrSegmentAlreadyExist,
		},
		{
			name: "Unexpected error",
			segment: entity.Segment{
				Slug: "slug",
			},
			poolErr:     entity.ErrInternalServer,
			expectedErr: entity.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("Error creating mock database connection: %v", err)
		}
		defer mockPool.Close()
		repo := &SegmentRepository{
			db: mockPool,
		}
		mockPool.ExpectExec("INSERT INTO segments").WithArgs(tc.segment.Slug).WillReturnResult(pgxmock.NewResult("INSERT 1 0", 1)).WillReturnError(tc.poolErr)
		err = repo.NewSegment(tc.segment)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestNewSegmentWithAutoAssign(t *testing.T) {
	testCases := []struct {
		name            string
		segment         entity.Segment
		percentAssigned int
		isCreated       bool
		poolErr         error
		expectedErr     error
	}{
		{
			name: "Success",
			segment: entity.Segment{
				Slug: "slug",
			},
			percentAssigned: 50,
			isCreated:       true,
			poolErr:         nil,
			expectedErr:     nil,
		},
		{
			name: "Unexpected error",
			segment: entity.Segment{
				Slug: "slug",
			},
			isCreated:       true,
			percentAssigned: 50,
			poolErr:         entity.ErrInternalServer,
			expectedErr:     entity.ErrInternalServer,
		},
		{
			name: "Segment already exists",
			segment: entity.Segment{
				Slug: "slug",
			},
			percentAssigned: 50,
			isCreated:       false,
			poolErr:         nil,
			expectedErr:     entity.ErrSegmentAlreadyExist,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("Error creating mock database connection: %v", err)
		}
		defer mockPool.Close()
		repo := &SegmentRepository{
			db: mockPool,
		}

		mockPool.ExpectQuery("SELECT").WithArgs(tc.segment.Slug, tc.percentAssigned).WillReturnRows(
			pgxmock.NewRows([]string{"user_id", "segment_created"}).
				AddRow(1, tc.isCreated)).WillReturnError(tc.poolErr)

		_, err = repo.NewSegmentWithAutoAssign(tc.segment, tc.percentAssigned)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestDeleteSegment(t *testing.T) {
	testCases := []struct {
		name             string
		slug             string
		poolRowsAffected int
		poolErr          error
		expectedErr      error
	}{
		{
			name:             "Success",
			slug:             "slug",
			poolRowsAffected: 1,
			poolErr:          nil,
			expectedErr:      nil,
		},
		{
			name:             "Unexpected error",
			slug:             "slug",
			poolRowsAffected: 0,
			poolErr:          entity.ErrInternalServer,
			expectedErr:      entity.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("Error creating mock database connection: %v", err)
		}
		defer mockPool.Close()
		repo := &SegmentRepository{
			db: mockPool,
		}
		mockPool.ExpectExec("DELETE").WithArgs(tc.slug).WillReturnResult(
			pgxmock.NewResult("DELETE 1", int64(tc.poolRowsAffected))).WillReturnError(tc.poolErr)

		err = repo.DeleteSegment(tc.slug)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}
