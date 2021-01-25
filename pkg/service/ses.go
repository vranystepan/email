package service

import "github.com/aws/aws-sdk-go/service/ses"

// SES contains AWS DynamoDB functions needed only in this codebase
type SES interface {
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}
