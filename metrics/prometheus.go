package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// PrometheusCounter create a new Counter of given name.
// It also starts a metric server on localhost:8081 for gathering purposes.
func PrometheusCounter(name, helpText string) prometheus.Counter {
	processedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: helpText,
	})
	err := prometheus.Register(processedCounter)
	if err != nil {
		log.Errorf("error in registering Prometheus handler: %v:", err)
	}

	go startPrometheusMetricsServer()

	log.Debug("init Prometheus()")

	return processedCounter
}

// startPrometheusMetricServer starts a metric server handling
// "/metrics" on "localhost:8081" exposing the registered metrics.
func startPrometheusMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Warningf("monitoring endpoint non available: %v: ", err)
	}
}
