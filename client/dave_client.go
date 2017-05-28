package client

import (
	"encoding/json"
	"net/http"
	"os"

	"fmt"

	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/Gamebuildr/Dave/pkg/gamebuildr_containers"
	"github.com/Gamebuildr/Dave/pkg/scaler"
	"github.com/Gamebuildr/Dave/pkg/watcher"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/papertrail"
)

// GogetaSubsystem is a representation for identifying a ScalableSystem as gogeta
const GogetaSubsystem = "GOGETA"

// MrRobotSubsystem is a representation for identifying a ScalableSystem as mr robot
const MrRobotSubsystem = "MRROBOT"

// DaveClient scales microservices remotely
type DaveClient struct {
	Watcher watcher.QueueMonitor
	Log     logger.Log
	DevMode bool
}

type basicEngineMessage struct {
	EngineName    string `json:"enginename"`
	EngineVersion string `json:"engineversion"`
}

const logFileName string = "dave_client_"

// Create a new DaveClient
func (client *DaveClient) Create() {
	logs := logger.SystemLogger{}
	if client.DevMode {
		fileLogger := logger.FileLogSave{
			LogFileName: logFileName,
			LogFileDir:  os.Getenv(config.LogEndpoint),
		}
		logs.LogSave = fileLogger
	} else {
		papertrailLog := papertrail.PapertrailLogSave{
			App: "Dave",
			URL: os.Getenv(config.LogEndpoint),
		}
		logs.LogSave = papertrailLog
	}

	// setup client
	clientWatcher := watcher.QueueWatcher{}
	clientWatcher.Setup()
	client.Watcher.Queue = clientWatcher
	client.Log = logs

	client.Log.Info("Dave Setup and Running")
}

// RunClient will read queue messages and send a response to an api endpoint to
// start up systems if the load is less than the max load
func (client *DaveClient) RunClient(subSystem string, queueURL string, maxLoad int) {
	message, err := client.Watcher.Queue.ReadNextMessage(queueURL)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	if message.MessageID == "" {
		client.Log.Info(fmt.Sprintf("Found no messages for %v", subSystem))
		return
	}
	client.Log.Info(fmt.Sprintf("[%v] Found message for %v, %v", message.MessageID, subSystem, message.Message))
	if err := client.Watcher.Queue.DeleteMessage(message.MessageID, queueURL); err != nil {
		client.Log.Error(err.Error())
	}
	rawMessageData := []byte(message.Message)
	engineData := new(basicEngineMessage)
	if err := json.Unmarshal(rawMessageData, &engineData); err != nil {
		client.Log.Error(err.Error())
		return
	}

	clientScaler, err := client.getScaler(subSystem, engineData.EngineName, engineData.EngineVersion)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	load, err := clientScaler.GetSystemLoad()
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	client.Log.Info(fmt.Sprintf("[%v] current load: %v", message.MessageID, load))
	if load > maxLoad {
		client.Log.Info("system at max load")
		return
	}

	client.Log.Info(fmt.Sprintf("[%v] Adding to system load", message.MessageID))
	resp, err := clientScaler.AddSystemLoad(message.Message)
	if err != nil {
		respinfo := fmt.Sprintf("Hal Response: %v", resp)
		client.Log.Info(respinfo)
		client.Log.Error(err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Adding container load failed with response %v", resp.Status)
		client.Log.Error(err.Error())
	}
}

func (client *DaveClient) getScaler(subSystem string, engineName string, engineVersion string) (scaler.System, error) {
	var container string
	var halBaseURL string
	if subSystem == GogetaSubsystem {
		container = "?image=gcr.io/gamebuildr-151415/gamebuildr-gogeta"
		halBaseURL = os.Getenv(config.HalGogetaAPI)
	} else if subSystem == MrRobotSubsystem {
		containers := gamebuildrContainers.GamebuildrContainers{}
		imageName, err := containers.GetContainerImageName(engineName, engineVersion)
		if err != nil {
			return nil, err
		}
		container = fmt.Sprintf("?image=gcr.io/gamebuildr-151415/%v", imageName)
		halBaseURL = os.Getenv(config.HalMrRobotAPI)
	} else {
		return nil, fmt.Errorf("Invalid system requested %v %v %v", subSystem, engineName, engineVersion)
	}

	clientScaler := scaler.HTTPScaler{
		Client:        &http.Client{},
		LoadAPIUrl:    halBaseURL + "api/container/count",
		AddLoadAPIUrl: halBaseURL + "api/container/run" + container,
	}

	return clientScaler, nil
}
