package response

import "github.com/aws/aws-lambda-go/events"

func BadRequest(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 400,
	}, nil
}

func TemporaryRedirect(location string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Location": location,
		},
		StatusCode: 307,
	}, nil
}
