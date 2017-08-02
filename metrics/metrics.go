package metrics

import (
	"os"
)

var registry MetricsRegistry = nil

func init() {
	metricsProvider := os.Getenv("METRICS_PROVIDER")
	if metricsProvider == "" || metricsProvider == "prometheus" {
		if registry == nil {
			registry = prometheusMetrics{}
			registry.Start()
		}
	}
}


type MetricsRegistry interface {
	Start()
	IncStripeErrors()
	IncAuth0Errors()
	IncPageLoadErrors()
	IncRequests()
	IncConfigLoads()
	IncLoginAttempts()
}

func Registry() MetricsRegistry {
	return registry
}
