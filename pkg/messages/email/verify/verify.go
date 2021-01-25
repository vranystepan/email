package verify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/pkg/service"
)

// Email contains all information necessary to send an email
type Email struct {
	ToAddresses string
	Sender      string
	Subject     string
	HTML        string
	Text        string
}

// Send sends an email via SES
func (e Email) Send(sesSvc service.SES) error {
	input := assembleSESInput(e)

	log.
		WithField("sender", e.Sender).
		WithField("ToAddresses", e.ToAddresses).
		Info("sending message")

	_, err := sesSvc.SendEmail(input)

	return err
}

// assembleSESInput creates a input for SES service
func assembleSESInput(e Email) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(e.ToAddresses),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.HTML),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(e.Text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.Sender),
	}
}
