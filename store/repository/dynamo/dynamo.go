package dynamo

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DeleteTable(dynamo *dynamodb.DynamoDB, name string) error {
	_, err := dynamo.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(name)})

	return err
}

func CreateTables(dynamo *dynamodb.DynamoDB, deleteTable bool, directory string) ([]string, error) {
	fs, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(fs))
	for i, f := range fs {
		name, err := CreateTable(dynamo, deleteTable, path.Join(directory, f.Name()))
		if err != nil {
			return nil, err
		}
		names[i] = name
	}

	return names, nil
}

func CreateTable(dynamo *dynamodb.DynamoDB, deleteTable bool, file string) (string, error) {
	t, err := getTableSchema(file)
	if err != nil {
		return "", err
	}

	if err := DeleteTable(dynamo, *t.TableName); err != nil {
		if !strings.Contains(err.Error(), dynamodb.ErrCodeResourceNotFoundException) {
			return "", err
		}
	}

	if _, err = dynamo.CreateTable(t); err != nil {
		return *t.TableName, err
	}

	return *t.TableName, nil
}

func getTableSchema(file string) (*dynamodb.CreateTableInput, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	t := &dynamodb.CreateTableInput{}
	if err = json.Unmarshal(bs, t); err != nil {
		return nil, err
	}

	return t, nil
}
