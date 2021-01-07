package centralrest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

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

// Gettoken Phase 1 auth
func Gettoken(username string, password string, clientID string, verbose bool) (AuthtokenStruct, error) {
	if verbose {
		log.Printf("Authentication Phase 1...")
	}
	postdatamap := make(map[string]string)
	postdatamap["username"] = username
	postdatamap["password"] = password
	postdatamapjson, err := json.Marshal(postdatamap)
	jsonStr := string(postdatamapjson)
	arubaauthdata := []byte(jsonStr)
	if verbose {
		log.Printf("Authentication user: %v clientID: %v", username, clientID)
	}
	url := fmt.Sprintf("%s/oauth2/authorize/central/api/login?client_id=%v", URIApi, clientID)
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
		if verbose {
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
	if verbose {
		log.Printf("Authentication user: %v clientID: %v csrftoken: %v ", username, clientID, csrftoken)
		log.Printf("Authentication user: %v clientID: %v session: %v ", username, clientID, sessionid)
	}
	authtoken := AuthtokenStruct{csrftoken, sessionid}
	return authtoken, nil
}

// Getauthcode Phase 2 auth
func Getauthcode(customerID string, sessionID string, csrfToken string, clientID string, verbose bool) (AuthcodeStruct, error) {
	if verbose {
		log.Printf("Authentication Phase 2...")
	}
	var err error

	postdataMap := make(map[string]string)
	postdataMap["customer_id"] = customerID
	postdataMapJSON, err := json.Marshal(postdataMap)
	postdataJSONStr := string(postdataMapJSON)
	postdata := []byte(postdataJSONStr)
	if verbose {
		log.Printf("Authentication customer id: %v", customerID)
	}
	url := fmt.Sprintf("%s/oauth2/authorize/central/api?client_id=%v&response_type=code&scope=all", URIApi, clientID)
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
		if verbose {
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
	if verbose {
		log.Printf("Authcode: %v", val)
	}
	return val, nil
}

// Getaccesstoken Phase 3 auth
func Getaccesstoken(clientID string, clientSecret string, authCode string, customerID string, verbose bool) (TokenStruct, error) {
	if verbose {
		log.Printf("Authentication Phase 3...")
		log.Printf("Authentication customer id: %v", customerID)
	}
	url := fmt.Sprintf("%s/oauth2/token?client_id=%v&client_secret=%v&grant_type=authorization_code&code=%v", URIApi, clientID, clientSecret, authCode)

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
		if verbose {
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
	if verbose {
		log.Printf("tokens: %v", val)
	}
	return val, nil
}
