package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var clientIDs = []string{}
var accessTokens = []string{}

func GetAddresses() []string {
	return []string{
		"123.123.123.123",
		"69.69.69.69",
		"1.2.3.4",
		"87.65.43.21",
	}
}

func TestGetClientID(t *testing.T) {
	//Number of times to repeat the test (suggested min 2).
	testRepeats := 2
	addrs := GetAddresses()

	//Create the request object. This will be modified for each test.
	req, err := http.NewRequest("GET", "/getclientid", nil)
	if err != nil {
		t.Errorf("GetClientID failed: Could not create request: %s", err)
	}

	//Run through the list of fake addresses and generate keys
	//Then rerun the test and compare the results with the first generated value
	for _, addr := range addrs {
		req.RemoteAddr = addr

		//Generate an initial key
		retVal := RunGetClientIDTest(req)
		if retVal == "" {
			t.Errorf("GetClientID failed: Result is blank for address %s", addr)
		} else {
			//Rerun generation and make sure the result is the same
			for i := 0; i < testRepeats; i++ {
				retVal2 := RunGetClientIDTest(req)
				if retVal2 != retVal {
					t.Errorf("GetClientID failed: Did not return the same values twice for address: %s", addrs[i])
				}
			}
		}
		clientIDs = append(clientIDs, strings.TrimSpace(retVal))
	}
}

func RunGetClientIDTest(req *http.Request) string {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetClientID)
	handler.ServeHTTP(rr, req)
	return rr.Body.String()
}

func TestGetAccessToken(t *testing.T) {
	addrs := GetAddresses()

	//Create the request object. This will be modified for each test.
	req, err := http.NewRequest("GET", "/getaccesstoken", nil)
	if err != nil {
		t.Errorf("GetAccessToken failed: Could not create request: %s", err)
	}

	//Test success conditions
	for i := 0; i < len(addrs); i++ {
		atr := &AccessTokenRequest{
			Response_Type: "None",
			Client_Id:     clientIDs[i],
			State:         "Active",
			Address:       addrs[i],
		}

		retVal := RunGetAccessTokenTest(req, atr)
		retVal = strings.TrimSpace(retVal)
		if retVal == "" {
			t.Errorf("GetAccessToken failed: Result is blank for address %s", atr.Address)
		} else {
			accessTokens = append(accessTokens, retVal)
		}
	}

	//Test failure conditions
	for i := 0; i < len(addrs); i++ {
		atr := &AccessTokenRequest{
			Response_Type: "None",
			Client_Id:     "THISISAFAKEANDBROKENCLIENTID",
			State:         "Active",
			Address:       addrs[i],
		}

		retVal := RunGetAccessTokenTest(req, atr)
		retVal = strings.TrimSpace(retVal)
		if retVal != "" {
			t.Errorf("GetAccessToken failed: Returned access token with invalid ClientID %s", atr.Address)
		}
	}
}

func RunGetAccessTokenTest(req *http.Request, atr *AccessTokenRequest) string {
	req.RemoteAddr = atr.Address
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(atr)

	if err == nil {
		req.Body = ioutil.NopCloser(b)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAccessToken)
		handler.ServeHTTP(rr, req)
		return rr.Body.String()
	}
	return ""
}

func TestAuthorise(t *testing.T) {
	addrs := GetAddresses()
	at := &AccessToken{}

	//Create the request object. This will be modified for each test.
	req, err := http.NewRequest("GET", "/authorise", nil)
	if err != nil {
		t.Errorf("Authorise failed: Could not create request: %s", err)
	}

	//Test success conditions
	for i := 0; i < len(addrs); i++ {
		err := json.Unmarshal([]byte(accessTokens[i]), at)

		if err == nil {
			retVal := RunAuthoriseTest(req, at)

			if !retVal {
				t.Errorf("Authorise failed: AccessToken authorisation failed for address %s", at.Address)
			}
		}
	}

	//Test fail conditions
	for i := 0; i < len(addrs); i++ {
		err := json.Unmarshal([]byte(accessTokens[i]), at)

		at.Access_Token = "THISISAFAKEANDBROKENACCESSTOKEN"

		if err == nil {
			retVal := RunAuthoriseTest(req, at)

			if retVal {
				t.Errorf("Authorise failed: AccessToken authorisation succeeded for a broken access token for address %s", at.Address)
			}
		}
	}
}

func RunAuthoriseTest(req *http.Request, at *AccessToken) bool {
	req.RemoteAddr = at.Address
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(at)

	if err == nil {
		req.Body = ioutil.NopCloser(b)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Authorise)
		handler.ServeHTTP(rr, req)
		return (strings.TrimSpace(rr.Body.String()) == "true")
	}
	return false
}
