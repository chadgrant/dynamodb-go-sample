package metrics

import (
	kit "github.com/chadgrant/kit/metrics"
	prov "github.com/chadgrant/kit/metrics/provider"
)

var AppMetrics = &Metrics{}

const appName = "dynamo-go-sample"

type Metrics struct {
	Repository Repository
}

type Repository struct {
	Counters RepositoryCounters
}

type RepositoryCounters struct {
	GetCateogries kit.Counter
}

func init() {
	p := prov.NewPrometheusProvider(appName, "repository")
	categoriesRepoMetrics(p)
}

func categoriesRepoMetrics(p prov.Provider) {
	AppMetrics.Repository.Counters.GetCateogries = p.NewCounter("categories_get")
}
