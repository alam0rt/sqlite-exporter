package exporter

import (
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	MetricsMap map[string]*Metric = make(map[string]*Metric)
)

type Metric struct {
	gauge prometheus.Gauge // holds our gauge objects
	Name  string           // name of metric
	help  string           // holds the metric description
	Value float64          // holds last obtained query result
	Query string           // holds the query to run against target db
}

// CreateMetric builds a Metric struct and pushes it into a global map
func CreateMetric(name string, help string, query string) Metric {
	g := promauto.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	m := Metric{
		g,
		name,
		help,
		0, // set metric value to 0
		query,
	}
	// now we push this metric to the global slice 'Metric'
	//	Metrics = append(Metrics, m)
	MetricsMap[name] = &m
	return m
}

// UpdateMetric takes a given pointer to a Metric struct
// and passes the new value of the metric to Prometheus so
// that the metric is updated.
func UpdateMetric(m *Metric) {
	//	fmt.Printf("UpdateMetric [%s] => %f (%s)\n", m.Name, m.Value, m.Query)
	m.gauge.Set(m.Value)
}

// SetMetric takes a name of a metric that has been created and exists
// in the MetricsMap map and sets its value to whatever float is provided.
func SetMetric(name string, v float64) {
	MetricsMap[name].Value = v
	//	fmt.Printf("SetMetric [%s] => %f\n", name, v)
}

// Listen takes a port number as a string and sets up a "/metrics" handler
// which everything is exported to.
func Listen(port string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logging.Error.Printf("unable to listen on %i: %s\n", port, err)
	}

}
