package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/google/uuid"
)

//mocked/faked
var categories = []string{"Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory"}

type Populator struct {
	repo ProductRepository
}

func NewPopulator(repo ProductRepository) *Populator {
	return &Populator{
		repo,
	}
}

func (p *Populator) Load(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var prds map[string][]*store.Product
	if err := json.Unmarshal(bs, &prds); err != nil {
		return err
	}

	for c, a := range prds {
		if err := p.add(c, a); err != nil {
			return err
		}
	}
	return nil
}

func (p *Populator) Create(max int) error {
	for _, c := range categories {
		prds := make([]*store.Product, max)
		for i := 0; i < max; i++ {
			id, _ := uuid.NewRandom()
			prds[i] = &store.Product{
				ID:          id.String(),
				Name:        fmt.Sprintf("Test %s %s", c, id.String()),
				Price:       randPrice(),
				Description: fmt.Sprintf("You should buy this %s, it's awesome. I have 2. You'll love it. Trust me.", c),
			}
		}
		if err := p.add(c, prds); err != nil {
			return err
		}
	}
	return nil
}

func (p *Populator) Export(path string) error {
	all := make(map[string][]*store.Product)

	for _, c := range categories {
		all[c] = make([]*store.Product, 0)
		last := ""
		lastPrice := float64(0)
		for {
			products, _, err := p.repo.GetPaged(c, 25, last, lastPrice)
			if err != nil {
				return err
			}

			all[c] = append(all[c], products...)

			if len(products) < 25 {
				break
			}

			last = products[len(products)-1].ID
			lastPrice = products[len(products)-1].Price
		}
	}

	bs, err := json.Marshal(all)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, bs, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Populator) add(category string, products []*store.Product) error {
	for _, pr := range products {
		if err := p.repo.Upsert(category, pr); err != nil {
			return fmt.Errorf("upsert failed: %v", err)
		}
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randPrice() float64 {
	min, max := 0.99, 999.99
	r := min + rand.Float64()*(max-min)
	return float64(int(r*100)) / 100
}
