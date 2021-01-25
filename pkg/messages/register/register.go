package register

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/pkg/service"
)

// Payload holds registration data
type Payload struct {
	TimeCreated time.Time `json:"timeCreated"`
	Source      string    `json:"source"`
	Email       string    `json:"email"`
}

// Parse converts JSON paylod to Payload struct
func Parse(message string) (Payload, error) {
	var m Payload
	err := json.Unmarshal([]byte(message), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// JSON converts RegisterPayload data to JSON string
func (r Payload) JSON() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Send sends message to the given SQS queue
func (r Payload) Send(svc service.SQS, queue string) error {
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
	log.WithField("payload", r).Info("sending message")

	_, err = svc.SendMessage(message)
	return err
}
