package client

import (
	"os"

	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/Gamebuildr/Dave/pkg/watcher"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type DaveClient struct {
	Watcher watcher.AmazonWatcher
}

func (client *DaveClient) Create() {
	sess := session.Must(session.NewSession())
	client.Watcher = watcher.AmazonWatcher{Client: sqs.New(sess)}
	client.Watcher.ReadQueueMessages(os.Getenv(config.gogetaSQSEndpoint))
}
