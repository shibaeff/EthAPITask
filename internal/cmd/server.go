package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ethereum-validator-api/handlers"
	"ethereum-validator-api/internal/docs"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Gin server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port := viper.GetString("server.port")
		logLevel := viper.GetString("logging.level")
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.WithError(err).Warn("Invalid log level, defaulting to info")
			level = logrus.InfoLevel
		}
		logrus.SetLevel(level)
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

		logrus.Infof("Starting server on port %s", port)

		router := gin.Default()
		router.Use(handlers.ConfigMiddleware(&handlers.AppConfig{
			BaseURL:       viper.GetString("server.ethnode"),
			EthScanAPIKey: viper.GetString("server.etherscankey"),
			Mode:          viper.GetString("server.mode"),
		}))
		docs.SwaggerInfo.BasePath = ""
		router.GET("/blockreward/:slot", handlers.GetBlockReward)
		router.GET("/syncduties/:slot", handlers.GetSyncDuties)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if err := router.Run(port); err != nil {
			logrus.Fatalf("Failed to start server: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
