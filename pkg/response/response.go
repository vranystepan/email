package response

import "github.com/aws/aws-lambda-go/events"

// BadRequest can be directly returned in case of some invalid data
func BadRequest(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 400,
	}, nil
}

// TemporaryRedirect can be directly returned in case of temporary redirection (code 307)
func TemporaryRedirect(location string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Location": location,
		},
		StatusCode: 307,
	}, nil
}
