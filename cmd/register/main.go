package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/vranystepan/email/internal/handler"
)

func main() {
	// initialize AWS SDK and services
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)

	// obtain configuration params
	queue := os.Getenv("CONFIG_SQS_EMAIL_REGISTRATION_QUEUE_URL")

	// start lambda function
	lambda.Start(handler.Register(svc, queue))
}
