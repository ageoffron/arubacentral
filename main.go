package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/tkanos/gonfig"
)

// Configuration Settings
type Configuration struct {
	ClientID     string
	CustomerID   string
	ClientSecret string
	Username     string
	Password     string
}

// AuthtokenStruct init csrf and session
type AuthtokenStruct struct {
	CsrfToken string
	SessionID string
}

// AuthcodeStruct auth code
type AuthcodeStruct struct {
	AuthCode string `json:"auth_code,omitempty"`
}

// TokenStruct List of tokens
type TokenStruct struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

var loglevel string
var configpath string

func init() {

	flag.StringVar(&loglevel, "loglevel", "NONE", "log level [NONE, INFO, DEBUG]")
	flag.StringVar(&configpath, "configpath", "./config/config.production.json", "path to config file ")

}
func loadconfig(configpath string) Configuration {
	configuration := Configuration{}
	var err error
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		panic(err)
	}
	err = gonfig.GetConf(configpath, &configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}

func main() {
	flag.Parse()
	var err error
	configuration := loadconfig(configpath)

	authToken, err := Gettoken(configuration.Username, configuration.Password, configuration.ClientID)
	if err != nil {
		panic(err)
	}
	authCode, err := getauthcode(configuration.CustomerID, authToken.SessionID, authToken.CsrfToken, configuration.ClientID)
	if err != nil {
		panic(err)
	}
	token, err := getaccesstoken(configuration.ClientID, configuration.ClientSecret, authCode.AuthCode, configuration.CustomerID)
	if err != nil {
		panic(err)
	}

	//log.Printf("AccessToken: %v", token.AccessToken)
	//log.Printf("RefreshToken: %v", token.RefreshToken)
	//log.Printf("TokenType: %v", token.TokenType)
	e, err := json.Marshal(token)
	fmt.Printf(string(e))

}

// Gettoken sdf
func Gettoken(username string, password string, clientID string) (AuthtokenStruct, error) {
	postdatamap := make(map[string]interface{})
	postdatamap["username"] = username
	postdatamap["password"] = password
	postdatamapjson, err := json.Marshal(postdatamap)
	jsonStr := string(postdatamapjson)
	arubaauthdata := []byte(jsonStr)
	if loglevel == "DEBUG" {
		log.Printf("Authentication user: %v clientID: %v", username, clientID)
	}
	url := fmt.Sprintf("https://apigw-prod2.central.arubanetworks.com/oauth2/authorize/central/api/login?client_id=%v", clientID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(arubaauthdata))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AuthtokenStruct{}, err
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	if resp.Status != "200 OK" {
		if loglevel == "DEBUG" {
			log.Printf("Authentication Error user: %v clientID: %v Response Status: %v", username, clientID, resp.Status)
		}
		err = errors.New(resp.Status)
		return AuthtokenStruct{}, err
	}

	csrftokencookie := resp.Header["Set-Cookie"][0]
	csrftokenrgx := regexp.MustCompile(`csrftoken=(.*?)\;`)
	csrftoken := csrftokenrgx.FindStringSubmatch(csrftokencookie)[1]

	sessionidcookie := resp.Header["Set-Cookie"][1]
	sessionidrgx := regexp.MustCompile(`session=(.*?)\;`)
	sessionid := sessionidrgx.FindStringSubmatch(sessionidcookie)[1]
	if loglevel == "DEBUG" {
		log.Printf("Authentication user: %v clientID: %v csrftoken: %v ", username, clientID, csrftoken)
		log.Printf("Authentication user: %v clientID: %v session: %v ", username, clientID, sessionid)
	}
	authtoken := AuthtokenStruct{csrftoken, sessionid}
	return authtoken, nil
}

func getauthcode(customerID string, sessionID string, csrfToken string, clientID string) (AuthcodeStruct, error) {
	var err error

	postdataMap := make(map[string]interface{})
	postdataMap["customer_id"] = customerID
	postdataMapJSON, err := json.Marshal(postdataMap)
	postdataJSONStr := string(postdataMapJSON)
	postdata := []byte(postdataJSONStr)
	if loglevel == "DEBUG" {
		log.Printf("Authentication customer id: %v", customerID)
	}
	url := fmt.Sprintf("https://apigw-prod2.central.arubanetworks.com/oauth2/authorize/central/api?client_id=%v&response_type=code&scope=all", clientID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postdata))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("session=%v", sessionID))
	req.Header.Set("X-CSRF-Token", csrfToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(resp.Status)
		return AuthcodeStruct{}, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		if loglevel == "DEBUG" {
			log.Printf("Authentication Error customerID: %v, sessionID: %v, csrfToken %v, Response Status: %v", customerID, sessionID, csrfToken, resp.Status)
		}
		err = errors.New(resp.Status)
		return AuthcodeStruct{}, err
	}

	decoder := json.NewDecoder(resp.Body)

	val := AuthcodeStruct{}
	err = decoder.Decode(&val)
	if err != nil {
		return AuthcodeStruct{}, err
	}
	if loglevel == "DEBUG" {
		log.Printf("Authcode: %v", val)
	}
	return val, nil
}

func getaccesstoken(clientID string, clientSecret string, authCode string, customerID string) (TokenStruct, error) {
	if loglevel == "DEBUG" {
		log.Printf("Authentication customer id: %v", customerID)
	}
	url := fmt.Sprintf("https://apigw-prod2.central.arubanetworks.com/oauth2/token?client_id=%v&client_secret=%v&grant_type=authorization_code&code=%v", clientID, clientSecret, authCode)

	var emptyData []byte
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(emptyData))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TokenStruct{}, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		if loglevel == "DEBUG" {
			log.Printf("Authentication Error clientID: %v, clientSecret: %v, authCode %v, Response Status: %v", clientID, clientSecret, authCode, resp.Status)
		}
		err = errors.New(resp.Status)
		return TokenStruct{}, err
	}

	decoder := json.NewDecoder(resp.Body)

	val := TokenStruct{}
	err = decoder.Decode(&val)
	if err != nil {
		return TokenStruct{}, err
	}
	if loglevel == "DEBUG" {
		log.Printf("tokens: %v", val)
	}
	return val, nil
}
