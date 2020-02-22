package repository_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/chadgrant/dynamodb-go-sample/store/repository/mock"
)

func TestMock(t *testing.T) {
	c := mock.NewCategoryRepository("Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory")
	runTests(c, mock.NewProductRepository(c, 100), t)
}

func TestDynamoIntegration(t *testing.T) {
	if len(os.Getenv("TEST_INTEGRATION")) == 0 {
		t.Log("Skipping integration tests, TEST_INTEGRATION environment variable not set")
		return
	}

	ep := os.Getenv("DYNAMO_ENDPOINT")
	if len(ep) == 0 {
		ep = "http://localhost:8000"
	}

	dyn := dynamo.New("us-east-1", "key", "secret", ep)

	if err := dynamo.WaitForTables(dyn, time.Second*30, "products"); err != nil {
		t.Fatalf("waiting on dynamodb %v", err)
	}

	runTests(dynamo.NewCategoryRepository("categories", dyn), dynamo.NewProductRepository("products", dyn), t)
}

func runTests(catRepo repository.CategoryRepository, repo repository.ProductRepository, t *testing.T) {

	cats, err := catRepo.GetAll()
	if err != nil {
		t.Fatalf("could not get categories %v", err)
	}

	t.Run("GetSingle", func(t *testing.T) {
		testGet(cats, repo, t)
	})

	t.Run("GetPaged", func(t *testing.T) {
		testGetPaged(cats, repo, t)
	})

	t.Run("Upsert", func(t *testing.T) {
		testUpsert(cats, repo, t)
	})

	t.Run("Add", func(t *testing.T) {
		testAdd(cats, repo, t)
	})

	t.Run("Delete", func(t *testing.T) {
		testDelete(cats, repo, t)
	})
}

func testAdd(categories []string, repo repository.ProductRepository, t *testing.T) {
	id, _ := uuid.NewRandom()
	p := &store.Product{
		ID:       id.String(),
		Category: strings.ToLower(categories[0]),
		Name:     "Test Product " + id.String(),
		Price:    1.00,
	}

	if err := repo.Upsert(p); err != nil {
		t.Error(err)
	}
}

func testGetPaged(categories []string, repo repository.ProductRepository, t *testing.T) {
	dupes := make(map[string]*store.Product)
	var all []*store.Product
	last := ""
	lastPrice := float64(0)
	size := 25
	pages := 0

	for {
		products, err := repo.GetPaged(categories[0], size, last, lastPrice)
		if err != nil {
			t.Fatalf("get paged %v", err)
		}
		all = append(all, products...)
		for i, p := range products {
			if i > 0 {
				if products[i-1].Price < products[i].Price {
					t.Fatal("products not sorted")
				}
			}
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
		pages++
	}

	if len(all) == 0 {
		t.Errorf("did not return products")
	}

	if pages < 3 {
		t.Errorf("did not get enough pages %d", pages)
	}
}

func testGet(categories []string, repo repository.ProductRepository, t *testing.T) {
	ps, err := repo.GetPaged(categories[0], 25, "", float64(0))
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

func testUpsert(categories []string, repo repository.ProductRepository, t *testing.T) {
	ps, err := repo.GetPaged(categories[0], 25, "", float64(0))
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

	if err := repo.Upsert(p); err != nil {
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

func testDelete(categories []string, repo repository.ProductRepository, t *testing.T) {
	ps, err := repo.GetPaged(categories[0], 25, "", float64(0))
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
