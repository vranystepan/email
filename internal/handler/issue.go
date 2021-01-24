package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	db "github.com/vranystepan/email/pkg/database/issue"
	"github.com/vranystepan/email/pkg/messages/issue"
	"github.com/vranystepan/email/pkg/messages/register"
	"github.com/vranystepan/email/pkg/service"
)

// IssueParams contains params for Issue function to minify signature
type IssueParams struct {
	Dynamo service.DynamoDB
	SQS    service.SQS
	Table  string
	Queue  string
}

// Issue creates a verification request based on the data from SQS
func Issue(p IssueParams) func(ctx context.Context, sqsEvent events.SQSEvent) error {
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
			item := db.Item{
				Email: payload.Email,
			}

			// here I need to create some throttling
			valid, err := item.Check(p.Dynamo, p.Table)
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
					Error("can't be sent now, discarding")
				return nil
			}

			// save item to db
			err = item.Save(p.Dynamo, p.Table)
			if err != nil {
				log.
					WithField("error", err).
					WithField("email", payload.Email).
					Error("could not save item")
				return err
			}

			err = issue.New(item.Email).Send(p.SQS, p.Queue)
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
