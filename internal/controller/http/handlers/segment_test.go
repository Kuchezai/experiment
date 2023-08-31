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

func TestNewSegment(t *testing.T) {
	testCases := []struct {
		name           string
		reqJSON        string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success",
			reqJSON:        `{"slug": "slug-name"}`,
			errUsecase:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Segment already exists",
			reqJSON:        `{"slug": "slug-name"}`,
			errUsecase:     entity.ErrSegmentAlreadyExist,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Unexpected usecase error",
			reqJSON:        `{"slug": "slug-name"}`,
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid request",
			reqJSON:        `{"slugggg": "slug-name"}`,
			errUsecase:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		logger := logger.New()
		mockUsecase := new(mocks.SegmentUsecase)
		mockContext := newMockGinContext()

		handler := segmentHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("NewSegment", mock.Anything).Return(tc.errUsecase)

		mockContext.Request = httptest.NewRequest("POST", "/segments", strings.NewReader(tc.reqJSON))
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.newSegment(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}

func TestDeleteSegments(t *testing.T) {
	testCase := []struct {
		name           string
		slug           string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success test",
			slug:           "slug",
			errUsecase:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-existent slug",
			slug:           "slug",
			errUsecase:     entity.ErrSegmentNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unexpected usecase error",
			slug:           "slug",
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCase {
		logger := logger.New()
		mockUsecase := new(mocks.SegmentUsecase)
		mockContext := newMockGinContext()

		handler := segmentHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("DeleteSegment", mock.Anything).Return(tc.errUsecase)

		mockContext.Params = []gin.Param{{Key: "slug", Value: tc.slug}}
		mockContext.Request = httptest.NewRequest("DELETE", "/segments/"+tc.slug+"/", nil)
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.deleteSegment(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}
func TestNewSegmentWithAutoAssign(t *testing.T) {
	testCases := []struct {
		name           string
		reqJSON        string
		errUsecase     error
		expectedStatus int
	}{
		{
			name:           "Success",
			reqJSON:        `{"slug": "slug-name", "percent": 10}`,
			errUsecase:     nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Segment already exists",
			reqJSON:        `{"slug": "slug-name", "percent": 10}`,
			errUsecase:     entity.ErrSegmentAlreadyExist,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Unexpected usecase error",
			reqJSON:        `{"slug": "slug-name", "percent": 10}`,
			errUsecase:     errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Percent is more then 100",
			reqJSON:        `{"slug": "slug-name", "percent": 101}`,
			errUsecase:     nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		logger := logger.New()
		mockUsecase := new(mocks.SegmentUsecase)
		mockContext := newMockGinContext()

		handler := segmentHandler{
			uc: mockUsecase,
			l:  logger,
		}
		mockUsecase.On("NewSegmentWithAutoAssign", mock.Anything, mock.Anything).Return([]int{1}, tc.errUsecase)

		mockContext.Request = httptest.NewRequest("POST", "/segments/auto-assign", strings.NewReader(tc.reqJSON))
		mockContext.Request.Header.Set("Accept", "application/json")

		handler.newSegmentWithAutoAssign(mockContext)
		require.Equal(t, tc.expectedStatus, mockContext.Writer.Status())
	}
}
