package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/internal/handler"
)

//nolint
func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	// initialize AWS SDK and services
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sesSvc := ses.New(sess)

	// create params struct for Issue handler
	params := handler.VerifyParams{
		SES: sesSvc,
	}

	// start lambda function
	lambda.Start(handler.Verify(params))
}
