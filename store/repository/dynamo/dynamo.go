package dynamo

import (
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
