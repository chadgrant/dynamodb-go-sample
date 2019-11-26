package repository

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

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

	repo := dynamo.NewProductRepository("products", sess, config)
	repo.DeleteTable()
	if err := repo.CreateTable(); err != nil {
		t.Errorf("could not create table %v", err)
		return
	}
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

	t.Run("GetProducts", func(t *testing.T) {
		testGetProducts(repo, t)
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

	if err := repo.Add(p); err != nil {
		t.Error(err)
	}
}

func testGetProducts(repo ProductRepository, t *testing.T) {
	p, err := repo.GetAll()
	if err != nil {
		t.Error(err)
		return
	}
	if len(p) == 0 {
		t.Fatal("no products returned")
	}
}

func testGetProductsPaged(repo ProductRepository, t *testing.T) {
	p, tot, err := repo.GetPaged("", 25)
	if err != nil {
		t.Error(err)
		return
	}
	if len(p) == 0 {
		t.Fatal("no products returned")
	}
	if tot <= 0 {
		t.Error("total is not correct")
	}

	if len(p) != 25 {
		t.Errorf("unexpected amount returned. expected 25, got %d", len(p))
	}

	next := p[len(p)-1].ID
	p, tot, err = repo.GetPaged(next, 25)
	if err != nil {
		t.Error(err)
		return
	}
	if len(p) == 0 {
		t.Fatal("no products returned")
	}
	if tot <= 0 {
		t.Error("total is not correct")
	}

	if len(p) != 25 {
		t.Errorf("unexpected amount returned. expected 25, got %d", len(p))
	}

	for _, v := range p {
		if v.ID == next {
			t.Errorf("shouldn't see next %s", next)
		}
	}
}

func testGetProduct(repo ProductRepository, t *testing.T) {
	p, err := repo.GetAll()
	if err != nil {
		t.Errorf("could not get all products %v", err)
		return
	}

	if len(p) == 0 {
		t.Fatal("no products returned")
	}

	for _, v := range p {
		sp, err := repo.Get(v.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if sp.ID != v.ID {
			t.Fatalf("wrong product returned expected %s got %s", v.ID, sp.ID)
		}
	}
}

func testUpsertProduct(repo ProductRepository, t *testing.T) {
	ps, err := repo.GetAll()
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

	if err := repo.Upsert(p); err != nil {
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
	ps, err := repo.GetAll()
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
