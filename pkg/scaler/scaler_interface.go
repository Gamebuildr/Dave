package scaler

import "net/http"

type System interface {
	GetSystemLoad() (int, error)
	AddSystemLoad() (*http.Response, error)
}

type ScalableSystem struct {
	System  System
	MaxLoad int
}
