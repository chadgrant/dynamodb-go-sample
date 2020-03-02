package category

import (
	"github.com/chadgrant/dynamodb-go-sample/store/metrics/util"
	kit "github.com/chadgrant/kit/metrics"
	"github.com/chadgrant/kit/metrics/multi"
	"github.com/chadgrant/kit/metrics/provider"
)

type (
	Category struct {
		Repository RepoMetrics
		Service    ServiceMetrics
	}

	ServiceMetrics struct {
		Counters   Counters
		Histograms Histograms
		Errors     Counters
	}

	RepoMetrics struct {
		Counters   RepoCounters
		Histograms RepoHistograms
		Errors     RepoCounters
	}

	Counters struct {
		Get    kit.Counter
		Add    kit.Counter
		Update kit.Counter
		Delete kit.Counter
	}

	RepoCounters struct {
		Get    kit.Counter
		Upsert kit.Counter
		Delete kit.Counter
	}

	Histograms struct {
		Get    kit.Histogram
		Add    kit.Histogram
		Update kit.Histogram
		Delete kit.Histogram
	}

	RepoHistograms struct {
		Get    kit.Histogram
		Upsert kit.Histogram
		Delete kit.Histogram
	}
)

func Build(errors kit.Counter, cat *Category, p ...provider.Provider) {
	total := multi.NewCounter(errors, counters("category_errors", p...))

	cat.Service.Counters = Counters{
		Get:    counters("category_get", p...),
		Add:    counters("category_add", p...),
		Update: counters("category_update", p...),
		Delete: counters("category_delete", p...),
	}

	cat.Service.Histograms = Histograms{
		Get:    histograms("category_get_dur", 10, p...),
		Add:    histograms("category_add_dur", 10, p...),
		Update: histograms("category_update_dur", 10, p...),
		Delete: histograms("category_delete_dur", 10, p...),
	}

	serviceErrs := multi.NewCounter(total, counters("category_service_errors", p...))
	cat.Service.Errors = Counters{
		Get:    multi.NewCounter(counters("category_get_error", p...), serviceErrs),
		Add:    multi.NewCounter(counters("category_add_error", p...), serviceErrs),
		Update: multi.NewCounter(counters("category_update_error", p...), serviceErrs),
		Delete: multi.NewCounter(counters("category_delete_error", p...), serviceErrs),
	}

	cat.Repository.Counters = RepoCounters{
		Get:    counters("category_repo_get", p...),
		Upsert: counters("category_repo_upsert", p...),
		Delete: counters("category_repo_delete", p...),
	}

	cat.Repository.Histograms = RepoHistograms{
		Get:    histograms("category_repo_get_dur", 10, p...),
		Upsert: histograms("category_repo_upsert_dur", 10, p...),
		Delete: histograms("category_repo_delete_dur", 10, p...),
	}

	repoErrs := multi.NewCounter(total, counters("category_repo_errors", p...))
	cat.Repository.Errors = RepoCounters{
		Get:    multi.NewCounter(counters("category_repo_get_errors", p...), repoErrs),
		Upsert: multi.NewCounter(counters("category_repo_upsert_errors", p...), repoErrs),
		Delete: multi.NewCounter(counters("category_repo_delete_errors", p...), repoErrs),
	}
}

func counters(name string, providers ...provider.Provider) kit.Counter {
	return util.Counters(name, providers...)
}

func histograms(name string, buckets int, providers ...provider.Provider) kit.Histogram {
	return util.Histograms(name, buckets, providers...)
}
