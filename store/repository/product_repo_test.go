package repository

import (
	"os"
	"testing"

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

	sess, config := session.Must(session.NewSession()), &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("123", "123", ""),
		Endpoint:    aws.String(ep),
	}

	if err := dynamo.CreateTables(dynamodb.New(sess, config), true, "../../data/schema"); err != nil {
		t.Fatalf("creating tables %v", err)
		return
	}

	runTests(dynamo.NewProductRepository("products", sess, config), t)
}

func runTests(repo ProductRepository, t *testing.T) {
	if err := setup(repo); err != nil {
		t.Fatalf("setup failed %v", err)
	}

	t.Run("Add", func(t *testing.T) {
		testAdd(repo, t)
	})

	t.Run("GetSingle", func(t *testing.T) {
		testGet(repo, t)
	})

	t.Run("GetPaged", func(t *testing.T) {
		testGetPaged(repo, t)
	})

	t.Run("Upsert", func(t *testing.T) {
		testUpsert(repo, t)
	})

	t.Run("Delete", func(t *testing.T) {
		testDelete(repo, t)
	})
}

func setup(repo ProductRepository) error {
	populator := NewPopulator(repo)

	// if err := populator.Create(100); err != nil {
	// 	return err
	// }

	// if err := populator.Export("../../data/products.json"); err != nil {
	// 	return err
	// }

	if err := populator.Load("../../data/products.json"); err != nil {
		return err
	}

	return nil
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
	var last string
	lastPrice := float64(0)
	total, visited := int64(0), int64(0)

	for {
		products, total, err = repo.GetPaged(categories[0], 25, last, lastPrice)
		if err != nil {
			t.Fatal(err)
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

		if len(products) < 25 {
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
