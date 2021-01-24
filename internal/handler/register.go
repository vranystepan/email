package handler

import (
	"context"
	"time"

	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vranystepan/email/pkg/email"
	"github.com/vranystepan/email/pkg/messages/register"
	"github.com/vranystepan/email/pkg/response"
	"github.com/vranystepan/email/pkg/service"
)

// RegisterParams contains params for Register function to minify signature
type RegisterParams struct {
	SQS   service.SQS
	Queue string
}

// Register is the main handler for the mail registration service
func Register(p RegisterParams) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { //nolint
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
		payload := register.Payload{
			Email:       userEmail,
			Source:      referer,
			TimeCreated: time.Now(),
		}

		// send payload to message broker
		err = payload.Send(p.SQS, p.Queue)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return response.TemporaryRedirect(referer)
	}
}
