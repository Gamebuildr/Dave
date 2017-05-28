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
						Body:          aws.String(`{"message":"one"}`),
						ReceiptHandle: &messageReceipt,
					},
					{
						Body:          aws.String(`{"message":"one"}`),
						ReceiptHandle: &messageReceipt,
					},
				},
			},
		},
	}
	for i, c := range mockMessages {
		queue := QueueWatcher{Client: MockedAmazonClient{Response: c.Resp}}
		count, err := queue.ReadNextMessage(fmt.Sprintf("mockUrl_%d", i))
		if err != nil {
			t.Fatalf("%d, amazon test error, %v", i, err)
		}
		if count != 2 {
			t.Errorf("Expected a count of %v, got %v", 2, count)
		}
	}
}
