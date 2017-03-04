package scaler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

// AddSystemLoad will increase the systems load by 1
func (system HTTPScaler) AddSystemLoad() (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, system.AddLoadAPIUrl, nil)
	if err != nil {
		return nil, err
	}
	w, err := system.Client.Do(r)
	if err != nil {
		return w, err
	}
	return w, nil
}
