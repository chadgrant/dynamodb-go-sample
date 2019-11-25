package dynamo

import (
	"fmt"

	"github.com/chadgrant/dynamodb-go-sample/store"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBProductRepository struct {
	table  string
	dynamo *dynamodb.DynamoDB
}

func NewProductRepository(endpoint, table string) *DynamoDBProductRepository {

	svc := dynamodb.New(session.Must(session.NewSession()), &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
		Endpoint:    aws.String(endpoint),
	})

	return &DynamoDBProductRepository{
		table:  table,
		dynamo: svc,
	}
}

func (r *DynamoDBProductRepository) CreateTable() error {

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(r.table),
	}

	_, err := r.dynamo.CreateTable(input)
	if err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBProductRepository) DeleteTable() error {
	input := &dynamodb.DeleteTableInput{TableName: aws.String(r.table)}

	_, err := r.dynamo.DeleteTable(input)
	if err != nil {
		return err
	}

	return nil
}

func (r *DynamoDBProductRepository) GetAll() ([]*store.Product, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.table),
	}
	rsp, err := r.dynamo.Scan(input)
	if err != nil {
		return nil, err
	}

	pr := make([]*store.Product, 0)

	for _, i := range rsp.Items {
		p := &store.Product{}
		p.ID = *i["id"].S
		p.Name = *i["name"].S
		//p.Description = *i["description"].S
		//p.Price = *i["price"].N
		pr = append(pr, p)
	}

	return pr, nil
}

func (r *DynamoDBProductRepository) Get(productID string) (*store.Product, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{S: aws.String(productID)},
		},
	}

	i, err := r.dynamo.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("could not get item %v", err)
	}

	if i.Item == nil {
		return nil, nil
	}

	p := &store.Product{}
	p.ID = *i.Item["id"].S
	p.Name = *i.Item["name"].S

	return p, nil
}

func (r *DynamoDBProductRepository) Add(product *store.Product) error {

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

	_, err := r.dynamo.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
