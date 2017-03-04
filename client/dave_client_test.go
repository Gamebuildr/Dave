package client

import (
	"net/http"
	"testing"

	"github.com/Gamebuildr/Dave/pkg/scaler"
)

type MockWatcher struct {
	messagecount int
}

func (watcher MockWatcher) ReadQueueMessagesCount(url string) (int, error) {
	return watcher.messagecount, nil
}

type MockScaler struct {
	GetLoadHit   bool
	AddSystemHit bool
	loadcount    int
}

func (system *MockScaler) GetSystemLoad() (int, error) {
	system.GetLoadHit = true
	return system.loadcount, nil
}

func (system *MockScaler) AddSystemLoad() (*http.Response, error) {
	system.AddSystemHit = true

	return &http.Response{StatusCode: http.StatusOK}, nil
}

type MockLogger struct{}

func (log MockLogger) Info(message string) string {
	return ""
}

func (log MockLogger) Error(message string) string {
	return ""
}

func TestClientDoesNotRunWhenMessagesZero(t *testing.T) {
	mockClient := DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{}
	scaler := scaler.ScalableSystem{System: &scalerSystem}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 0}
	mockClient.RunClient(&scaler, "mock/url")

	if scalerSystem.GetLoadHit != false {
		t.Errorf("Expected GetSystemLoad to not be called")
	}
	if scalerSystem.AddSystemHit != false {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
}

func TestClientDoesNotRunOverMaxLoad(t *testing.T) {
	mockClient := DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{loadcount: 5}
	scaler := scaler.ScalableSystem{
		System:  &scalerSystem,
		MaxLoad: 4,
	}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 1}
	mockClient.RunClient(&scaler, "mock/url")

	if scalerSystem.GetLoadHit != true {
		t.Errorf("Expected GetSystemLoad to be called")
	}
	if scalerSystem.AddSystemHit != false {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
}

func TestClientAddsToLoadWhenLoadUnderMax(t *testing.T) {
	mockClient := DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{loadcount: 1}
	scaler := scaler.ScalableSystem{
		System:  &scalerSystem,
		MaxLoad: 4,
	}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 1}
	mockClient.RunClient(&scaler, "mock/url")

	if scalerSystem.GetLoadHit != true {
		t.Errorf("Expected GetSystemLoad to be called")
	}
	if scalerSystem.AddSystemHit != true {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
}
