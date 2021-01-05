package main

import (
	"testing"

	"github.com/ageoffron/arubacentral/centralrest"
)

// TestGettokenBadCreds Test Bad credentials
func TestGettokenBadCreds(t *testing.T) {
	authcode, err := centralrest.Gettoken("anthony", "fakepassword", "fakeclienid")
	if err == nil {
		t.Errorf("Error should not be Null, got: %v", err)
	}
	if authcode.CsrfToken != "" {
		t.Errorf("CsrfToken should be empty , %v", err)
	}
	if authcode.SessionID != "" {
		t.Errorf("SessionID should be empty , %v", err)
	}

}

// authCode, err := getauthcode(customerID, authToken.SessionID, authToken.CsrfToken, clientID)
func TestGetauthcodeFailure(t *testing.T) {
	_, err := centralrest.Getauthcode("fakecustomerID", "fakeSessionID", "fakeCsrfToken", "fakeclientID")
	if err == nil {
		t.Errorf("Error should not be Null, got: %v", err)
	}
}

// token, err := getaccesstoken(clientID, clientSecret, authCode.AuthCode, customerID)
func TestGetaccesstokenFailure(t *testing.T) {
	_, err := centralrest.Getaccesstoken("fakeclientID", "fakeclientSecret", "fakeAuthcode", "fakecustomerID")
	if err == nil {
		t.Errorf("Error should not be Null, got: %v", err)
	}
}
