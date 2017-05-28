package scaler

import "net/http"

type System interface {
	GetSystemLoad() (int, error)
	AddSystemLoad(message string) (*http.Response, error)
}
