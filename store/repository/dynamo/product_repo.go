package dynamo

import (
	"fmt"
	"strconv"

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

func NewProductRepository(table string, configProvider client.ConfigProvider, config *aws.Config) *DynamoDBProductRepository {
	return &DynamoDBProductRepository{
		table:  table,
		dynamo: dynamodb.New(configProvider, config),
	}
}

func (r *DynamoDBProductRepository) CreateTable() error {

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{{
			AttributeName: aws.String("id"),
			AttributeType: aws.String("S"),
		}},
		KeySchema: []*dynamodb.KeySchemaElement{{
			AttributeName: aws.String("id"),
			KeyType:       aws.String("HASH"),
		}},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{{
			IndexName: aws.String("PriceIndex"),
			KeySchema: []*dynamodb.KeySchemaElement{{
				AttributeName: aws.String("price"),
				KeyType:       aws.String("RANGE"),
			}},
		}},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(r.table),
	}

	if _, err := r.dynamo.CreateTable(input); err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBProductRepository) DeleteTable() error {
	_, err := r.dynamo.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(r.table)})

	return err
}

func (r *DynamoDBProductRepository) GetAll() ([]*store.Product, error) {
	resp, err := r.dynamo.Scan(&dynamodb.ScanInput{
		TableName: aws.String(r.table),
	})
	if err != nil {
		return nil, err
	}

	return mapItems(resp.Items)
}

func (r *DynamoDBProductRepository) GetPaged(start string, limit int) ([]*store.Product, int, error) {
	//fake
	prds, err := r.GetAll()
	return prds, 100, err
}

func (r *DynamoDBProductRepository) Get(productID string) (*store.Product, error) {
	i, err := r.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String(productID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("could not get item %v", err)
	}

	if i.Item == nil {
		return nil, nil
	}

	return mapItem(i.Item)
}

func (r *DynamoDBProductRepository) Add(product *store.Product) error {

	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return err
	}

	if _, err = r.dynamo.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	}); err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBProductRepository) Upsert(product *store.Product) error {

	av, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.table),
	}

	_, err = r.dynamo.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBProductRepository) Delete(productID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String(productID)},
		},
	}

	if _, err := r.dynamo.DeleteItem(input); err != nil {
		return err
	}

	return nil
}

func mapItems(items []map[string]*dynamodb.AttributeValue) ([]*store.Product, error) {
	pr := make([]*store.Product, len(items))

	for i, item := range items {
		p, err := mapItem(item)
		if err != nil {
			return nil, err
		}
		pr[i] = p
	}

	return pr, nil
}

func mapItem(item map[string]*dynamodb.AttributeValue) (*store.Product, error) {
	p := &store.Product{}

	p.ID = *item["id"].S
	p.Name = *item["name"].S
	p.Description = string(item["description"].B)
	price, err := strconv.ParseFloat(*item["price"].N, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse price %v", err)
	}
	p.Price = price

	return p, err
}
