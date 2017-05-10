package authorisation

import (
	"testing"
)

func TestGetClientID(t *testing.T) {
	//Test whether the service writes to the file, and can read from it
	initialValue, err := GetClientID()

	if err != nil {
		t.Errorf("GetClientID failed: Error on initial GET (%s)", err)
	}

	if initialValue == "" {
		t.Error("GetClientID fail: Returned a blank string.")
	}

	fileReadValue, err := GetClientID()

	if err != nil {
		t.Errorf("GetClientID failed: Error reading from file (%s)", err)
	}

	if initialValue != fileReadValue {
		t.Error("GetClientID fail: Did not successfully read from the file.")
	}

	ClearCachedClientID()

	secondGetValue, err := GetClientID()

	if err != nil {
		t.Errorf("GetClientID failed: Error on second server GET (%s)", err)
	}

	if initialValue != secondGetValue {
		t.Errorf("GetClientID fail: Service did not return the same value twice. FirstVal: %s, SecondVal: %s", initialValue, secondGetValue)
	}
}

func TestAccessToken(t *testing.T) {
	value, err := GetAccessToken()
	if value == "" {
		t.Errorf("GetAccessToken failed: Could not get access token (%s)", err)
	}
	if ValidateTokenString(value) {
		t.Error("ValidateTokenString failed: Validation failed.")
	}
}
