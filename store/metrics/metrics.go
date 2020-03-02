package metrics

import (
	cmetrics "github.com/chadgrant/dynamodb-go-sample/store/metrics/category"
	pmetrics "github.com/chadgrant/dynamodb-go-sample/store/metrics/product"
	"github.com/chadgrant/dynamodb-go-sample/store/metrics/util"

	kit "github.com/chadgrant/kit/metrics"
	"github.com/chadgrant/kit/metrics/provider"
)

var AppMetrics = &Metrics{}

const appName = "dynamo_go_sample"

type (
	Metrics struct {
		TotalErrors kit.Counter
		Category    cmetrics.Category
		Product     pmetrics.Product
	}
)

func init() {
	providers := []provider.Provider{provider.NewExpvarProvider(), provider.NewPrometheusProvider(appName, "category")}

	AppMetrics.TotalErrors = util.Counters("errors", providers...)
	cmetrics.Build(AppMetrics.TotalErrors, &AppMetrics.Category, providers...)
	pmetrics.Build(AppMetrics.TotalErrors, &AppMetrics.Product, providers...)
}
