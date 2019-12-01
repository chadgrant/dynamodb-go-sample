package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
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
		return fmt.Errorf("reading products from disk %s %v", path, err)
	}

	var prds []*store.Product
	if err := json.Unmarshal(bs, &prds); err != nil {
		return fmt.Errorf("unmarshaling products from disk %s, %v", path, err)
	}

	if err := p.add(prds); err != nil {
		return fmt.Errorf("saving products %v", err)
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
				Category:    strings.ToLower(c),
				Name:        fmt.Sprintf("Test %s %s", c, id.String()),
				Price:       randPrice(),
				Description: fmt.Sprintf("You should buy this %s, it's awesome. I have 2. You'll love it. Trust me.", c),
			}
		}
		if err := p.add(prds); err != nil {
			return fmt.Errorf("saving products %v", err)
		}
	}
	return nil
}

func (p *Populator) Export(path string) error {
	all := make([]*store.Product, 0)

	for _, c := range categories {
		last := ""
		lastPrice := float64(0)
		for {
			products, _, err := p.repo.GetPaged(c, 25, last, lastPrice)
			if err != nil {
				return err
			}

			all = append(all, products...)

			if len(products) < 25 {
				break
			}

			last = products[len(products)-1].ID
			lastPrice = products[len(products)-1].Price
		}
	}

	bs, err := json.Marshal(all)
	if err != nil {
		return fmt.Errorf("export marshaling products %v", err)
	}
	if err := ioutil.WriteFile(path, bs, 0644); err != nil {
		return fmt.Errorf("export saving products %v", err)
	}

	return nil
}

func (p *Populator) add(products []*store.Product) error {
	for _, pr := range products {
		if err := p.repo.Upsert(pr); err != nil {
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
