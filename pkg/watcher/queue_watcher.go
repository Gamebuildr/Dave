package watcher

import (
	"os"

	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// QueueWatcher is the aws sns implementation of the NotificationService
type QueueWatcher struct {
	Client sqsiface.SQSAPI
}

// Setup gets the watcher ready to use
func (watcher QueueWatcher) Setup() {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(os.Getenv(config.Region))
	watcher.Client = sqs.New(sess)
}

// ReadNextMessage from the specified amazon queue
func (watcher QueueWatcher) ReadNextMessage(url string) (int, error) {
	params := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}
	response, err := watcher.Client.ReceiveMessage(&params)
	if err != nil {
		return 0, err
	}
	return len(response.Messages), nil
}
