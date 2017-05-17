package main

import (
	"encoding/json"
	"github.com/imryano/utils/webservice"
	"net/http"
	"testing"
)

var clientIDs = []string{}
var accessTokens = []string{}

//GetClientID Tests
func TestPassGetClientID(t *testing.T) {
	//Number of times to repeat the test (suggested min 2).
	testRepeats := 2
	addrs := GetTestAddresses()

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
		retVal := webservice.RunWebServiceTest(req, nil, addr, getClientID)
		if retVal == "" {
			t.Errorf("GetClientID failed: Result is blank for address %s", addr)
		} else {
			//Rerun generation and make sure the result is the same
			for i := 0; i < testRepeats; i++ {
				retVal2 := webservice.RunWebServiceTest(req, nil, addr, getClientID)
				if retVal2 != retVal {
					t.Errorf("GetClientID failed: Did not return the same values twice for address: %s", addrs[i])
				}
			}
		}
		clientIDs = append(clientIDs, retVal)
	}
}

//GetAccessToken Tests
func TestPassGetAccessToken(t *testing.T) {
	addrs := GetTestAddresses()

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

		retVal := webservice.RunWebServiceTest(req, atr, addrs[i], getAccessToken)
		if retVal == "" {
			t.Errorf("GetAccessToken failed: Result is blank for address %s", atr.Address)
		} else {
			accessTokens = append(accessTokens, retVal)
		}
	}
}

func TestFailGetAccessToken(t *testing.T) {
	addrs := GetTestAddresses()

	//Create the request object. This will be modified for each test.
	req, err := http.NewRequest("GET", "/getaccesstoken", nil)
	if err != nil {
		t.Errorf("GetAccessToken failed: Could not create request: %s", err)
	}

	//Test failure conditions
	for i := 0; i < len(addrs); i++ {
		atr := &AccessTokenRequest{
			Response_Type: "None",
			Client_Id:     "THISISAFAKEANDBROKENCLIENTID",
			State:         "Active",
			Address:       addrs[i],
		}

		retVal := webservice.RunWebServiceTest(req, atr, addrs[i], getAccessToken)
		if retVal != "" {
			t.Errorf("GetAccessToken failed: Returned access token with invalid ClientID %s", atr.Address)
		}
	}
}

//Authorisation Tests
func TestPassAuthorisation(t *testing.T) {
	addrs := GetTestAddresses()
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
			retVal := (webservice.RunWebServiceTest(req, at, at.Address, authorise) == "true")

			if !retVal {
				t.Errorf("Authorise failed: AccessToken authorisation failed for address %s", at.Address)
			}
		}
	}
}

func TestFailAuthorisation(t *testing.T) {
	addrs := GetTestAddresses()
	at := &AccessToken{}

	//Create the request object. This will be modified for each test.
	req, err := http.NewRequest("GET", "/authorise", nil)
	if err != nil {
		t.Errorf("Authorise failed: Could not create request: %s", err)
	}

	//Test fail conditions
	for i := 0; i < len(addrs); i++ {
		err := json.Unmarshal([]byte(accessTokens[i]), at)

		at.Access_Token = "THISISAFAKEANDBROKENACCESSTOKEN"

		if err == nil {
			retVal := (webservice.RunWebServiceTest(req, at, at.Address, authorise) == "true")

			if retVal {
				t.Errorf("Authorise failed: AccessToken authorisation succeeded for a broken access token for address %s", at.Address)
			}
		}
	}
}
