package handlers

import (
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

func newMockGinContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}
