package repository

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/google/uuid"
)

type Populator struct {
	repo ProductRepository
}

func NewPopulator(repo ProductRepository) *Populator {
	return &Populator{
		repo,
	}
}

func (p *Populator) LoadProducts(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var prds []*store.Product
	if err := json.Unmarshal(bs, &prds); err != nil {
		return err
	}

	return p.addProducts(prds)
}

func (p *Populator) CreateProducts(max int) error {
	prds := make([]*store.Product, max)
	for i := 0; i < max; i++ {
		id, _ := uuid.NewRandom()
		pr := &store.Product{
			ID:          id.String(),
			Name:        "Test Product " + id.String(),
			Price:       randPrice(),
			Description: "You should buy this product, it's awesome. I have 2. You'll love it. Trust me.'",
		}
		prds[i] = pr
	}
	return p.addProducts(prds)
}

func (p *Populator) ExportProducts(path string) error {
	prods, err := p.repo.GetAll()
	if err != nil {
		return err
	}
	bs, err := json.Marshal(prods)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, bs, 0644); err != nil {
		return err
	}
	return nil
}

func (p *Populator) addProducts(products []*store.Product) error {
	for _, pr := range products {
		if err := p.repo.Add(pr); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randPrice() float64 {
	var min float64 = 0.99
	var max float64 = 999.99
	ret := min + rand.Float64()*(max-min)
	return float64(int(ret*100)) / 100
}
