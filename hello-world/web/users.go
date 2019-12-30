package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type UserDao interface {
	RegisterUser(username, password string) error
	GetPassword(username string) (string, error)
}

// UserUnencryptedDao implementation of UserDao which keeps
// user data in DynamoDB unencrypted.
// It's meant to be a temporary solution for the prototype
type UserUnencryptedDao struct {
	tableName string
}

func (t *UserUnencryptedDao) RegisterUser(username, password string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
			"password": {
				S: aws.String(password),
			},
		},
		TableName:              aws.String(t.tableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}

	_, e := svc.PutItem(input)
	if e != nil {
		return e
	}
	return nil
}

func (t *UserUnencryptedDao) GetPassword(username string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	result, e := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(t.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				N: aws.String(username),
			},
		},
	})
	if e != nil {
		return "", e
	}
	return *result.Item["password"].S, nil
}

