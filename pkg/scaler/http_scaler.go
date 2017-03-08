package scaler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Gamebuildr/Hal/pkg/config"
	jwt "github.com/dgrijalva/jwt-go"
)

// HTTPScaler is a scaling system to increase remote api system scale
type HTTPScaler struct {
	LoadAPIUrl    string
	AddLoadAPIUrl string
	Client        *http.Client
}

// Response is the data that is returned from the API
type Response struct {
	LoadCount int
}

// GetSystemLoad returns to load count of the running system
func (system HTTPScaler) GetSystemLoad() (int, error) {
	r, err := http.NewRequest(http.MethodPost, system.LoadAPIUrl, nil)
	if err != nil {
		return 0, err
	}
	r.Header.Add("Content-Type", "application/json")
	if err = authenticateRoute(r); err != nil {
		return 0, err
	}

	w, err := system.Client.Do(r)
	if err != nil {
		return 0, err
	}
	resp, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return 0, err
	}
	jsonResp := Response{}
	if err := json.Unmarshal(resp, &jsonResp); err != nil {
		return 0, err
	}
	return jsonResp.LoadCount, nil
}

// AddSystemLoad will increase the systems load by one
func (system HTTPScaler) AddSystemLoad() (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, system.AddLoadAPIUrl, nil)
	if err != nil {
		return nil, err
	}
	if err = authenticateRoute(r); err != nil {
		return nil, err
	}
	w, err := system.Client.Do(r)
	if err != nil {
		return w, err
	}
	return w, nil
}

func authenticateRoute(r *http.Request) error {
	token, err := getStringToken()
	if err != nil {
		return err
	}
	bearer := "Bearer " + token
	r.Header.Add("Authorization", bearer)
	return nil
}

func getStringToken() (string, error) {
	tokenValue := os.Getenv(config.Auth0ClientSecret)
	secretKey := []byte(tokenValue)
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
