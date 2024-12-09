package handlers

import (
	"github.com/gin-gonic/gin"
)

type AppConfig struct {
	BaseURL       string `json:"base_url"`
	EthScanAPIKey string `json:"eth_scan_api_key"`
	Mode          string `json:"mode"`
}

func ConfigMiddleware(cfg *AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}
