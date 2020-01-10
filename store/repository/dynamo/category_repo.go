package dynamo

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBCategoryRepository struct {
	table  string
	dynamo *dynamodb.DynamoDB
}

func NewCategoryRepository(table string, dyn *dynamodb.DynamoDB) *DynamoDBCategoryRepository {
	return &DynamoDBCategoryRepository{
		table:  table,
		dynamo: dyn,
	}
}

func (r *DynamoDBCategoryRepository) GetAll() ([]string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.table),
	}

	resp, err := r.dynamo.Scan(input)
	if err != nil {
		return nil, err
	}

	arr := make([]string, len(resp.Items))
	for i, r := range resp.Items {
		arr[i] = *r["category"].S
	}

	sort.Strings(arr)
	return arr, nil
}

func (r *DynamoDBCategoryRepository) Upsert(category string) error {
	av := make(map[string]*dynamodb.AttributeValue)

	av["category"] = &dynamodb.AttributeValue{S: aws.String(category)}

	_, err := r.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	})

	return err
}

func (r *DynamoDBCategoryRepository) Delete(category string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"category": &dynamodb.AttributeValue{S: aws.String(category)},
		},
	}

	_, err := r.dynamo.DeleteItem(input)

	return err
}
