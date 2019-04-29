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
	"time"
)

var DB *sql.DB // set up package level DB reference

func main() {
	// set up a log handler
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "sqlite_exporter: ", log.Lshortfile)
		err    error
	)
	logger.Print("starting sqlite-exporter...")

	// Set up argument parsing
	portArg := flag.String("port", "9001", "a port to listen on")
	dbArg := flag.String("database", "Algo.db", "a sqlite3 database")
	intervalArg := flag.Float64("interval", 60, "per second interval to query the database")
	configArg := flag.String("config", "configuration.yml", "sqlite-exporter configuration file")

	flag.Parse()

	// Open database
	DB, err = sql.Open("sqlite3", *dbArg)
	if err != nil {
		logger.Fatal("Unable to open database")
	}

	// load and process our config
	config.ProcessConfig(*configArg)

	c := config.Config // config.ProcessConfig returns a struct with all of the unmarshalled data
	for k, v := range c {
		exporter.CreateMetric(
			k, // name of metric
			v.Description,
			v.Query,
		)
	}

	// loop over our metrics
	metricsLoop(*intervalArg)

	logger.Print("Listening on 0.0.0.0:" + *portArg + "...")
	logger.Print("Opened " + *dbArg)
	fmt.Print(&buf)

	// finally we can listen on the provided TCP port
	exporter.Listen(*portArg)
}

// iterates over our collection of metrics every n seconds and updates them
func metricsLoop(i float64) {
	go func() {
		for {
			// iterate over Metrics slice and increase()
			for _, m := range exporter.MetricsMap {
				exporter.SetMetric(
					m.Name,
					// todo: start here tomorrow
					database.QueryMetric(DB, m.Query),
				)
				exporter.UpdateMetric(m)
			}
			d := time.Duration(i) * time.Second
			time.Sleep(d)
		}
	}()
}
