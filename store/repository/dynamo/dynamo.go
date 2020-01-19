package dynamo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func IsTableActive(dynamo *dynamodb.DynamoDB, table string, timeout time.Duration) bool {
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

func Health(dynamo *dynamodb.DynamoDB, table string, timeout time.Duration) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req := &dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		}
		resp, err := dynamo.DescribeTableWithContext(ctx, req)
		if err != nil {
			return err
		}
		if *resp.Table.TableStatus != "ACTIVE" {
			return fmt.Errorf("table is not active %s", *resp.Table.TableStatus)
		}
		return nil
	}
}
