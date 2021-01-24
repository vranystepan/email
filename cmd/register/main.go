package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

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
	sqsSvc := sqs.New(sess)

	// obtain configuration params
	queue := os.Getenv("CONFIG_SQS_EMAIL_REGISTRATION_QUEUE_URL")

	// create params struct for Issue handler
	params := handler.RegisterParams{
		SQS:   sqsSvc,
		Queue: queue,
	}

	// start lambda function
	lambda.Start(handler.Register(params))
}
