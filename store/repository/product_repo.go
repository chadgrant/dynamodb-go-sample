package repository

import "github.com/chadgrant/dynamodb-go-sample/store"

type ProductRepository interface {
	GetAll() ([]*store.Product, error)
	GetPaged(start string, limit int) ([]*store.Product, int, error)
	Get(productID string) (*store.Product, error)
	Add(product *store.Product) error
	Upsert(product *store.Product) error
	Delete(productID string) error
}
