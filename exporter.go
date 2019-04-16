package exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	Metrics []prometheus.Gauge // contains addresses of successful CreateMetric()
)

func CreateMetric(name string, help string) prometheus.Gauge {
	metric := promauto.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	Metrics = append(Metrics, metric) // append new gauge address to Metrics
	return metric
}

func RecordMetrics() {
	go func() {
		for {
			// iterate over Metrics slice and increase()
			for _, metric := range Metrics {
				go metric.Inc()
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func Listen(port string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("IMPLEMENT")
	}

}
