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
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
	"github.com/google/uuid"
)

func TestMockRepository(t *testing.T) {
	repo := mock.NewProductRepository()
	runTests(repo, t)
}

func TestIntegrationRepository(t *testing.T) {
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

	if _, err := dynamo.CreateTables(dynamodb.New(sess, config), true, "../../data/schema"); err != nil {
		t.Errorf("error creating tables %v", err)
		return
	}

	repo := dynamo.NewProductRepository("products", sess, config)
	runTests(repo, t)
}

func runTests(repo ProductRepository, t *testing.T) {
	if err := setup(repo); err != nil {
		t.Fatalf("setup failed %v", err)
	}

	t.Run("AddProduct", func(t *testing.T) {
		testAddProduct(repo, t)
	})

	t.Run("GetSingleProduct", func(t *testing.T) {
		testGetProduct(repo, t)
	})

	t.Run("GetProductsPaged", func(t *testing.T) {
		testGetProductsPaged(repo, t)
	})

	t.Run("UpsertProduct", func(t *testing.T) {
		testUpsertProduct(repo, t)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		testDeleteProduct(repo, t)
	})
}

func setup(repo ProductRepository) error {
	populator := NewPopulator(repo)

	// if err := populator.CreateProducts(100); err != nil {
	// 	return err
	// }

	// if err := populator.ExportProducts("../../data/products.json"); err != nil {
	// 	return err
	// }

	if err := populator.LoadProducts("../../data/products.json"); err != nil {
		return err
	}

	return nil
}

func testAddProduct(repo ProductRepository, t *testing.T) {
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

func testGetProductsPaged(repo ProductRepository, t *testing.T) {

	dupes := make(map[string]*store.Product)
	var products []*store.Product
	var err error
	last := ""
	lastPrice := float64(0)
	total := int64(0)
	visited := int64(0)

	for {
		products, total, err = repo.GetPaged(categories[0], 25, last, lastPrice)
		if err != nil {
			t.Fatal(err)
		}

		for i, p := range products {
			if i > 0 {
				if products[i-1].Price < products[i].Price {
					t.Error("products not sorted")
					return
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

func testGetProduct(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Errorf("could not get products %v", err)
		return
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
	}

	for _, v := range ps {
		p, err := repo.Get(v.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if p.ID != v.ID {
			t.Fatalf("wrong product returned expected %s got %s", v.ID, p.ID)
		}
	}
}

func testUpsertProduct(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Errorf("could not get all %v", err)
		return
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
	}

	p, err := repo.Get(ps[0].ID)
	if err != nil {
		t.Errorf("could not get product %v", err)
		t.Error(err)
		return
	}

	p.Name = p.Name + " Updated"

	if err := repo.Upsert(categories[0], p); err != nil {
		t.Errorf("could not upsert product %v", err)
		return
	}

	n, err := repo.Get(p.ID)
	if err != nil {
		t.Errorf("could not get updated product %v", err)
		return
	}

	if n.Name != p.Name {
		t.Errorf("product name not updated got %s expected %s", n.Name, p.Name)
	}
}

func testDeleteProduct(repo ProductRepository, t *testing.T) {
	ps, _, err := repo.GetPaged(categories[0], 25, "", float64(0))
	if err != nil {
		t.Errorf("could not get all products %v", err)
		return
	}

	if len(ps) == 0 {
		t.Fatal("no products returned")
		return
	}

	if err := repo.Delete(ps[0].ID); err != nil {
		t.Errorf("could not delete product %s %v", ps[0].ID, err)
		return
	}

	p, err := repo.Get(ps[0].ID)
	if err != nil {
		t.Errorf("could not get product %s %v", ps[0].ID, err)
		return
	}

	if p != nil {
		t.Fatalf("product not deleted %s", p.ID)
	}
}
