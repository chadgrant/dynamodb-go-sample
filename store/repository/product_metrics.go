package repository

import (
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/metrics/product"
)

type metricsProductRepository struct {
	*product.RepoMetrics
	repo ProductRepository
}

func NewMetricsProductRepository(m *product.RepoMetrics, r ProductRepository) ProductRepository {
	return &metricsProductRepository{m, r}
}

func (r *metricsProductRepository) GetPaged(category string, limit int, last string, lastprice float64) ([]*store.Product, error) {
	start := time.Now()
	defer r.Histograms.Paged.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Paged.Add(1)

	prods, err := r.repo.GetPaged(category, limit, last, lastprice)
	if err != nil {
		r.Errors.Paged.Add(1)
	}
	return prods, err
}

func (r *metricsProductRepository) Get(productID string) (*store.Product, error) {
	start := time.Now()
	defer r.Histograms.Get.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Get.Add(1)

	prod, err := r.repo.Get(productID)
	if err != nil {
		r.Errors.Get.Add(1)
	}
	return prod, err
}

func (r *metricsProductRepository) Upsert(product *store.Product) error {
	start := time.Now()
	defer r.Histograms.Upsert.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Upsert.Add(1)

	err := r.repo.Upsert(product)
	if err != nil {
		r.Errors.Upsert.Add(1)
	}
	return err
}

func (r *metricsProductRepository) Delete(productID string) error {
	start := time.Now()
	defer r.Histograms.Delete.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Delete.Add(1)

	err := r.repo.Delete(productID)
	if err != nil {
		r.Errors.Delete.Add(1)
	}
	return err
}
