package tests

import (
	"github.com/HasanShahjahan/go-guest/api/auth"
	"github.com/HasanShahjahan/go-guest/api/utils"
	"net/http"
	"testing"
)

const (
	logTag = "[Token]"
)

func TestCreateToken(t *testing.T) {
	result, err := auth.CreateToken(12345)
	if err != nil {
		logging.Error(logTag, "Create JWT token is failed", err)
	}

	if result == "" {
		t.Errorf("Expected JWT. Got %s", result)
	}
}

func TestVerifyToken(t *testing.T) {

	response, err := auth.CreateToken(12345)
	if err != nil {
		logging.Error(logTag, "Create JWT token is failed", err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + response

	// Create a new request using http
	req, err := http.NewRequest("GET", "", nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	error := auth.TokenValid(req)
	if error != nil {
		logging.Error(logTag, "Create JWT token is failed", err)
	}
}
