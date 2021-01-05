package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Configuration struct
type Configuration struct {
	ClientID     string
	CustomerID   string
	ClientSecret string
	Username     string
	Password     string
}

// ViperCentralConfig Configuration File
var ViperCentralConfig = viper.New()

// ViperCentralAuth Token file
var ViperCentralAuth = viper.New()

const (
	configFilename = "arubacentral.json"
	description    = `Aruba Central management tool

	Aruba Central cli to communicate with Aruba Central REST API
`
)

var configpath = []string{"./config/", "$HOME/.config/", "."}
var configuration = Configuration{}
var loglevel string

var rootCmd = &cobra.Command{
	Use:   "central",
	Short: "Aruba Central management tool",
	Long:  description,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
}

// Execute entry point
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(loadconfig)
	rootCmd.PersistentFlags().StringVar(&loglevel, "loglevel", "NONE", "log level [NONE, INFO, DEBUG]")
	rootCmd.AddCommand(authCmd)
	rootCmd.DisableSuggestions = false
	rootCmd.SuggestionsMinimumDistance = 1
}

func loadconfig() {
	ViperCentralConfig.SetConfigName("arubacentral")
	for _, cp := range configpath {
		ViperCentralConfig.AddConfigPath(cp)
	}
	ViperCentralConfig.SetEnvPrefix("ARUBACENTRAL")
	ViperCentralConfig.AutomaticEnv()
	if err := ViperCentralConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found: %s", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}
