package metrics

import (
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Map of all the registered Counters.
var registeredCounters = make(map[string]prometheus.Counter)

// Valid regex for prometheus model name.
// (Prometheus model reference: https://github.com/prometheus/common)
const validPrometheusName = "[^a-zA-Z_][^a-zA-Z0-9_]*"

// GetCounter return the prometheus counter of given name.
func GetCounter(name string) prometheus.Counter {
	// Validate and fix name (replace invalid chars with underscore "_").
	name = validateAndFix(name, validPrometheusName)
	return registeredCounters[name]
}

// PrometheusCounter create a new Counter of given name with help text.
func RegisterPrometheusCounter(name, helpText string) {
	// Validate and fix name (replace invalid chars with underscore "_").
	name = validateAndFix(name, validPrometheusName)

	// Register counter in the map.
	registeredCounters[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: helpText,
	})
	err = prometheus.Register(registeredCounters[name])
	if err != nil {
		log.Warningf("Error in metrics RegisterPrometheusCounter: %v", err)
	}
}

// StartPrometheusMetricServer starts a metric server handling
// "/metrics" on "localhost:8081" exposing the registered metrics.
func StartPrometheusMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Warningf("monitoring endpoint non available: %v: ", err)
	}
}

// Validate and fix name (replace invalid chars with underscore "_").
func validateAndFix(name, regex string) string {
	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Warningf("Error in metrics regex RegisterPrometheusCounter: %v", err)
	}
	name = reg.ReplaceAllString(name, "_")

	return name
}
