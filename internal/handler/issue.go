package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/pkg/database/verification"
	"github.com/vranystepan/email/pkg/messages/issue"
	"github.com/vranystepan/email/pkg/messages/register"
	"github.com/vranystepan/email/pkg/service"
)

// Issue creates a verification request based on the data from SQS
func Issue(dynamo service.DynamoDB, sqs service.SQS, table string, queue string) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			payload, err := register.Parse(message.Body)
			if err != nil {
				log.
					WithField("error", err).
					Error("could not parse incoming message")
				return err
			}
			log.
				WithField("email", payload.Email).
				Info("processing payload")

			// prepare database item
			item := verification.Item{
				Email: payload.Email,
			}

			// here I need to create some throttling
			valid, err := item.Check(dynamo, table)
			if err != nil {
				log.
					WithField("error", err).
					WithField("email", payload.Email).
					Error("could not validate")
				return err
			}

			// if verification request is not valid, discard it
			if !valid {
				log.
					WithField("email", payload.Email).
					Error("verification request can't be sent now")
				return nil
			}

			// save item to db
			err = item.Save(dynamo, table)
			if err != nil {
				log.
					WithField("error", err).
					WithField("email", payload.Email).
					Error("could not save item")
				return err
			}

			// here I need to send message to email system
			err = issue.New(item.Email).Send(sqs, queue)
			if err != nil {
				log.
					WithField("error", err).
					WithField("email", item.Email).
					Error("could not send message")
				return err
			}
			log.
				WithField("email", item.Email).
				Info("message sent")
		}
		return nil
	}
}
