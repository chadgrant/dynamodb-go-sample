package repository

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/google/uuid"
)

func TestMock(t *testing.T) {
	runTests(NewMockProductRepository(), t)
}

func TestIntegration(t *testing.T) {
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
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
		Endpoint:    aws.String(ep),
	})

	if err := dynamo.CreateTables(dyn, true, "../../data/schema"); err != nil {
		t.Fatalf("couldn't create tables %v", err)
	}

	if !dynamo.IsTableActive(dyn, "products", 5*time.Second) {
		t.Fatal("timeout waiting for table to be active")
	}

	runTests(dynamo.NewProductRepository("products", dyn), t)
}

func runTests(repo ProductRepository, t *testing.T) {
	pop := NewPopulator(repo)

	// if err := pop.Create(100); err != nil {
	// 	t.Fatalf("couldnt create data %v", err)
	// }

	// if err := pop.Export("../../data/products.json"); err != nil {
	// 	t.Fatalf("couldnt export data %v", err)
	// }

	if err := pop.Load("../../data/products.json"); err != nil {
		t.Fatalf("couldnt load data %v", err)
	}

	t.Run("GetSingle", func(t *testing.T) {
		testGet(repo, t)
	})

	t.Run("GetPaged", func(t *testing.T) {
		testGetPaged(repo, t)
	})

	t.Run("Upsert", func(t *testing.T) {
		testUpsert(repo, t)
	})

	t.Run("Add", func(t *testing.T) {
		testAdd(repo, t)
	})

	t.Run("Delete", func(t *testing.T) {
		testDelete(repo, t)
	})
}

func testAdd(repo ProductRepository, t *testing.T) {
	id, _ := uuid.NewRandom()
	p := &store.Product{
		ID:    id.String(),
		Name:  "Test Product " + id.String(),
		Price: 1.00,
	}

	if err := repo.Upsert(categories[0], p); err != nil {
		t.Error(err)
	}
}

func testGetPaged(repo ProductRepository, t *testing.T) {

	dupes := make(map[string]*store.Product)
	var products []*store.Product
	var err error
	last := ""
	lastPrice := float64(0)
	total, visited := int64(0), int64(0)
	size := 25

	for {
		products, total, err = repo.GetPaged(categories[0], size, last, lastPrice)
		if err != nil {
			t.Fatalf("get paged %v", err)
		}

		for i, p := range products {
			if i > 0 {
				if products[i-1].Price < products[i].Price {
					t.Fatal("products not sorted")
				}
			}
			visited++
			if dupes[p.ID] != nil {
				t.Errorf("duplicate %s", p.ID)
			}
			dupes[p.ID] = p
		}

		if len(products) < size {
			break
		}

		last = products[len(products)-1].ID
		lastPrice = products[len(products)-1].Price
	}

	if visited < total {
		t.Errorf("did not visit all items expected %d got %d", total, visited)
	}
}

func testGet(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Fatalf("could not get products %v", err)
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
	}

	for _, v := range ps {
		p, err := repo.Get(v.ID)
		if err != nil {
			t.Fatalf("getting product %s %v", v.ID, err)
		}
		if p.ID != v.ID {
			t.Fatalf("wrong product returned expected %s got %s", v.ID, p.ID)
		}
	}
}

func testUpsert(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Fatalf("getting products %v", err)
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
	}

	p, err := repo.Get(ps[0].ID)
	if err != nil {
		t.Fatalf("getting product %s %v", ps[0].ID, err)
	}

	p.Name = p.Name + " Updated"

	if err := repo.Upsert(categories[0], p); err != nil {
		t.Fatalf("upserting product %v", err)
	}

	n, err := repo.Get(p.ID)
	if err != nil {
		t.Fatalf("could not get updated product %v", err)
	}

	if n.Name != p.Name {
		t.Fatalf("product name not updated got %s expected %s", n.Name, p.Name)
	}
}

func testDelete(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Fatalf("could not get products %v", err)
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
	}

	if err := repo.Delete(ps[0].ID); err != nil {
		t.Fatalf("deleting product %s %v", ps[0].ID, err)
	}

	p, err := repo.Get(ps[0].ID)
	if err != nil {
		t.Fatalf("getting product %s %v", ps[0].ID, err)
	}

	if p != nil {
		t.Fatalf("product not deleted %s", p.ID)
	}
}
