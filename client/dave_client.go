package client

import (
	"net/http"
	"os"

	"fmt"

	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/Gamebuildr/Dave/pkg/scaler"
	"github.com/Gamebuildr/Dave/pkg/watcher"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/papertrail"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// DaveClient scales microservices remotely
type DaveClient struct {
	Watcher watcher.QueueMonitor
	Log     logger.Log
}

const logFileName string = "dave_client_"

// Create a new DaveClient
func (client *DaveClient) Create() {
	// New logger
	fileLogger := papertrail.PapertrailLogSave{
		App: "Dave",
		URL: os.Getenv(config.LogEndpoint),
	}
	logs := logger.SystemLogger{LogSave: fileLogger}

	// AWS SQS session
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(os.Getenv(config.Region))

	// setup client
	client.Watcher.Queue = watcher.AmazonWatcher{Client: sqs.New(sess)}
	client.Log = logs

	client.Log.Info("Dave Setup and Running")
}

// RunClient will read queue messages and send a response to an api endpoint to
// start up systems if the load is less than the max load
func (client *DaveClient) RunClient(system *scaler.ScalableSystem, queueURL string) {
	messageCount, err := client.Watcher.Queue.ReadQueueMessagesCount(queueURL)
	if err != nil {
		client.Log.Error(err.Error())
	}
	if messageCount <= 0 {
		return
	}
	load, err := system.System.GetSystemLoad()
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	if load > system.MaxLoad {
		client.Log.Info("system at max load")
		return
	}
	resp, err := system.System.AddSystemLoad()
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Adding container load failed with response %v", resp.Status)
		client.Log.Error(err.Error())
	}
}
