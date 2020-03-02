package service

import (
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store"
	"github.com/chadgrant/dynamodb-go-sample/store/metrics"
)

type serviceMetrics struct {
	*metrics.Metrics
	svc Service
}

func NewServiceMetrics(m *metrics.Metrics, service Service) Service {
	return &serviceMetrics{m, service}
}

func (m *serviceMetrics) Categories() ([]string, error) {
	start := time.Now()
	defer m.Category.Service.Histograms.Get.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Category.Service.Counters.Get.Add(1)

	cats, err := m.svc.Categories()
	if err != nil {
		m.Category.Service.Errors.Get.Add(1)
	}
	return cats, err
}

func (m *serviceMetrics) ProductsPaged(category string, limit int, last string, lastprice float64) ([]*store.Product, error) {
	start := time.Now()
	defer m.Product.Service.Histograms.Paged.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Product.Service.Counters.Paged.Add(1)

	prods, err := m.svc.ProductsPaged(category, limit, last, lastprice)
	if err != nil {
		m.Product.Service.Errors.Paged.Add(1)
	}
	return prods, err
}

func (m *serviceMetrics) ProductById(productID string) (*store.Product, error) {
	start := time.Now()
	defer m.Product.Service.Histograms.Get.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Product.Service.Counters.Get.Add(1)

	prod, err := m.svc.ProductById(productID)
	if err != nil {
		m.Product.Service.Errors.Get.Add(1)
	}
	return prod, err
}

func (m *serviceMetrics) AddProduct(p *store.Product) error {
	start := time.Now()
	defer m.Product.Service.Histograms.Add.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Product.Service.Counters.Add.Add(1)

	if err := m.svc.AddProduct(p); err != nil {
		m.Product.Service.Errors.Add.Add(1)
		return err
	}
	return nil
}

func (m *serviceMetrics) UpdateProduct(p *store.Product) error {
	start := time.Now()
	defer m.Product.Service.Histograms.Update.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Product.Service.Counters.Update.Add(1)

	if err := m.svc.UpdateProduct(p); err != nil {
		m.Product.Service.Errors.Update.Add(1)
		return err
	}
	return nil
}

func (m *serviceMetrics) DeleteProduct(productID string) error {
	start := time.Now()
	defer m.Product.Service.Histograms.Delete.Observe(float64(time.Now().Sub(start).Milliseconds()))
	m.Product.Service.Counters.Delete.Add(1)

	if err := m.svc.DeleteProduct(productID); err != nil {
		m.Product.Service.Errors.Add.Add(1)
		return err
	}
	return nil
}
