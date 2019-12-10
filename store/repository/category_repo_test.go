package repository

import (
	"os"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
)

func TestCategoryMock(t *testing.T) {
	c := NewMockCategoryRepository("Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory")
	runTests(c, NewMockProductRepository(c, 100), t)
}

func TestCategoryDynamoIntegration(t *testing.T) {
	if len(os.Getenv("TEST_INTEGRATION")) == 0 {
		t.Log("Skipping integration tests, TEST_INTEGRATION environment variable not set")
		return
	}

	ep := os.Getenv("DYNAMO_ENDPOINT")
	if len(ep) == 0 {
		ep = "http://localhost:8000"
	}

	dyn := dynamodb.New(session.Must(session.NewSession()), &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("key", "secret", ""),
		Endpoint:    aws.String(ep),
	})

	runCategoryTests(dynamo.NewCategoryRepository("categories", dyn), t)
}

func runCategoryTests(repo CategoryRepository, t *testing.T) {

	t.Run("GetAll", func(t *testing.T) {
		testCategoryGet(repo, t)
	})

	t.Run("Upsert", func(t *testing.T) {
		testCategoryUpsert(repo, t)
	})

	t.Run("Delete", func(t *testing.T) {
		testCategoryDelete(repo, t)
	})
}

func testCategoryGet(repo CategoryRepository, t *testing.T) {
	cats, err := repo.GetAll()
	if err != nil {
		t.Errorf("error getting categories %v", err)
	}
	if len(cats) == 0 {
		t.Errorf("no categories returned")
	}

	if !sort.StringsAreSorted(cats) {
		t.Errorf("categories are not sorted")
	}
}

func testCategoryUpsert(repo CategoryRepository, t *testing.T) {
	cats, err := repo.GetAll()
	if err != nil {
		t.Errorf("error getting categories %v", err)
	}
	if len(cats) == 0 {
		t.Errorf("no categories returned")
	}

	if err := repo.Upsert("testing"); err != nil {
		t.Fatalf("upsert failed %v", err)
	}

	cats, err = repo.GetAll()
	for _, v := range cats {
		if v == "testing" {
			return
		}
	}
	t.Errorf("could not find new key upserted : testing")
}

func testCategoryDelete(repo CategoryRepository, t *testing.T) {
	cats, err := repo.GetAll()
	if err != nil {
		t.Errorf("error getting categories %v", err)
	}
	if len(cats) == 0 {
		t.Errorf("no categories returned")
	}

	if err := repo.Upsert("testing"); err != nil {
		t.Fatalf("upsert failed %v", err)
	}

	if err := repo.Delete("testing"); err != nil {
		t.Fatalf("upsert failed %v", err)
	}

	cats, err = repo.GetAll()
	for _, v := range cats {
		if v == "testing" {
			t.Errorf("should not have found category : testing")
		}
	}
}
