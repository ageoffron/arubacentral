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

// Verbose true/false
var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "central",
	Short: "Aruba Central cli management tool",
	Long:  description,
}

var validgetCmdArgs = []string{"devices", "swarms"}

var getCmd = &cobra.Command{
	Use:       "get [devices, swarms, aps]",
	Short:     "get [devices, swarms, aps]",
	ValidArgs: validgetCmdArgs,
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
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(swarmsCmd)
	getCmd.AddCommand(devicesCmd)
	getCmd.AddCommand(apsCmd)
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
