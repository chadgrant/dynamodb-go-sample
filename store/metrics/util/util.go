package util

import (
	kit "github.com/chadgrant/kit/metrics"
	"github.com/chadgrant/kit/metrics/multi"
	"github.com/chadgrant/kit/metrics/provider"
)

func Counters(name string, p ...provider.Provider) kit.Counter {
	counters := make([]kit.Counter, len(p))
	for i := range p {
		counters[i] = p[i].NewCounter(name)
	}

	return multi.NewCounter(counters...)
}

func Histograms(name string, buckets int, p ...provider.Provider) kit.Histogram {
	histograms := make([]kit.Histogram, len(p))
	for i := range p {
		histograms[i] = p[i].NewHistogram(name, buckets)
	}

	return multi.NewHistogram(histograms...)
}
