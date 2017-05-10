package main

import (
	"encoding/json"
	"fmt"
	"github.com/imryano/utils/random"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const dbUrl string = "127.0.0.1"
const dbName string = "authDB"
const accessTokenCol string = "accessTokens"
const clientCol string = "clients"

//Structs

type AccessTokenRequest struct {
	Response_Type string
	Client_Id     string
	State         string
	Address       string
}

func (atr AccessTokenRequest) String() string {
	format := `
		response_type: 	%s
		client_id:		%s	
		state:			%s 
		address:		%s	
	`

	return fmt.Sprintf(format, atr.Response_Type, atr.Client_Id, atr.State, atr.Address)
}

type AccessToken struct {
	Id            bson.ObjectId `bson:"_id,omitempty"`
	Client_Id     string        `bson:"client_id"`
	Address       string        `bson:"address"`
	Access_Token  string        `bson:"access_token"`
	Refresh_Token string        `bson:"refresh_token"`
	Token_Type    string        `bson:"token_type"`
	Expires       int           `bson:"expires"`
}

func (at AccessToken) String() string {
	format := `
		ID:            	%d
		client_id:    	%s
		address:       	%s
		access_token:  	%s
		refresh_token: 	%s
		token_type:    	%s
		expires:		%d
	`

	return fmt.Sprintf(format, at.Id, at.Client_Id, at.Address, at.Access_Token, at.Refresh_Token, at.Token_Type, at.Expires)
}

type Client struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Client_Id string        `bson:"client_id"`
	Address   string        `bson:"address"`
}

//Functions

// Generates and returns an AccessToken string
// Will return nil if there is any kind of error
func (atr *AccessTokenRequest) GenerateAccessToken() *AccessToken {
	accessToken := &AccessToken{}

	session, err := mgo.Dial(dbUrl)
	if err == nil {
		colClients := session.DB(dbName).C(clientCol)
		numResults, err := colClients.Find(bson.M{"client_id": atr.Client_Id, "address": atr.Address}).Count()

		if numResults > 0 && err == nil {
			c := session.DB(dbName).C(accessTokenCol)
			numResults, err := c.Find(bson.M{"client_id": atr.Client_Id, "address": atr.Address}).Count()

			if numResults > 0 {
				err = c.Find(bson.M{"client_id": atr.Client_Id, "address": atr.Address}).One(accessToken)
			} else {
				accessToken = atr.GenerateAccessTokenObject()
				if accessToken != nil {
					err = c.Insert(accessToken)
				}
			}
			if err == nil {
				return accessToken
			}
		}
	}

	return &AccessToken{}
}

// GenerateAccessToken creates an AccessToken object from an AccessToken request
// Does not handle Client validation
func (atr *AccessTokenRequest) GenerateAccessTokenObject() *AccessToken {
	accessToken := &AccessToken{}
	var err error

	accessToken.Client_Id = atr.Client_Id
	accessToken.Access_Token, err = random.GenerateRandomString(50)

	if err == nil {
		accessToken.Refresh_Token, err = random.GenerateRandomString(50)
		accessToken.Expires = 600
		accessToken.Token_Type = "token"
		accessToken.Address = atr.Address
		return accessToken
	}

	return nil
}

// Validates the access key
func (accessToken *AccessToken) Validate() bool {
	session, err := mgo.Dial(dbUrl)
	if err == nil {
		col := session.DB(dbName).C(accessTokenCol)
		numResults, dbErr := col.Find(bson.M{"access_token": accessToken.Access_Token, "client_id": accessToken.Client_Id, "address": accessToken.Address}).Count()

		return (numResults > 0 && dbErr == nil)
	}
	return false
}

// GenerateClientID creates a client ID for  new service
// Generate a client ID if one doesn't exist, or return the cleint ID if one does.
// Returns blank string if there is a failure
func GetClientID(w http.ResponseWriter, r *http.Request) {
	client := &Client{}
	result := &Client{}

	session, err := mgo.Dial(dbUrl)
	if err == nil {
		col := session.DB(dbName).C(clientCol)

		client.Address = r.RemoteAddr

		err = col.Find(bson.M{"address": client.Address}).One(result)
		if err != nil {
			client.Client_Id, err = random.GenerateRandomString(50)
			if err == nil {
				err = col.Insert(client)
				if err == nil {
					fmt.Fprintln(w, client.Client_Id)
				}
			}
		} else {
			client.Client_Id = result.Client_Id
			fmt.Fprintln(w, client.Client_Id)
		}
	}
	fmt.Fprintln(w, "")
}

// GenerateAccessToken creates a key for a validated client
// Generate and return either an AccessToken or an error
func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	atr := &AccessTokenRequest{}

	err := json.NewDecoder(r.Body).Decode(&atr)
	if err == nil {
		accessToken := atr.GenerateAccessToken()
		if accessToken.Validate() {
			err = json.NewEncoder(w).Encode(accessToken)
		}
	} else {
		fmt.Println(err.Error())
	}
}

// Authorise returns true if the key matches the client id, address and access_code
// Returns false if any validation fails
func Authorise(w http.ResponseWriter, r *http.Request) {
	var accessToken AccessToken
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&accessToken)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintln(w, accessToken.Validate())
}

// Create WebServer
func main() {
	http.HandleFunc("/getclientid", GetClientID)
	http.HandleFunc("/getaccesstoken", GetAccessToken)
	http.HandleFunc("/authorise", Authorise)
	http.ListenAndServe(":8080", nil)
}
