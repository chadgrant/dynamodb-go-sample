package repository

import (
	"sort"
	"strings"

	"github.com/chadgrant/dynamodb-go-sample/store"
)

type MockRepository struct {
	products []*store.Product
	lookup   map[string]string
}

func NewMockProductRepository() *MockRepository {
	return &MockRepository{
		products: make([]*store.Product, 0),
		lookup:   make(map[string]string),
	}
}

func (r *MockRepository) GetPaged(category string, limit int, lastID string, lastPrice float64) ([]*store.Product, int64, error) {

	prds := make([]*store.Product, 0)
	for _, p := range r.products {
		if strings.EqualFold(r.lookup[p.ID], category) {
			prds = append(prds, p)
		}
	}

	start := 0
	if len(lastID) > 0 {
		start, _ = find(prds, lastID)
		start++
	}

	end := start + limit
	if end > len(prds) {
		end = len(prds)
	}

	return prds[start:end], int64(len(prds)), nil
}

func (r *MockRepository) Get(productID string) (*store.Product, error) {
	_, p := find(r.products, productID)
	return p, nil
}

func (r *MockRepository) Upsert(category string, product *store.Product) error {
	i, _ := find(r.products, product.ID)
	if i >= 0 {
		r.products[i] = product
	} else {
		r.products = append(r.products, product)
	}
	r.lookup[product.ID] = category
	r.sort()
	return nil
}

func (r *MockRepository) Delete(productID string) error {
	ps := make([]*store.Product, 0)
	for _, p := range r.products {
		if p.ID != productID {
			ps = append(ps, p)
		}
	}
	r.products = ps
	return nil
}

func (r *MockRepository) sort() {
	sort.Slice(r.products, func(i, j int) bool {
		return r.products[i].Price > r.products[j].Price
	})
}

func find(prds []*store.Product, id string) (int, *store.Product) {
	for i := 0; i < len(prds); i++ {
		if prds[i].ID == id {
			return i, prds[i]
		}
	}
	return -1, nil
}