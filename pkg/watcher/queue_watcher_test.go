package watcher

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// MockedAmazonClient allows mocking of the Amazon SQS client
type MockedAmazonClient struct {
	sqsiface.SQSAPI
	Response sqs.ReceiveMessageOutput
}

func (m MockedAmazonClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &m.Response, nil
}

func TestGetQueueMessages(t *testing.T) {
	messageReceipt := "mockReceipts"
	mockMessages := []struct {
		Resp sqs.ReceiveMessageOutput
	}{
		{
			Resp: sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						Body:          aws.String(`{"MessageId": "f6720a18-fb39-56e2-8504-5bbf2f75c021", "message":"one"}`),
						ReceiptHandle: &messageReceipt,
					},
					{
						Body:          aws.String(`{"MessageId": "f6720a18-fb39-56e2-8504-5bbf2f75c022", message":"one"}`),
						ReceiptHandle: &messageReceipt,
					},
				},
			},
		},
	}
	for i, c := range mockMessages {
		queue := QueueWatcher{client: MockedAmazonClient{Response: c.Resp}}
		messageInfo, err := queue.ReadNextMessage(fmt.Sprintf("mockUrl_%d", i))
		if err != nil {
			t.Fatalf("%d, amazon test error, %v", i, err)
		}
		if messageInfo.MessageID != "f6720a18-fb39-56e2-8504-5bbf2f75c021" {
			t.Errorf("Expected ID %v, got %v", "f6720a18-fb39-56e2-8504-5bbf2f75c021", messageInfo.MessageID)
		}
		if messageInfo.Message != "one" {
			t.Errorf("Expected message %v, got %v", "one", messageInfo.Message)
		}
	}
}
