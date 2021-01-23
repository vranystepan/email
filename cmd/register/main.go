package main

import (
	"context"
	"os"
	"time"

	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/vranystepan/email/pkg/email"
	"github.com/vranystepan/email/pkg/messages/register"
	"github.com/vranystepan/email/pkg/response"
)

func HandleRequest(svc register.SQSService, queue string) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		params, err := url.ParseQuery(req.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		// get important values from request
		userEmail := params.Get("email")
		referer := req.Headers["referer"]

		// validate email
		if !email.Valid(userEmail) {
			return response.BadRequest("incorrect email address")
		}

		// create registration payload
		payload := register.RegisterPayload{
			Email:       userEmail,
			Source:      referer,
			TimeCreated: time.Now(),
		}

		// send payload to message broker
		err = payload.Send(svc, queue)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return response.TemporaryRedirect(referer)
	}
}

func main() {
	// initialize AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sqs.New(sess)
	queue := os.Getenv("CONFIG_SQS_EMAIL_REGISTRATION_QUEUE_URL")

	lambda.Start(HandleRequest(svc, queue))
}
