package dynamo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DeleteTable(dynamo *dynamodb.DynamoDB, name string) error {
	_, err := dynamo.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(name)})
	return err
}

func CreateTables(dynamo *dynamodb.DynamoDB, deleteTable bool, directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("reading files %s %v", directory, err)
	}

	for _, f := range files {
		if err := CreateTable(dynamo, deleteTable, path.Join(directory, f.Name())); err != nil {
			return err
		}
	}

	return nil
}

func CreateTable(dynamo *dynamodb.DynamoDB, deleteTable bool, file string) error {
	t, err := loadTableSchema(file)
	if err != nil {
		return err
	}

	if err := DeleteTable(dynamo, *t.TableName); err != nil {
		if !strings.Contains(err.Error(), dynamodb.ErrCodeResourceNotFoundException) {
			return fmt.Errorf("deleting table %v", err)
		}
	}

	if _, err = dynamo.CreateTable(t); err != nil {
		return fmt.Errorf("couldn't create table %v", err)
	}

	return nil
}

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

func loadTableSchema(file string) (*dynamodb.CreateTableInput, error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading schema file %s %v", file, err)
	}

	t := &dynamodb.CreateTableInput{}
	if err = json.Unmarshal(bs, t); err != nil {
		return nil, fmt.Errorf("unmarshaling schema %s %v", file, err)
	}

	return t, nil
}
