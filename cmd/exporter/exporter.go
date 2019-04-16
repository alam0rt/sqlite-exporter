package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mattn/go-sqlite3"
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
	dbArg := flag.String("db", "", "a sqlite3 database")

	flag.Parse()
	database.Open("Algo.db")
	recordMetrics()
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	fmt.Print(&buf)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+*portArg, nil)
}
