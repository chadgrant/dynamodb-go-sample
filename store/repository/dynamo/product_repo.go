package dynamo

import (
	"fmt"

	"github.com/chadgrant/dynamodb-go-sample/store"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBProductRepository struct {
	table  string
	dynamo *dynamodb.DynamoDB
}

func NewProductRepository(table string, sess client.ConfigProvider, config *aws.Config) *DynamoDBProductRepository {
	return &DynamoDBProductRepository{
		table:  table,
		dynamo: dynamodb.New(sess, config),
	}
}

func (r *DynamoDBProductRepository) GetPaged(category string, limit int, lastID string, lastPrice float64) ([]*store.Product, int64, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.table),
		IndexName:              aws.String("price-index"),
		KeyConditionExpression: aws.String("category = :c"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {S: aws.String(category)},
		},
		Limit:            aws.Int64(int64(limit)),
		ScanIndexForward: aws.Bool(false),
	}
	if len(lastID) > 0 {

		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":       {S: aws.String(lastID)},
			"category": {S: aws.String(category)},
			"price":    {N: aws.String(fmt.Sprint(lastPrice))},
		}
	}
	resp, err := r.dynamo.Query(input)
	if err != nil {
		return nil, 0, err
	}

	prds := make([]*store.Product, len(resp.Items))
	for i, item := range resp.Items {
		p := &store.Product{}
		if err := dynamodbattribute.UnmarshalMap(item, &p); err != nil {
			return nil, 0, fmt.Errorf("error mapping item %v", err)
		}
		prds[i] = p
	}

	return prds, *resp.Count, nil
}

func (r *DynamoDBProductRepository) Get(productID string) (*store.Product, error) {
	result, err := r.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String(productID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not get item %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	p := &store.Product{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &p); err != nil {
		return nil, fmt.Errorf("error mapping item %v", err)
	}
	return p, nil
}

func (r *DynamoDBProductRepository) Upsert(category string, product *store.Product) error {

	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return fmt.Errorf("error marshalling %v", err)
	}

	av["category"] = &dynamodb.AttributeValue{S: aws.String(category)}

	_, err = r.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	})

	return err
}

func (r *DynamoDBProductRepository) Delete(productID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String(productID)},
		},
	}

	_, err := r.dynamo.DeleteItem(input)

	return err
}