package service

import (
	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/repository"
)

type Service interface {
	Categories() ([]string, error)

	ProductsPaged(category string, limit int, last string, lastprice float64) ([]*store.Product, error)
	ProductById(productID string) (*store.Product, error)
	AddProduct(p *store.Product) error
	UpdateProduct(p *store.Product) error
	DeleteProduct(productID string) error
}

type service struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

func NewService(cat repository.CategoryRepository, prod repository.ProductRepository) Service {
	return &service{cat, prod}
}

func (s *service) Categories() ([]string, error) {
	return s.categoryRepo.GetAll()
}

func (s *service) ProductsPaged(category string, limit int, last string, lastprice float64) ([]*store.Product, error) {
	return s.productRepo.GetPaged(category, limit, last, lastprice)
}

func (s *service) ProductById(productID string) (*store.Product, error) {
	return s.productRepo.Get(productID)
}

func (s *service) AddProduct(p *store.Product) error {
	// p.Created = time.Now().UTC()
	return s.productRepo.Upsert(p)
}

func (s *service) UpdateProduct(p *store.Product) error {
	// p.Updated = time.Now().UTC()
	return s.productRepo.Upsert(p)
}

func (s *service) DeleteProduct(productID string) error {
	return s.productRepo.Delete(productID)
}
