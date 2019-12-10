package repository

import "github.com/chadgrant/dynamodb-go-sample/store"

type ProductRepository interface {
	GetPaged(category string, limit int, last string, lastprice float64) ([]*store.Product, error)
	Get(productID string) (*store.Product, error)
	Upsert(product *store.Product) error
	Delete(productID string) error
}
