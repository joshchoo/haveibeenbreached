package haveibeenbreached

import (
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
	QueueURL    string
}

func (m MessageQueue) SendMessage(input SendMessageInput) error {
	_, err := m.queue.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(input.MessageBody),
		QueueUrl:    &input.QueueURL,
	})
	if err != nil {
		return err
	}
	return nil
}
