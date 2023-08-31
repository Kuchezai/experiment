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

func TestRegistration(t *testing.T) {
	testCases := []struct {
		name           string
		reqJSON        string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success",
			reqJSON:        `{"name": "testuser", "pass": "testpass"}`,
			errUsecase:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid pass",
			reqJSON:        `{"name": "testuser", "pass": "testpass"}`,
			errUsecase:     entity.ErrInvalidPassString,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unexpected usecase error",
			reqJSON:        `{"name": "testuser", "pass": "testpass"}`,
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "User already exists",
			reqJSON:        `{"name": "existinguser", "pass": "testpass"}`,
			errUsecase:     entity.ErrUserAlreadyExist,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Invalid request",
			reqJSON:        `{"name": "existinguser"}`,
			errUsecase:     entity.ErrUserAlreadyExist,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.New()
			mockUsecase := new(mocks.AuthUsecase)
			mockUsecase.On("Registration", mock.Anything, mock.Anything).Return(1, tc.errUsecase)

			mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
			handler := authHandler{
				uc: mockUsecase,
				l:  logger,
			}

			mockContext.Request = httptest.NewRequest("POST", "/registration", strings.NewReader(tc.reqJSON))
			mockContext.Request.Header.Set("Content-Type", "application/json")

			handler.registration(mockContext)

			require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
		})
	}
}

func TestLogin(t *testing.T) {
	testCases := []struct {
		name           string
		reqJSON        string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success",
			reqJSON:        `{"name": "testuser", "pass": "t¦st¶pa❻ss"}`,
			errUsecase:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid pass",
			reqJSON:        `{"name": "testuser", "pass": "testpass"}`,
			errUsecase:     entity.ErrInvalidPassString,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unexpected usecase error",
			reqJSON:        `{"name": "testuser", "pass": "testpass"}`,
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "User already exists",
			reqJSON:        `{"name": "existinguser", "pass": "testpass"}`,
			errUsecase:     entity.ErrInvalidNameOrPass,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid request",
			reqJSON:        `{"name": "existinguser"}`,
			errUsecase:     entity.ErrUserAlreadyExist,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.New()
			mockUsecase := new(mocks.AuthUsecase)
			mockUsecase.On("Login", mock.Anything, mock.Anything).Return("token", tc.errUsecase)

			mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
			handler := authHandler{
				uc: mockUsecase,
				l:  logger,
			}

			mockContext.Request = httptest.NewRequest("POST", "/registration", strings.NewReader(tc.reqJSON))
			mockContext.Request.Header.Set("Content-Type", "application/json")

			handler.login(mockContext)

			require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
		})
	}
}
