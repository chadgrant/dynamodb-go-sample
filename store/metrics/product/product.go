package product

import (
	"github.com/chadgrant/dynamodb-go-sample/store/metrics/util"

	kit "github.com/chadgrant/kit/metrics"
	"github.com/chadgrant/kit/metrics/multi"
	"github.com/chadgrant/kit/metrics/provider"
)

type (
	Product struct {
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

	RepoCounters struct {
		Paged  kit.Counter
		Get    kit.Counter
		Upsert kit.Counter
		Delete kit.Counter
	}

	Counters struct {
		Paged  kit.Counter
		Get    kit.Counter
		Add    kit.Counter
		Update kit.Counter
		Delete kit.Counter
	}

	RepoHistograms struct {
		Paged  kit.Histogram
		Get    kit.Histogram
		Upsert kit.Histogram
		Delete kit.Histogram
	}

	Histograms struct {
		Paged  kit.Histogram
		Get    kit.Histogram
		Add    kit.Histogram
		Update kit.Histogram
		Delete kit.Histogram
	}
)

func Build(errors kit.Counter, prod *Product, p ...provider.Provider) {
	total := multi.NewCounter(errors, counters("product_errors", p...))

	prod.Service.Counters = Counters{
		Paged:  counters("product_paged", p...),
		Get:    counters("product_get", p...),
		Add:    counters("product_add", p...),
		Update: counters("product_update", p...),
		Delete: counters("product_delete", p...),
	}

	prod.Service.Histograms = Histograms{
		Paged:  histograms("product_paged_dur", 10, p...),
		Get:    histograms("product_get_dur", 10, p...),
		Add:    histograms("product_add_dur", 10, p...),
		Update: histograms("product_update_dur", 10, p...),
		Delete: histograms("product_delete_dur", 10, p...),
	}

	serviceErrs := multi.NewCounter(total, counters("product_service_errors", p...))
	prod.Service.Errors = Counters{
		Paged:  multi.NewCounter(counters("product_paged_error", p...), serviceErrs),
		Get:    multi.NewCounter(counters("product_get_error", p...), serviceErrs),
		Add:    multi.NewCounter(counters("product_add_error", p...), serviceErrs),
		Update: multi.NewCounter(counters("product_update_error", p...), serviceErrs),
		Delete: multi.NewCounter(counters("product_delete_error", p...), serviceErrs),
	}

	prod.Repository.Counters = RepoCounters{
		Paged:  counters("product_repo_paged", p...),
		Get:    counters("product_repo_get", p...),
		Upsert: counters("product_repo_upsert", p...),
		Delete: counters("product_repo_delete", p...),
	}

	prod.Repository.Histograms = RepoHistograms{
		Paged:  histograms("product_repo_paged_dur", 10, p...),
		Get:    histograms("product_repo_get_dur", 10, p...),
		Upsert: histograms("product_repo_upsert_dur", 10, p...),
		Delete: histograms("product_repo_delete_dur", 10, p...),
	}

	repoErrs := multi.NewCounter(total, counters("product_repo_errors", p...))
	prod.Repository.Errors = RepoCounters{
		Paged:  multi.NewCounter(counters("product_repo_paged_errors", p...), repoErrs),
		Get:    multi.NewCounter(counters("product_repo_get_errors", p...), repoErrs),
		Upsert: multi.NewCounter(counters("product_repo_upsert_errors", p...), repoErrs),
		Delete: multi.NewCounter(counters("product_repo_delete_errors", p...), repoErrs),
	}
}

func counters(name string, providers ...provider.Provider) kit.Counter {
	return util.Counters(name, providers...)
}

func histograms(name string, buckets int, providers ...provider.Provider) kit.Histogram {
	return util.Histograms(name, buckets, providers...)
}
