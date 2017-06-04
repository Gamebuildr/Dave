package scaler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Gamebuildr/Hal/pkg/config"
)

func TestHTTPScalerReturnsHTTPLoad(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		m := Response{LoadCount: 2}
		resp, _ := json.Marshal(m)
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))
	defer ts.Close()

	httpScaler := HTTPScaler{
		LoadAPIUrl: ts.URL,
		Client:     &http.Client{},
	}

	load, err := httpScaler.GetSystemLoad()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if load != 2 {
		t.Errorf("Expected %v, but got %v", 2, load)
	}
}

func TestHTTPScalerHitScalingRequestEndpoint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	httpScaler := HTTPScaler{
		AddLoadAPIUrl: ts.URL,
		Client:        &http.Client{},
	}

	resp, err := httpScaler.AddSystemLoad("")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected http status %v, but got %v", http.StatusOK, resp.Status)
	}
}

func TestHTTPScalerCanAuthenticateARouteWithCorrectToken(t *testing.T) {
	mockToken := "mockTest"
	os.Setenv(config.Auth0ClientSecret, mockToken)

	r, err := http.NewRequest(http.MethodPost, "/mock/url", nil)
	authenticateRoute(r)

	if err != nil {
		t.Fatalf(err.Error())
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		t.Errorf("Expected url to be authorized")
	}
}

func TestHTTPTokenValueReturnsSignedString(t *testing.T) {
	mockToken := "mockTest"
	os.Setenv(config.Auth0ClientSecret, mockToken)

	tokenString, err := getStringToken()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if tokenString == mockToken {
		t.Errorf("Expected %v token to be hashed, got %v", mockToken, tokenString)
	}
}
