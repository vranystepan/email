package issue

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	log "github.com/sirupsen/logrus"
	"github.com/vranystepan/email/pkg/service"
)

// Item is representation of verification record
type Item struct {
	Email string
	// Last is the last time of verification request
	Last time.Time
}

// Check checks if the previous request is older than 5 minutes / 300 seconds
func (i *Item) Check(dynamo service.DynamoDB, table string) (bool, error) {

	// prepare input for search
	input := prepareQuery(table, i.Email)

	// get item from DB
	resp, err := dynamo.GetItem(input)
	if err != nil {
		log.
			WithField("email", i.Email).
			Error("failed to get item from database")
		return false, err
	}

	// if len is 0 then email does not exist in the database yet - safe to process
	if len(resp.Item) == 0 {
		log.
			WithField("email", i.Email).
			Info("item does not exist yet")
		return true, err
	}

	// get the actual item from DB
	err = dynamodbattribute.UnmarshalMap(resp.Item, i)
	if err != nil {
		log.
			WithField("email", i.Email).
			Error("could not unmasrshall DB entry")
		return false, err
	}

	// check the time limit here
	log.
		WithField("email", i.Email).
		WithField("last", i.Last).
		Info("email is already in DB")
	diff := time.Now().Sub(i.Last)

	// duration can be parametrized in future
	if diff < time.Minute*5 {
		log.
			WithField("email", i.Email).
			WithField("last", i.Last).
			WithField("diff", diff.String()).
			Info("email has been already sent")
		return false, nil
	}

	return true, err
}

// Save updates item in the DynamoDB table
func (i Item) Save(dynamo service.DynamoDB, table string) error {
	// update time
	i.Last = time.Now()
	// serialize item
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		log.WithField("email", i.Email).Error("failed to serialize verification item")
		return err
	}

	// prepare put input
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	// put item
	_, err = dynamo.PutItem(input)
	if err != nil {
		log.WithField("email", i.Email).Error("failed to put item")
		return err
	}

	return nil
}

func prepareQuery(table string, email string) *dynamodb.GetItemInput {
	return &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	}
}
