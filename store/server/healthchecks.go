package server

import (
	"github.com/chadgrant/dynamodb-go-sample/store/repository/dynamo"
	"github.com/chadgrant/go-http-infra/infra/health"
)

func RegisterHealthChecks(cfg *Configuration) health.HealthChecker {
	h := cfg.HealthChecks
	hc := health.NewHealthChecker()

	if h.Dynamo.Enabled {
		dyn := dynamo.New(cfg.AWS.Region, cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.Dynamo.Endpoint)
		hc.AddReadiness("dynamo", h.Dynamo.Interval, dynamo.Health(dyn, h.Dynamo.Timeout, cfg.Dynamo.Tables.Categories, cfg.Dynamo.Tables.Products))
	}

	if h.TCP.Enabled {
		hc.AddReadiness("google tcp connection", h.TCP.Interval, health.TCPDialCheck(h.TCP.Value, h.TCP.Timeout))
	}

	if h.HTTP.Enabled {
		hc.AddReadiness("http get", h.HTTP.Interval, health.HTTPGetCheck(h.HTTP.Value, h.HTTP.Timeout))
	}

	if h.DNS.Enabled {
		hc.AddReadiness("dns lookup", h.DNS.Interval, health.DNSResolveCheck(h.DNS.Value, h.DNS.Timeout))
	}

	return hc
}
