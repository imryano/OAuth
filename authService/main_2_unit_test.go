package main

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

//CheckClientExists Tests
func TestPassCheckClientExists(t *testing.T) {
	clientList := GetTestClients()

	if !InsertClients() {
		t.Errorf("CheckClientExists failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(clientCol); success {
		for _, client := range clientList {
			retVal := checkClientExists(c, client.Address, client.Client_Id)
			if !retVal {
				t.Errorf("CheckClientExists failed: Could not find address: %s with client_id: %s", client.Address, client.Client_Id)
				return
			}

		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("CheckClientExists failed: Could not connect to database.")
	}
}

func TestFailCheckClientExists(t *testing.T) {
	clientList := GetTestClients()

	brokenClientID := "ThisIsAnotherFakeStringThatShouldBreak"

	if !InsertClients() {
		t.Errorf("CheckClientExists failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(clientCol); success {
		for _, client := range clientList {
			retVal := checkClientExists(c, client.Address, brokenClientID)
			if retVal {
				t.Errorf("CheckClientExists failed: Found address/client_id that does not exist.")
				return
			}
		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("CheckClientExists failed: Could not connect to database.")
	}
}

//GetExistingAccessToken Tests
func TestPassGetExistingAccessToken(t *testing.T) {
	accessTokenList := GetTestAccessTokens()

	if !InsertAccessTokens() {
		t.Errorf("GetExistingAccessToken failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(accessTokenCol); success {
		for _, accessToken := range accessTokenList {
			retVal, resultAT := getExistingAccessToken(c, accessToken.Address, accessToken.Client_Id)
			if !retVal {
				t.Errorf("GetExistingAccessToken failed: Could not find address: %s with client_id: %s", accessToken.Address, accessToken.Client_Id)
				return
			}

			if resultAT.Access_Token != accessToken.Access_Token || resultAT.Refresh_Token != accessToken.Refresh_Token || resultAT.Client_Id != accessToken.Client_Id || resultAT.Address != accessToken.Address {
				t.Errorf("GetExistingAccessToken failed: AccessToken returned (%s) does not match AccessToken sent (%s)", resultAT.String(), accessToken.String())
				return
			}
		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("GetExistingAccessToken failed: Could not connect to database.")
	}
}

func TestFailGetExistingAccessToken(t *testing.T) {
	accessTokenList := GetTestAccessTokens()

	brokenClientID := "ThisIsAnotherFakeStringThatShouldBreak"

	if !InsertAccessTokens() {
		t.Errorf("GetExistingAccessToken failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(accessTokenCol); success {
		for _, accessToken := range accessTokenList {
			retVal, resultAT := getExistingAccessToken(c, accessToken.Address, brokenClientID)
			if retVal {
				t.Errorf("GetExistingAccessToken failed: Found address/client_id that does not exist.")
				return
			}

			if resultAT != nil {
				t.Errorf("GetExistingAccessToken failed: Returned non-nil value on failure")
				return
			}
		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("GetExistingAccessToken failed: Could not connect to database.")
	}
}

//CreateAccessToken Tests
func TestPassCreateAccessToken(t *testing.T) {
	atrs := GetTestAccessTokenRequests()

	if success, c := GetTestCollection(accessTokenCol); success {
		for _, atr := range atrs {
			at := atr.createAccessToken(c)

			if at.Address != atr.Address || at.Client_Id != atr.Client_Id {
				t.Errorf("CreateAccessToken failed: AccessToken did not have the same address and client id for address %s", at.Address)
			} else {
				if at.Access_Token == "" {
					t.Errorf("CreateAccessToken failed: Access token did not have an AccessToken string")
				} else {
					if at.Refresh_Token == "" {
						t.Errorf("CreateAccessToken failed: Access token did not have a RefreshToken string")
					}
				}
			}
		}
	} else {
		t.Errorf("CreateAccessToken failed: Could not connect to database.")
	}
}

//ValidateAccessToken Tests
func TestPassValidateAccessToken(t *testing.T) {
	accessTokenList := GetTestAccessTokens()

	if !InsertAccessTokens() {
		t.Errorf("ValidateAccessToken failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(accessTokenCol); success {
		for _, accessToken := range accessTokenList {
			retVal := validateAccessToken(c, accessToken)
			if !retVal {
				t.Errorf("ValidateAccessToken failed: Could not find address: %s with client_id: %s and access_token: %s", accessToken.Address, accessToken.Client_Id, accessToken.Access_Token)
				return
			}
		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("ValidateAccessToken failed: Could not connect to database.")
	}
}

func TestFailValidateAccessToken(t *testing.T) {
	accessTokenList := GetTestAccessTokens()

	brokenAccessToken := "ThisIsAnotherFakeStringThatShouldBreak"

	if !InsertAccessTokens() {
		t.Errorf("ValidateAccessToken failed: Could not insert into database.")
		return
	}

	if success, c := GetTestCollection(accessTokenCol); success {
		for _, accessToken := range accessTokenList {
			accessToken.Access_Token = brokenAccessToken
			retVal := validateAccessToken(c, accessToken)
			if retVal {
				t.Errorf("ValidateAccessToken failed: Found address/client_id/access_token that does not exist.")
				return
			}
		}
		c.RemoveAll(bson.M{})
	} else {
		t.Errorf("ValidateAccessToken failed: Could not connect to database.")
	}
}
