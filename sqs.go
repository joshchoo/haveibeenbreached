package haveibeenbreached

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageQueue struct {
	queue *sqs.SQS
}

func NewMessageQueue(queue *sqs.SQS) MessageQueue {
	return MessageQueue{queue}
}

type SendMessageInput struct {
	MessageBody string
	QueueName   string
}

func (m MessageQueue) SendMessage(input SendMessageInput) error {
	result, err := m.queue.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &input.QueueName,
	})
	if err != nil {
		return err
	}
	if result.QueueUrl == nil {
		return fmt.Errorf("queue URL is nil")
	}
	_, err = m.queue.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(input.MessageBody),
		QueueUrl:    result.QueueUrl,
	})
	if err != nil {
		return err
	}
	return nil
}
