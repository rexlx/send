package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_readJSON(t *testing.T) {
	var decodedJson struct {
		Foo string `json:"foo"`
	}
	// create sample json
	jason := map[string]interface{}{
		"foo": "bar",
	}
	body, _ := json.Marshal(jason)
	// create req
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Log(err)
	}

	// need  a test recorder
	rr := httptest.NewRecorder()
	defer req.Body.Close()

	err = testApp.readJSON(rr, req, &decodedJson)
	if err != nil {
		t.Error("coudlnt decode json", err)
	}

}

func Test_writeJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := jsonResponse{
		Error:   false,
		Message: "yeyeyeyeye",
	}
	headers := make(http.Header)
	headers.Add("YO", "MAN")
	err := testApp.writeJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Errorf("test failed %v", err)
	}
	testApp.environment = "production"
	err = testApp.writeJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Errorf("test failed (env specific) %v", err)
	}
	testApp.environment = "dev"
}

func Test_errorJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	err := testApp.errorJSON(rr, errors.New("your fink isn't in the bluq"))
	if err != nil {
		t.Error(err)
	}
}
