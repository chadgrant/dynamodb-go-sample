package repository

import (
	"time"

	"github.com/chadgrant/dynamodb-go-sample/store/metrics/category"
)

type metricsCategoryRepository struct {
	*category.RepoMetrics
	repo CategoryRepository
}

func NewMetricsCategoryRepository(m *category.RepoMetrics, r CategoryRepository) CategoryRepository {
	return &metricsCategoryRepository{m, r}
}

func (r *metricsCategoryRepository) GetAll() ([]string, error) {
	start := time.Now()
	defer r.Histograms.Get.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Get.Add(1)

	cats, err := r.repo.GetAll()
	if err != nil {
		r.Errors.Get.Add(1)
	}
	return cats, err
}

func (r *metricsCategoryRepository) Upsert(category string) error {
	start := time.Now()
	defer r.Histograms.Upsert.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Upsert.Add(1)

	err := r.repo.Upsert(category)
	if err != nil {
		r.Errors.Upsert.Add(1)
	}
	return err
}

func (r *metricsCategoryRepository) Delete(category string) error {
	start := time.Now()
	defer r.Histograms.Delete.Observe(float64(time.Now().Sub(start).Milliseconds()))
	r.Counters.Delete.Add(1)

	err := r.repo.Delete(category)
	if err != nil {
		r.Errors.Delete.Add(1)
	}
	return err
}
