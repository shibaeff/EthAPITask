package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd is the base command without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mewatcher",
	Short: "mewatcher is an Ethereum validator service",
	Long:  `A simple CLI application that runs a Gin server to provide Ethereum validator-related endpoints.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("config", "", "config file (default is ./config.yaml)")
	//nolint:errcheck // That's expected
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.SetEnvPrefix("mewatcher")
	viper.AutomaticEnv()
	viper.SetDefault("server.port", ":8000")
	viper.SetDefault("logging.level", "info")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
