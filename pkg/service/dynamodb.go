package service

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB contains AWS DynamoDB functions needed only in this codebase
type DynamoDB interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}
