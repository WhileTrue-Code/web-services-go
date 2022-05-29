package main

import (
	"fmt"
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
			Name: "my_app_http_hit_total",
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

	httpGetGroupsHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "getGroups_hits",
			Help: "getGroups_hits",
		},
	)

	httpPostGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "postGroup_hits",
			Help: "postGroup_hits",
		},
	)

	httpDeleteGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "deleteGroup_hits",
			Help: "deleteGroup_hits",
		},
	)

	httpPutGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "putGroup_hits",
			Help: "putGroup_hits",
		},
	)

	httpGetConfigFromGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "getConfigFromGroup_hits",
			Help: "getConfigFromGroup_hits",
		},
	)

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{httpHits, httpGetConfigHits, httpPostConfigHits, httpDeleteConfigHits,
		httpGetConfigsHits, httpGetGroupsHits, httpPostGroupHits, httpDeleteGroupHits, httpPutGroupHits, httpGetConfigFromGroupHits}

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
		if URL == "group" {
			fmt.Println(strings.Split(r.URL.String(), "/"))
			if len(strings.Split(r.URL.String(), "/")) > 3 && r.Method == "GET" {
				httpGetConfigFromGroupHits.Inc()
			} else {
				if r.Method == "GET" {
					httpGetGroupsHits.Inc()
				} else if r.Method == "POST" {
					httpPostGroupHits.Inc()
				} else if r.Method == "DELETE" {
					httpDeleteGroupHits.Inc()
				} else if r.Method == "PUT" {
					httpPutGroupHits.Inc()
				}
			}
		}
		httpHits.Inc()
		f(w, r) // original function call
	}
}
