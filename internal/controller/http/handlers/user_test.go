package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"experiment.io/internal/entity"
	"experiment.io/internal/mocks"
	"experiment.io/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEditUserSegments(t *testing.T) {
	testCase := []struct {
		name              string
		userID            string
		reqJSON           string
		errUsecaseAdded   error
		errUsecaseRemoved error
		expectedStatus    int
	}{
		{
			name:   "Success test",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusOK,
		},
		{
			name:   "Segments intersect",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment1"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusBadRequest,
		},
		{
			name:   "Non-existent added segment",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   entity.ErrSegmentNotFound,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusUnprocessableEntity,
		},

		{
			name:   "Non-existent user",
			userID: "0",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   entity.ErrUserNotFound,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusNotFound,
		},
		{
			name:   "Invalid user",
			userID: "not_int",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusNotFound,
		},
		{
			name:   "User already assigned to a segment",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   entity.ErrUserAlreadyAssigned,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusConflict,
		},
		{
			name:   "Invalid json",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segmen
					
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusBadRequest,
		},
		{
			name:   "ttl is less then 0",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": -1
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: nil,
			expectedStatus:    http.StatusBadRequest,
		},
		{
			name:   "Non-existent removed segment",
			userID: "1",
			reqJSON: `{
				"add_segments": 
				[{
					"slug": "segment1",
					"ttl": 7
				}],
				"remove_segments": ["segment2"]
				}`,
			errUsecaseAdded:   nil,
			errUsecaseRemoved: entity.ErrUserToSegmentNotFound,
			expectedStatus:    http.StatusNotFound,
		},
	}

	for _, tc := range testCase {
		logger := logger.New()
		mockUsecase := new(mocks.UserUsecase)
		mockContext := newMockGinContext()

		handler := userHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("RemoveUserSegments", mock.Anything, mock.Anything).Return(tc.errUsecaseRemoved)
		mockUsecase.On("AddUserSegments", mock.Anything, mock.Anything).Return(tc.errUsecaseAdded)

		mockContext.Params = []gin.Param{{Key: "user_id", Value: tc.userID}}
		mockContext.Request = httptest.NewRequest("PATCH", "/users/"+tc.userID+"/segments", strings.NewReader(tc.reqJSON))
		mockContext.Request.Header.Set("Content-Type", "application/json")
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.editUserSegments(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}

func TestUserSegments(t *testing.T) {
	testCase := []struct {
		name           string
		userID         string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success test",
			userID:         "1",
			errUsecase:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent user",
			userID:         "1",
			errUsecase:     entity.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unexpected error",
			userID:         "1",
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid user id",
			userID:         "1invalid",
			errUsecase:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		logger := logger.New()
		mockUsecase := new(mocks.UserUsecase)
		mockContext := newMockGinContext()

		handler := userHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("UserSegments", mock.Anything).Return([]entity.SlugWithExpiredDate{{}}, tc.errUsecase)

		mockContext.Params = []gin.Param{{Key: "user_id", Value: tc.userID}}
		mockContext.Request = httptest.NewRequest("GET", "/users/"+tc.userID+"/segments", nil)
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.userSegments(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}

func TestUsersHistoryInCSVByDate(t *testing.T) {
	testCase := []struct {
		name           string
		reqJSON        string
		errUsecase     error
		expectedStatus int
	}{
		{
			name: "Success test",
			reqJSON: `
				{
					"year": 2023,
					"month": 8
				}
				`,
			errUsecase:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Usecase error",
			reqJSON: `
				{
					"year": 2023,
					"month": 8
				}`,
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Invalid json",
			reqJSON: `
				{
					"year": 2023,
					"mon
				}`,
			errUsecase:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid month or year",
			reqJSON: `
				{
					"year": 1999,
					"month": 13
				}`,
			errUsecase:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCase {
		logger := logger.New()
		mockUsecase := new(mocks.UserUsecase)
		mockContext := newMockGinContext()

		handler := userHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("UsersHistoryInCSVByDate", mock.Anything, mock.Anything).Return("link", tc.errUsecase)

		mockContext.Request = httptest.NewRequest("POST", "/users/segments/history", strings.NewReader(tc.reqJSON))
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.createUsersHistoryInCSVByDate(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}

func newMockGinContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}
