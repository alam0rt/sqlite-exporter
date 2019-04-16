package main

import (
	"bitbucket.org/dragontailcom/sqlite-exporter/internal/config"
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/database"
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			wow := float64(database.QueryMetric("Algo.db", "SELECT random() as Metric"))
			fmt.Println(wow)
			sqlite_random.Set(wow)
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	Random_gauge  int64 // test
	sqlite_random = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sqlite_random",
		Help: "Result of 'select random()'",
	})
)

func main() {
	// set up a log handler
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)
	logger.Print("starting sqlite-exporter...")

	// Set up argument parsing
	portArg := flag.String("port", "9001", "a port to listen on")
	dbArg := flag.String("database", "Algo.db", "a sqlite3 database")
	configArg := flag.String("config", "configuration.yml", "sqlite-exporter configuration file")

	flag.Parse()

	// DEBUG
	//Random_gauge = database.QueryMetric(*dbArg, "SELECT random() as Metric")
	c := config.ProcessConfig(*configArg)
	fmt.Println(c)
	//
	recordMetrics()
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	fmt.Print(&buf)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*portArg, nil)
}
