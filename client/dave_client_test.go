package client_test

import (
	"net/http"
	"testing"

	"os"

	"github.com/Gamebuildr/Dave/client"
	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/Gamebuildr/Dave/pkg/watcher"
)

type MockWatcher struct {
	MessageID    string
	Message      string
	messagecount int
}

func (mockWatcher MockWatcher) ReadNextMessage(url string) (*watcher.MessageInfo, error) {
	return &watcher.MessageInfo{
		MessageID: mockWatcher.MessageID,
		Message:   mockWatcher.Message,
	}, nil
}

func (mockWatcher MockWatcher) DeleteMessage(messageID string, url string) error {
	return nil
}

func (mockWatcher MockWatcher) Setup() {
}

type MockScaler struct {
	GetLoadHit   bool
	AddSystemHit bool
	loadcount    int
	addLoadURL   string
	getLoadURL   string
}

func (system *MockScaler) GetSystemLoad() (int, error) {
	system.GetLoadHit = true
	return system.loadcount, nil
}

func (system *MockScaler) AddSystemLoad(message string) (*http.Response, error) {
	system.AddSystemHit = true

	return &http.Response{StatusCode: http.StatusOK}, nil
}

func (system *MockScaler) SetLoadURLs(loadAPIUrl string, addLoadAPIUrl string) {
	system.getLoadURL = loadAPIUrl
	system.addLoadURL = addLoadAPIUrl
}

type MockLogger struct{}

func (log MockLogger) Info(message string) string {
	return ""
}

func (log MockLogger) Error(message string) string {
	return ""
}

func TestClientDoesNotRunWhenMessagesZero(t *testing.T) {
	mockClient := client.DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 0, MessageID: "", Message: "{}"}
	mockClient.ClientScaler = &scalerSystem
	mockClient.RunClient(client.GogetaSubsystem, "mock/url", 1)

	if scalerSystem.GetLoadHit != false {
		t.Errorf("Expected GetSystemLoad to not be called")
	}
	if scalerSystem.AddSystemHit != false {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
}

func TestClientDoesNotRunOverMaxLoad(t *testing.T) {
	mockClient := client.DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{loadcount: 5}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 1, MessageID: "1", Message: "{}"}
	mockClient.ClientScaler = &scalerSystem
	mockClient.RunClient(client.GogetaSubsystem, "mock/url", 4)

	if scalerSystem.GetLoadHit != true {
		t.Errorf("Expected GetSystemLoad to be called")
	}
	if scalerSystem.AddSystemHit != false {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
}

func TestClientAddsToLoadWhenLoadUnderMaxGogetaImage(t *testing.T) {
	mockClient := client.DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{loadcount: 1}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 1, MessageID: "1", Message: "{}"}
	mockClient.ClientScaler = &scalerSystem
	halAPI := "http://localhost/"
	os.Setenv(config.HalGogetaAPI, halAPI)
	mockClient.RunClient(client.GogetaSubsystem, "mock/url", 4)

	countAPI := halAPI + "api/container/count"
	runAPI := halAPI + "api/container/run?image=gcr.io/gamebuildr-151415/gamebuildr-gogeta"
	if scalerSystem.GetLoadHit != true {
		t.Errorf("Expected GetSystemLoad to be called")
	}
	if scalerSystem.AddSystemHit != true {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
	if scalerSystem.getLoadURL != countAPI {
		t.Errorf("Expected getLoadURL to be %v, got %v", countAPI, scalerSystem.getLoadURL)
	}
	if scalerSystem.addLoadURL != runAPI {
		t.Errorf("Expected getLoadURL to be %v, got %v", runAPI, scalerSystem.addLoadURL)
	}
}

func TestClientAddsToLoadWhenLoadUnderMaxMrRobotImage(t *testing.T) {
	mockClient := client.DaveClient{}
	mockClient.Log = MockLogger{}
	scalerSystem := MockScaler{loadcount: 1}
	mockClient.Watcher.Queue = MockWatcher{messagecount: 1, MessageID: "1", Message: `{"enginename": "godot engine", "engineversion": "2.1"}`}
	mockClient.ClientScaler = &scalerSystem
	halAPI := "http://localhost/"
	os.Setenv(config.HalMrRobotAPI, halAPI)
	mockClient.RunClient(client.MrRobotSubsystem, "mock/url", 4)

	countAPI := halAPI + "api/container/count"
	runAPI := halAPI + "api/container/run?image=gcr.io/gamebuildr-151415/mr.robot-godot-2.1.2"
	if scalerSystem.GetLoadHit != true {
		t.Errorf("Expected GetSystemLoad to be called")
	}
	if scalerSystem.AddSystemHit != true {
		t.Errorf("Expected AddSystemLoad to not be called")
	}
	if scalerSystem.getLoadURL != countAPI {
		t.Errorf("Expected getLoadURL to be %v, got %v", countAPI, scalerSystem.getLoadURL)
	}
	if scalerSystem.addLoadURL != runAPI {
		t.Errorf("Expected getLoadURL to be %v, got %v", runAPI, scalerSystem.addLoadURL)
	}
}
