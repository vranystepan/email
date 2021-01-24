package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"

	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/internal/handler"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	// initialize AWS SDK and services
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamoSvc := dynamodb.New(sess)
	sqsSvc := sqs.New(sess)

	// obtain configuration params
	table := os.Getenv("CONFIG_DYNAMODB_TABLE_NAME")
	queue := os.Getenv("CONFIG_SQS_EMAIL_ISSUE_QUEUE_URL")

	// start lambda function
	lambda.Start(handler.Issue(dynamoSvc, sqsSvc, table, queue))
}
