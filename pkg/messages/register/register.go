package register

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSService interface {
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type RegisterPayload struct {
	TimeCreated time.Time `json:"timeCreated"`
	Source      string    `json:"source"`
	Email       string    `json:"email"`
}

func (r RegisterPayload) JSON() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r RegisterPayload) Send(svc SQSService, queue string) error {
	// convert body to JSON
	body, err := r.JSON()
	if err != nil {
		return err
	}

	// send message to message broker
	message := &sqs.SendMessageInput{
		//MessageAttributes: map[string]*sqs.MessageAttributeValue{},
		MessageBody: aws.String(body),
		QueueUrl:    aws.String(queue),
	}
	_, err = svc.SendMessage(message)
	return err
}
