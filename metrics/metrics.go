package metrics

import "os"

var registry MetricsRegistry

func init() {
	metricsProvider := os.Getenv("METRICS_PROVIDER")
	if metricsProvider == "" || metricsProvider == "prometheus" {
		registry = prometheusMetrics{}
		registry.Start()
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
