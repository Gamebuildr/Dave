package watcher

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// AmazonWatcher is the aws sns implementation of the NotificationService
type AmazonWatcher struct {
	Client sqsiface.SQSAPI
}

// ReadQueueMessages from the specified amazon queue
func (watcher AmazonWatcher) ReadQueueMessages(url string) (int, error) {
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
