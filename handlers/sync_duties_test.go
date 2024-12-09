package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSyncDuties(t *testing.T) {
	gin.SetMode(gin.TestMode)
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	testCases := []struct {
		name           string
		slot           string
		expectedStatus int
		expectJSON     bool
		expectedLength int
	}{
		{
			name:           "Valid slot",
			slot:           "10566687",
			expectedStatus: http.StatusOK,
			expectJSON:     true,
			expectedLength: 512,
		},
		{
			name:           "Missed slot (404)",
			slot:           "10564787",
			expectedStatus: http.StatusNotFound,
			expectJSON:     false,
		},
		{
			name:           "Invalid slot parameter format (400)",
			slot:           "invalid", // Simulate invalid input
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name:           "Future slot (400)",
			slot:           "4503137824400", // Simulate a slot in the future
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(ConfigMiddleware(&AppConfig{
				BaseURL: baseUrl,
			}))
			router.GET("/syncduties/:slot", GetSyncDuties)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/syncduties/%s", tc.slot), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			body := w.Body.String()
			if tc.expectJSON {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
			}
		})
	}
}
