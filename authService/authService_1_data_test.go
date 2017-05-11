package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const dbTest = "testAuthDB"

func GetTestAddresses() []string {
	return []string{
		"123.123.123.123",
		"69.69.69.69",
		"1.2.3.4",
		"87.65.43.21",
	}
}

func GetTestClients() []Client {
	return []Client{
		Client{Client_Id: "FAKECLIENTIDFIRST", Address: "123.123.123.123"},
		Client{Client_Id: "FAKECLIENTIDSECOND", Address: "69.69.69.69"},
		Client{Client_Id: "FAKECLIENTIDTHIRD", Address: "1.2.3.4"},
		Client{Client_Id: "FAKECLIENTIDFOURTH", Address: "87.65.43.21"},
	}
}

func GetTestAccessTokenRequests() []AccessTokenRequest {
	return []AccessTokenRequest{
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDFIRST", State: "New", Address: "123.123.123.123"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDSECOND", State: "New", Address: "69.69.69.69"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDTHIRD", State: "New", Address: "1.2.3.4"},
		AccessTokenRequest{Response_Type: "Test", Client_Id: "FAKECLIENTIDFOURTH", State: "New", Address: "87.65.43.21"},
	}
}

func GetTestAccessTokens() []AccessToken {
	return []AccessToken{
		AccessToken{Client_Id: "FAKECLIENTIDFIRST", Address: "123.123.123.123", Access_Token: "123123123123", Refresh_Token: "321321321321", Expires: 600, Token_Type: "token"},
		AccessToken{Client_Id: "FAKECLIENTIDSECOND", Address: "69.69.69.69", Access_Token: "69696969", Refresh_Token: "96969696", Expires: 600, Token_Type: "token"},
		AccessToken{Client_Id: "FAKECLIENTIDTHIRD", Address: "1.2.3.4", Access_Token: "1234", Refresh_Token: "4321", Expires: 600, Token_Type: "token"},
		AccessToken{Client_Id: "FAKECLIENTIDFOURTH", Address: "87.65.43.21", Access_Token: "87654321", Refresh_Token: "12345678", Expires: 600, Token_Type: "token"},
	}
}

//Data Insertion
func InsertClients() bool {
	clientList := GetTestClients()

	if success, c := GetTestCollection(clientCol); success {
		c.RemoveAll(bson.M{})

		for _, client := range clientList {
			err := c.Insert(client)
			if err != nil {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func InsertAccessTokens() bool {
	accessTokenList := GetTestAccessTokens()

	if success, c := GetTestCollection(accessTokenCol); success {
		c.RemoveAll(bson.M{})

		for _, accessToken := range accessTokenList {
			err := c.Insert(accessToken)
			if err != nil {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func GetTestCollection(colName string) (bool, *mgo.Collection) {
	session, err := mgo.Dial(dbUrl)
	if err == nil {
		c := session.DB(dbTest).C(colName)
		return true, c
	}
	return false, nil
}
