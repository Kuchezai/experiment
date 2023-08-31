package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthorizedMiddleware(t *testing.T) {
	testCases := []struct {
		name           string
		secretKey      string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "Success",
			secretKey:      "secret-key",
			authorization:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdDEifQ.UR82ZHfp4E5oKqXNIx11NUEEwuhmDnR9l2ns4EH204w",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No authorization header",
			secretKey:      "secret-key",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Token Signature",
			secretKey:      "secret-key",
			authorization:  "Bearer invalid-jwt-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()

			router.GET("/protected", Authorized(tc.secretKey), func(c *gin.Context) {
				c.String(http.StatusOK, "Authorized")
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.authorization)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			router.ServeHTTP(c.Writer, c.Request)

			require.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
