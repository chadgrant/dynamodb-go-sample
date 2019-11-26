package mock

import (
	"github.com/chadgrant/dynamodb-go-sample/store"
)

type MockRepository struct {
	Products map[string]*store.Product
}

func NewProductRepository() *MockRepository {
	repo := &MockRepository{
		Products: make(map[string]*store.Product),
	}

	return repo
}

func (r *MockRepository) GetAll() ([]*store.Product, error) {
	p := make([]*store.Product, len(r.Products))
	i := 0
	for _, v := range r.Products {
		p[i] = v
		i++
	}
	return p, nil
}

func (r *MockRepository) GetPaged(start string, limit int) ([]*store.Product, int, error) {
	prds := make([]*store.Product, 0)
	collect := len(start) == 0

	for _, v := range r.Products {
		if !collect {
			if v.ID == start {
				collect = true
			}
			continue
		}

		prds = append(prds, v)

		if len(prds) == limit {
			break
		}
	}
	return prds, len(r.Products), nil
}

func (r *MockRepository) Get(productID string) (*store.Product, error) {
	return r.Products[productID], nil
}

func (r *MockRepository) Add(product *store.Product) error {
	r.Products[product.ID] = product
	return nil
}

func (r *MockRepository) Upsert(product *store.Product) error {
	r.Products[product.ID] = product
	return nil
}

func (r *MockRepository) Delete(productID string) error {
	delete(r.Products, productID)
	return nil
}
