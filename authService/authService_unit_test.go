package main

import (
	"testing"
)

func GetTestAccessTokenRequests() []AccessTokenRequest {
	return []AccessTokenRequest{
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDFIRST", State: "New", Address: "123.123.123.123"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDSECOND", State: "New", Address: "69.69.69.69"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDTHIRD", State: "New", Address: "1.2.3.4"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDFOURTH", State: "New", Address: "87.65.43.21"},
	}
}

//Generate Access Token Tests
func TestPassGenerateAccessToken(t *testing.T) {
	atrs := GetTestAccessTokenRequests()

	for _, atr := range atrs {
		at := atr.GenerateAccessToken()

		if at.Address != atr.Address || at.Client_Id != atr.Client_Id {
			t.Errorf("GenerateAccessToken failed: AccessToken did not have the same address and client id for address %s", at.Address)
		} else {
			if at.Access_Token == "" {
				t.Errorf("GenerateAccessToken failed: Access token did not have an AccessToken string")
			} else {
				if at.Refresh_Token == "" {
					t.Errorf("GenerateAccessToken failed: Access token did not have a RefreshToken string")
				}
			}
		}
	}
}
