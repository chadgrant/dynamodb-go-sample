package repository

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/google/uuid"
)

type MockRepository struct {
	products   []*store.Product
	lookup     map[string]string
	categories CategoryRepository
}

func NewMockProductRepository(repo CategoryRepository, max int) *MockRepository {
	m := &MockRepository{
		products:   make([]*store.Product, 0),
		lookup:     make(map[string]string),
		categories: repo,
	}
	m.create(max)
	return m
}

func (r *MockRepository) GetPaged(category string, limit int, lastID string, lastPrice float64) ([]*store.Product, error) {

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

	return prds[start:end], nil
}

func (r *MockRepository) Get(productID string) (*store.Product, error) {
	_, p := find(r.products, productID)
	return p, nil
}

func (r *MockRepository) Upsert(product *store.Product) error {
	i, _ := find(r.products, product.ID)
	if i >= 0 {
		r.products[i] = product
	} else {
		r.products = append(r.products, product)
	}
	r.lookup[product.ID] = product.Category
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

func (r *MockRepository) create(max int) error {
	cats, err := r.categories.GetAll()
	if err != nil {
		return err
	}
	for _, c := range cats {
		for i := 0; i < max; i++ {
			id, _ := uuid.NewRandom()
			p := &store.Product{
				ID:          id.String(),
				Category:    strings.ToLower(c),
				Name:        fmt.Sprintf("Test %s %s", c, id.String()),
				Price:       randPrice(),
				Description: fmt.Sprintf("You should buy this %s, it's awesome. I have 2. You'll love it. Trust me.", c),
			}
			if err := r.Upsert(p); err != nil {
				return fmt.Errorf("saving products %v", err)
			}
		}

	}
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randPrice() float64 {
	min, max := 0.99, 999.99
	r := min + rand.Float64()*(max-min)
	return float64(int(r*100)) / 100
}
