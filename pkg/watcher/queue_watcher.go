package watcher

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gamebuildr/Dave/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// QueueWatcher is the aws sns implementation of the NotificationService
type QueueWatcher struct {
	client sqsiface.SQSAPI
}

var messageIDToReceiptHandle = map[string]*string{}

// Setup gets the watcher ready to use
func (watcher *QueueWatcher) Setup() {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(os.Getenv(config.Region))
	watcher.client = sqs.New(sess)
}

// ReadNextMessage from the specified queue
func (watcher QueueWatcher) ReadNextMessage(url string) (*MessageInfo, error) {
	response, err := watcher.getMessageFromSqs(url)

	if err != nil {
		return new(MessageInfo), err
	}

	messages := response.Messages

	if len(messages) == 0 {
		return new(MessageInfo), nil
	}

	message := messages[0]
	var messageInfo MessageInfo

	data := []byte(*message.Body)
	json.Unmarshal(data, &messageInfo)
	messageIDToReceiptHandle[messageInfo.MessageID] = message.ReceiptHandle

	return &messageInfo, err
}

// DeleteMessage with ID from the specified queue
func (watcher QueueWatcher) DeleteMessage(messageID string, url string) error {
	handle, exists := messageIDToReceiptHandle[messageID]
	if !exists {
		return fmt.Errorf("Message requested for delete does not exist: %v", messageID)
	}

	deleteMsg := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(url),
		ReceiptHandle: handle,
	}

	_, err := watcher.client.DeleteMessage(deleteMsg)
	if err != nil {
		return err
	}

	return nil
}

func (watcher *QueueWatcher) getMessageFromSqs(url string) (*sqs.ReceiveMessageOutput, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}

	response, err := watcher.client.ReceiveMessage(params)

	return response, err
}
