package authorisation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imryano/utils/random"
	"io/ioutil"
	"net/http"
	"os"
)

const authServiceAddress string = "http://127.0.0.1:8080"

type AccessTokenRequest struct {
	response_type string
	client_id     string
	state         string
	address       string
}

type AccessToken struct {
	client_id     string
	address       string
	access_token  string
	refresh_token string
	token_type    string
	expires       int
}

func (accessToken *AccessToken) ValidateToken() bool {
	var result bool
	url := fmt.Sprintf(authServiceAddress + "/authorise")
	jsonAT, err := json.Marshal(accessToken)
	if err == nil {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonAT))
		if err == nil {
			client := &http.Client{}
			resp, err := client.Do(req)
			if err == nil {
				err := json.NewDecoder(resp.Body).Decode(result)
				defer resp.Body.Close()
				if err == nil {
					return result
				}
			}
		}
	}
	return false
}

func ValidateTokenString(accessTokenString string) bool {
	var accessToken AccessToken
	accessTokenBytes := []byte(accessTokenString)
	err := json.NewDecoder(bytes.NewBuffer(accessTokenBytes)).Decode(accessToken)

	if err == nil {
		return accessToken.ValidateToken()
	}

	return false
}

func GetClientID() (string, error) {
	//Check clientID file first
	dat, err := ioutil.ReadFile("/data/clientid")
	if err == nil {
		return string(dat), nil
	}

	url := fmt.Sprintf(authServiceAddress + "/getclientid")
	req, err := http.NewRequest("GET", url, nil)
	if err == nil {
		result := ""
		client := &http.Client{}
		numChecks := 0
		for numChecks < 10 && result == "" {
			resp, err := client.Do(req)
			if err == nil {
				resultArr, err := ioutil.ReadAll(resp.Body)
				result = string(resultArr)
				//err := json.NewDecoder(resp.Body).Decode(result)
				if err == nil {
					err = ioutil.WriteFile("/data/clientid", resultArr, 0644)
					return result, nil
				}
			}
			numChecks += 1
			defer resp.Body.Close()
		}
	}
	return "", err
}

func ClearCachedClientID() {
	err := os.Remove("/data/clientid")
	if err == nil {
		return
	}
}

func GetAccessToken() (string, error) {
	var err error
	accessToken := &AccessToken{}
	url := fmt.Sprintf(authServiceAddress + "/gettoken")
	accessTokenRequest := &AccessTokenRequest{}

	accessTokenRequest.client_id, err = GetClientID()
	if accessTokenRequest.client_id != "" {
		accessTokenRequest.state, err = random.GenerateRandomString(50)
		if err == nil {
			accessTokenRequest.response_type = "code"
			jsonAT, err := json.Marshal(accessTokenRequest)
			if err == nil {
				req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonAT))
				if err == nil {
					client := &http.Client{}
					numChecks := 0
					for numChecks < 10 && err == nil {
						resp, err := client.Do(req)
						if err == nil {
							err := json.NewDecoder(resp.Body).Decode(accessToken)
							if err == nil {
								jsonAT, err = json.Marshal(accessToken)
								if err == nil {
									return string(jsonAT), nil
								}
							}
						}
						numChecks += 1
						defer resp.Body.Close()
					}
				}
			}
		}
	} else {
		err = errors.New("Could not get Client ID")
	}
	return "", err
}
