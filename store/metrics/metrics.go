package metrics

import (
	cmetrics "github.com/chadgrant/dynamodb-go-sample/store/metrics/category"
	pmetrics "github.com/chadgrant/dynamodb-go-sample/store/metrics/product"
	"github.com/chadgrant/dynamodb-go-sample/store/metrics/util"

	kit "github.com/chadgrant/kit/metrics"
	"github.com/chadgrant/kit/metrics/provider"
)

var appMetrics = &Metrics{}

type Metrics struct {
	TotalErrors kit.Counter
	Category    cmetrics.Category
	Product     pmetrics.Product
}

func New() *Metrics {
	// intentional singleton
	return appMetrics
}

func init() {
	providers := []provider.Provider{
		provider.NewExpvarProvider(),
		provider.NewPrometheusProvider("sample", "service"),
	}

	appMetrics.TotalErrors = util.Counters("errors", providers...)
	cmetrics.Build(appMetrics.TotalErrors, &appMetrics.Category, providers...)
	pmetrics.Build(appMetrics.TotalErrors, &appMetrics.Product, providers...)
}
