package pg

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"experiment.io/internal/entity"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	testCases := []struct {
		name        string
		user        entity.User
		poolErr     error
		expectedID  int
		expectedErr error
	}{
		{
			name: "Success",
			user: entity.User{
				Name:     "testuser",
				Password: "hashedpassword",
			},
			poolErr:     nil,
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name: "User already exists",
			user: entity.User{
				Name:     "existinguser",
				Password: "hashedpassword",
			},
			poolErr:     &pgconn.PgError{Code: DuplicatePKErrCode},
			expectedID:  0,
			expectedErr: entity.ErrUserAlreadyExist,
		},
		{
			name: "Unexpected error",
			user: entity.User{
				Name:     "testuser",
				Password: "hashedpassword",
			},
			poolErr:     entity.ErrInternalServer,
			expectedID:  0,
			expectedErr: entity.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("Error creating mock database connection: %v", err)
		}
		defer mockPool.Close()
		repo := &UserRepository{
			db: mockPool,
		}
		mockPool.ExpectQuery("INSERT INTO users").WithArgs(tc.user.Name, tc.user.Password).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(tc.expectedID)).
			WillReturnError(tc.poolErr)

		id, err := repo.NewUser(tc.user)

		require.Equal(t, tc.expectedID, id)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestPassword(t *testing.T) {
	testCases := []struct {
		name           string
		username       string
		queryResult    string
		queryErr       error
		expectedResult string
		expectedErr    error
	}{
		{
			name:           "Success",
			username:       "testuser",
			queryResult:    "hashedpassword",
			queryErr:       nil,
			expectedResult: "hashedpassword",
			expectedErr:    nil,
		},
		{
			name:           "User not found",
			username:       "nonexistentuser",
			queryResult:    "",
			queryErr:       pgx.ErrNoRows,
			expectedResult: "",
			expectedErr:    entity.ErrInvalidNameOrPass,
		},
		{
			name:           "Unexpected error",
			username:       "testuser",
			queryResult:    "",
			queryErr:       entity.ErrInternalServer,
			expectedResult: "",
			expectedErr:    entity.ErrInternalServer,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("Error creating mock database connection: %v", err)
		}
		defer mockPool.Close()
		repo := &UserRepository{
			db: mockPool,
		}

		mockPool.ExpectQuery("SELECT").WithArgs(tc.username).
			WillReturnRows(pgxmock.NewRows([]string{"encrypted_pwd"}).AddRow(tc.queryResult)).
			WillReturnError(tc.queryErr)

		result, err := repo.Password(tc.username)

		require.Equal(t, tc.expectedResult, result)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestAddUserSegments(t *testing.T) {
	testCases := []struct {
		name          string
		userID        int
		segmentsToAdd []entity.SlugWithExpiredDate
		execErr       error
		commitErr     error
		expectedErr   error
	}{
		{
			name:   "Success",
			userID: 1,
			segmentsToAdd: []entity.SlugWithExpiredDate{
				{Slug: "segment", ExpiredDate: time.Now().Add(24 * time.Hour)},
			},
			execErr:     nil,
			commitErr:   nil,
			expectedErr: nil,
		},
		{
			name:   "Error during commit",
			userID: 1,
			segmentsToAdd: []entity.SlugWithExpiredDate{
				{Slug: "segment", ExpiredDate: time.Now().Add(24 * time.Hour)},
			},
			execErr:     nil,
			commitErr:   entity.ErrInternalServer,
			expectedErr: entity.ErrInternalServer,
		},
		{
			name:   "Error during exec",
			userID: 1,
			segmentsToAdd: []entity.SlugWithExpiredDate{
				{Slug: "segment", ExpiredDate: time.Now().Add(24 * time.Hour)},
			},
			execErr:     &pgconn.PgError{Code: DuplicatePKErrCode},
			commitErr:   nil,
			expectedErr: entity.ErrUserAlreadyAssigned,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		mockPool.ExpectBegin()

		for _, segment := range tc.segmentsToAdd {
			mockPool.ExpectExec("INSERT INTO segments_to_users").
				WithArgs(segment.Slug, tc.userID, segment.ExpiredDate).
				WillReturnResult(pgxmock.NewResult("INSERT", 1)).
				WillReturnError(tc.execErr)
		}
		mockPool.ExpectCommit().WillReturnError(tc.commitErr)

		repo := &UserRepository{
			db: mockPool,
		}

		err = repo.AddUserSegments(tc.userID, tc.segmentsToAdd)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestRemoveUserSegments(t *testing.T) {
	testCases := []struct {
		name         string
		userID       int
		rowsAffected int
		segments     []string
		execErr      error
		commitErr    error
		expectedErr  error
	}{
		{
			name:         "Success",
			userID:       1,
			rowsAffected: 2,
			segments:     []string{"segment1", "segment2"},
			execErr:      nil,
			commitErr:    nil,
			expectedErr:  nil,
		},
		{
			name:         "Error during query execution",
			userID:       1,
			rowsAffected: 1,
			segments:     []string{"segment3"},
			execErr:      entity.ErrInternalServer,
			commitErr:    nil,
			expectedErr:  entity.ErrInternalServer,
		},
		{
			name:         "Error during commit",
			userID:       1,
			rowsAffected: 1,
			segments:     []string{"segment4"},
			execErr:      nil,
			commitErr:    entity.ErrInternalServer,
			expectedErr:  entity.ErrInternalServer,
		},
		{
			name:         "Zero rows affected",
			userID:       1,
			rowsAffected: 0,
			segments:     []string{"segment4"},
			execErr:      nil,
			commitErr:    nil,
			expectedErr:  entity.ErrUserToSegmentNotFound,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		mockPool.ExpectBegin()
		for _, segment := range tc.segments {
			mockPool.ExpectExec("DELETE FROM segments_to_users").
				WithArgs(tc.userID, segment).
				WillReturnResult(pgxmock.NewResult("DELETE", int64(tc.rowsAffected))).
				WillReturnError(tc.execErr)
		}
		mockPool.ExpectCommit().WillReturnError(tc.commitErr)

		repo := &UserRepository{
			db: mockPool,
		}

		err = repo.RemoveUserSegments(tc.userID, tc.segments)
		require.ErrorIs(t, err, tc.expectedErr)
	}
}

func TestUserSegments(t *testing.T) {
	testCases := []struct {
		name        string
		userID      int
		queryResult *pgxmock.Rows
		queryErr    error
		expected    []entity.SlugWithExpiredDate
		expectedErr error
	}{
		{
			name:        "Success",
			userID:      1,
			queryResult: pgxmock.NewRows([]string{"segment_slug", "expiration_date"}).AddRow("segment1", time.Now().Add(24*time.Hour)),
			queryErr:    nil,
			expected: []entity.SlugWithExpiredDate{
				{Slug: "segment1", ExpiredDate: time.Now().Add(24 * time.Hour)},
			},
			expectedErr: nil,
		},
		{
			name:        "Error during query execution",
			userID:      2,
			queryResult: nil,
			queryErr:    entity.ErrInternalServer,
			expected:    nil,
			expectedErr: entity.ErrInternalServer,
		},
		{
			name:        "No segments found",
			userID:      3,
			queryResult: pgxmock.NewRows([]string{"segment_slug", "expiration_date"}),
			queryErr:    nil,
			expected:    nil,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		mockPool.ExpectQuery("SELECT segment_slug, expiration_date FROM segments_to_users").
			WithArgs(tc.userID).
			WillReturnRows(tc.queryResult).
			WillReturnError(tc.queryErr)

		repo := &UserRepository{
			db: mockPool,
		}

		segments, err := repo.UserSegments(tc.userID)
		require.ErrorIs(t, err, tc.expectedErr)
		require.Equal(t, tc.expected, segments)
	}
}

func TestUsersHistoryInByDate(t *testing.T) {
	testCases := []struct {
		name        string
		year        int
		month       int
		queryResult *pgxmock.Rows
		queryErr    error
		expected    []entity.UserSegmentsHistory
		expectedErr error
	}{
		{
			name:  "Success",
			year:  2023,
			month: 9,
			queryResult: pgxmock.NewRows([]string{"operation_id", "user_id", "segment_slug", "isAdded", "operation_date"}).
				AddRow(1, 1, "segment1", true, time.Now()),
			queryErr: nil,
			expected: []entity.UserSegmentsHistory{
				{
					OperationID: 1,
					UserID:      1,
					SegmentSlug: "segment1",
					IsAdded:     true,
					Date:        time.Now(),
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Error during query execution",
			year:        2023,
			month:       9,
			queryResult: nil,
			queryErr:    entity.ErrInternalServer,
			expected:    nil,
			expectedErr: entity.ErrInternalServer,
		},
		{
			name:        "No history found",
			year:        2023,
			month:       10,
			queryResult: pgxmock.NewRows([]string{"operation_id", "user_id", "segment_slug", "isAdded", "operation_date"}),
			queryErr:    nil,
			expected:    nil,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		mockPool.ExpectQuery("SELECT operation_id, user_id, segment_slug, isAdded, operation_date FROM segment_user_operations").
			WithArgs(time.Date(tc.year, time.Month(tc.month), 1, 0, 0, 0, 0, time.UTC), time.Date(tc.year, time.Month(tc.month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)).
			WillReturnRows(tc.queryResult).
			WillReturnError(tc.queryErr)

		repo := &UserRepository{
			db: mockPool,
		}

		history, err := repo.UsersHistoryInByDate(tc.year, tc.month)
		require.ErrorIs(t, err, tc.expectedErr)
		require.Equal(t, tc.expected, history)
	}
}

func TestWriteHistoryToCSV(t *testing.T) {
	testCases := []struct {
		name        string
		history     []entity.UserSegmentsHistory
		year        int
		month       int
		expectedErr error
	}{
		{
			name: "Success",
			history: []entity.UserSegmentsHistory{
				{
					OperationID: 1,
					UserID:      101,
					SegmentSlug: "segment1",
					IsAdded:     true,
					Date:        time.Now(),
				},
				{
					OperationID: 2,
					UserID:      102,
					SegmentSlug: "segment2",
					IsAdded:     false,
					Date:        time.Now().Add(24 * time.Hour),
				},
			},
			year:        time.Now().Year(),
			month:       int(time.Now().Month()),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", "test-csv")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)
			repo := &UserRepository{
				dirToStoreCSV: tmpDir,
			}
			filePath, err := repo.WriteHistoryToCSV(tc.history, tc.year, tc.month)

			require.Equal(t, tc.expectedErr, err)

			if err == nil {
				csvContent, err := ioutil.ReadFile(filePath)
				require.NoError(t, err)

				now := time.Now()
				expectedCSVContent := "operation_id,user_id,segment_slug,is_added,date\n" +
					"1,101,segment1,true," + now.Format(time.RFC3339) + "\n" +
					"2,102,segment2,false," + now.Add(24*time.Hour).Format(time.RFC3339)

				require.Equal(t, strings.TrimSpace(expectedCSVContent), strings.TrimSpace(string(csvContent)))
			}
		})
	}
}

func TestCheckUserToSegmentError(t *testing.T) {
	testCases := []struct {
		name          string
		errorToCheck  error
		expectedError error
	}{
		{
			name: "Invalid Segment FK",
			errorToCheck: &pgconn.PgError{
				Code:           NonExistentFKErrCode,
				ConstraintName: InvalidSegmentFK,
			},
			expectedError: entity.ErrSegmentNotFound,
		},
		{
			name: "Invalid User FK",
			errorToCheck: &pgconn.PgError{
				Code:           NonExistentFKErrCode,
				ConstraintName: InvalidUserFK,
			},
			expectedError: entity.ErrUserNotFound,
		},
		{
			name: "Duplicate PK Error",
			errorToCheck: &pgconn.PgError{
				Code: DuplicatePKErrCode,
			},
			expectedError: entity.ErrUserAlreadyAssigned,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &UserRepository{}

			op := "test.operation"
			err := repo.checkUserToSegmentError(op, tc.errorToCheck)

			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
