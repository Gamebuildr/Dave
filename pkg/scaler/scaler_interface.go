package scaler

import "net/http"

type System interface {
	GetSystemLoad() (int, error)
	AddSystemLoad(message string) (*http.Response, error)
	SetLoadURLs(loadAPIUrl string, addLoadAPIUrl string)
}
