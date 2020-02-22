package dynamo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func New(region, accesskey, secret, endpoint string) *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(session.NewSession()), &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accesskey, secret, ""),
		Endpoint:    aws.String(endpoint),
	})
}

func WaitForTables(dynamo dynamodbiface.DynamoDBAPI, timeout time.Duration, tables ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, table := range tables {
		req := &dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		}

		if err := dynamo.WaitUntilTableExistsWithContext(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func IsTableActive(dynamo dynamodbiface.DynamoDBAPI, table string, timeout time.Duration) bool {
	tick := time.NewTicker(500 * time.Millisecond)
	timeoutC := time.After(timeout)
	defer tick.Stop()

	for {
		select {
		case <-timeoutC:
			return false

		case <-tick.C:
			if resp, err := dynamo.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(table)}); err == nil {
				if *resp.Table.TableStatus == "ACTIVE" {
					return true
				}
			}
		}
	}
}

func Health(dynamo dynamodbiface.DynamoDBAPI, timeout time.Duration, tables ...string) func() error {
	return func() error {
		return WaitForTables(dynamo, timeout, tables...)
	}
}
