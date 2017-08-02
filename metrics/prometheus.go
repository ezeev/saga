package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	stripeErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
		Name: "stripe_errors_total",
		Help: "Number of stripe errors",
	})
	auth0Errors = prometheus.NewCounter(
		prometheus.CounterOpts{
		Name: "auth0_errors_total",
		Help: "Number of auth0 errors",
	})
	pageLoadErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
		Name: "page_errors_total",
		Help: "Number of auth0 errors",
	})
	totalErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
		Name: "errors_total",
		Help: "Number of total errors",
	})
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Number of total requests",
	})
	totalConfigLoads = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "config_loads_total",
			Help: "Number of total config loads",
		})
	totalLoginAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Number of total login attemps",
		})
)



type prometheusMetrics struct {
}



func (p prometheusMetrics) Start() {

	prometheus.MustRegister(stripeErrors)
	prometheus.MustRegister(auth0Errors)
	prometheus.MustRegister(totalErrors)
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(totalConfigLoads)
	prometheus.MustRegister(totalLoginAttempts)
	uri := os.Getenv("METRICS_URI")
	log.Printf("Registering Metrics URI: %s", uri)
	http.Handle(uri, promhttp.Handler())

}

func (p prometheusMetrics) IncStripeErrors() {
	stripeErrors.Inc()
	totalErrors.Inc()
}

func (p prometheusMetrics) IncAuth0Errors() {
	auth0Errors.Inc()
	totalErrors.Inc()
}

func (p prometheusMetrics) IncPageLoadErrors() {
	pageLoadErrors.Inc()
	totalErrors.Inc()
}

func (p prometheusMetrics) IncRequests() {
	totalRequests.Inc()
}

func (p prometheusMetrics) IncConfigLoads() {
	totalConfigLoads.Inc()
}

func (p prometheusMetrics) IncLoginAttempts() {
	totalLoginAttempts.Inc()
}