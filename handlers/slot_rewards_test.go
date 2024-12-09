package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func loadConfig() (string, error) {
	viper.SetConfigFile("../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}
	return viper.GetString("server.ethnode"), nil
}

func TestGetSlotRewardLight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	baseUrl, err := loadConfig()
	if err != nil {
		t.Skip()
	}
	testCases := []struct {
		name           string
		slot           string
		expectedStatus int
		expectedBody   string // For partial matching or error messages
		expectJSON     bool
		expectedJSON   map[string]interface{}
	}{
		{
			name:           "Valid slot",
			slot:           "4700013",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectJSON:     true,
			expectedJSON: map[string]interface{}{
				"status": false,
				"reward": 45031378244,
			},
		},
		{
			name:           "Non-existent slot (404)",
			slot:           "-1",
			expectedStatus: http.StatusNotFound,
			expectJSON:     false,
		},
		{
			name:           "Invalid slot format (400)",
			slot:           "invalid", // Simulate invalid input
			expectedStatus: http.StatusBadRequest,
			expectJSON:     false,
		},
		{
			name:           "Future slot (400)",
			slot:           "4503137824400", // Simulate a slot in the future
			expectedStatus: http.StatusBadRequest,
			expectedBody:   constSlotInFuture,
			expectJSON:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(ConfigMiddleware(&AppConfig{
				BaseURL: baseUrl,
			}))
			router.GET("/slotreward/:slot", GetBlockReward)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/slotreward/%s", tc.slot), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			body := w.Body.String()
			if tc.expectJSON {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(body), &response)
				require.NoError(t, err)
			} else {
				require.Contains(t, body, tc.expectedBody)
			}
		})
	}

}
