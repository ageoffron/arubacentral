package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ageoffron/arubacentral/centralrest"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func writeTokenConfig(token centralrest.TokenStruct) {
	ViperCentralAuth.SetConfigName("arubacentraltoken")
	ViperCentralAuth.AddConfigPath("./config/")
	ViperCentralAuth.SetConfigType("json")
	ViperCentralAuth.Set("access_token", token.AccessToken)
	ViperCentralAuth.Set("refresh_token", token.RefreshToken)
	ViperCentralAuth.Set("token_type", token.TokenType)
	// keys := ViperCentralAuth.AllKeys()
	ViperCentralAuth.WriteConfig()
}

func readTokenConfig() {
	ViperCentralAuth.SetConfigName("arubacentraltoken")
	ViperCentralAuth.AddConfigPath("./config/")
	ViperCentralAuth.SetConfigType("json")
	if err := ViperCentralAuth.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found: %s", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}
}

func arubaAuth() centralrest.TokenStruct {

	var err error
	authToken, err := centralrest.Gettoken(ViperCentralConfig.GetString("username"), ViperCentralConfig.GetString("password"), ViperCentralConfig.GetString("clientID"), Verbose)
	if err != nil {
		panic(err)
	}
	authCode, err := centralrest.Getauthcode(ViperCentralConfig.GetString("CustomerID"), authToken.SessionID, authToken.CsrfToken, ViperCentralConfig.GetString("clientID"), Verbose)
	if err != nil {
		panic(err)
	}
	accesstoken, err := centralrest.Getaccesstoken(ViperCentralConfig.GetString("clientID"), ViperCentralConfig.GetString("ClientSecret"), authCode.AuthCode, ViperCentralConfig.GetString("CustomerID"), Verbose)
	if err != nil {
		panic(err)
	}

	if Verbose {
		e, err := json.Marshal(accesstoken)
		if err != nil {
			panic(err)
		}
		log.Printf("tokens: %v", string(e))
	}
	return accesstoken
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth using secrets from config file",
	Long: `
            auth using creds and secret from config file`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		token := arubaAuth()
		e, err := json.Marshal(token)
		if err != nil {
			panic(err)
		}
		fmt.Printf(string(e))
		writeTokenConfig(token)
	},
}
