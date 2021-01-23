package service

import (
	sqsService "github.com/aws/aws-sdk-go/service/sqs"
)

// SQS contains functions needed only in this codebase
type SQS interface {
	SendMessage(*sqsService.SendMessageInput) (*sqsService.SendMessageOutput, error)
}
