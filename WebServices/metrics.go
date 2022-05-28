package main

import (
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Initial count.
	currentCount = 0

	// The Prometheus metric that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_hits",
			Help: "Total number of http hits.",
		},
	)

	httpGetConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "getConfig_hits",
			Help: "getConfig_hits",
		},
	)

	httpPostConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "postConfig_hits",
			Help: "postConfig_hits",
		},
	)

	httpDeleteConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "deleteConfig_hits",
			Help: "deleteConfig_hits",
		},
	)

	httpGetConfigsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "getConfigs_hits",
			Help: "getConfigs_hits",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{httpHits, httpGetConfigHits, httpPostConfigHits, httpDeleteConfigHits, httpGetConfigsHits}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		URL := strings.Split(r.URL.String(), "/")[1]
		if URL == "config" {
			if r.Method == "GET" {
				httpGetConfigHits.Inc()
			} else if r.Method == "POST" {
				httpPostConfigHits.Inc()
			} else if r.Method == "DELETE" {
				httpDeleteConfigHits.Inc()
			}
		}
		if URL == "configs" {
			if r.Method == "GET" {
				httpGetConfigsHits.Inc()
			}
		}
		httpHits.Inc()
		f(w, r) // original function call
	}
}
