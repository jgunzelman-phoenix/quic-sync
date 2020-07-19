package main

import (
	"encoding/json"
	"net/http"
	"testing"

	model "github.com/jgunzelman-phoenix/quic-sync/qs-model"
)

var urlRoot string = "https://localhost:8443/"

func TestAPI_GetVersion(t *testing.T) {
	expectedVal := "0.1.0"

	response, err := http.Get(urlRoot + "quic-sync/version")
	if err != nil {
		t.Error("Error requesting Version: " + err.Error())
		return
	}
	version := model.Version{}
	err = json.NewDecoder(response.Body).Decode(&version)
	if err != nil {
		t.Error("Error parsing response: " + err.Error())
		return
	}
	if version.Version != expectedVal {
		t.Errorf("Expected %s recieved %s", version.Version, expectedVal)
		return
	}
}
