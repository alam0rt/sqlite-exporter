package main

import (
	"bitbucket.org/dragontailcom/sqlite-exporter"
	"bitbucket.org/dragontailcom/sqlite-exporter/internal/config"
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/database"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"log"
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
	exporter.CreateMetric("testy", "this is a test")
	exporter.CreateMetric("testo", "thi is a test")
	for _, metric := range exporter.Metrics {
		metric.Inc()
	}
	exporter.RecordMetrics()
	fmt.Print(exporter.Metrics)
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	value := database.QueryMetric(db, *dbArg, "SELECT random() as Metric")
	fmt.Print(value)
	exporter.Listen(*portArg)
}
