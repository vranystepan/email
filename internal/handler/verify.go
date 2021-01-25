package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/pkg/messages/email/verify"
	"github.com/vranystepan/email/pkg/messages/issue"
	"github.com/vranystepan/email/pkg/service"
)

// VerifyParams contains params for Verify function to minify signature
type VerifyParams struct {
	SES service.SES
}

// Verify sends verification email via SES
func Verify(p VerifyParams) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			payload, err := issue.Parse(message.Body)
			if err != nil {
				log.
					WithField("error", err).
					Error("could not parse incoming message")
				return err
			}
			log.
				WithField("email", payload.Email).
				Info("processing payload")

			email := verify.Email{
				ToAddresses: payload.Email,
				Sender:      "stepan@vrany.dev",
				HTML:        "<b>hello</b>",
				Text:        "hello",
				Subject:     "verification!",
			}
			err = email.Send(p.SES)
			if err != nil {
				log.
					WithField("error", err).
					Error("could not send email message")
				return err
			}
		}
		return nil
	}
}
