package main

import (
	"bitbucket.org/dragontailcom/sqlite-exporter/internal/config"
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/database"
	"bytes"
	"database/sql"
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
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

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
	c := config.ProcessConfig(*configArg)
	fmt.Println(c)
	// open db
	db, err := sql.Open("sqlite3", *dbArg)
	if err != nil {
		logger.Fatal("Unable to open database")
	}
	recordMetrics()
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	value := database.QueryMetric(db, *dbArg, "SELECT random() as Metric")
	fmt.Print(value)
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":"+*portArg, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
