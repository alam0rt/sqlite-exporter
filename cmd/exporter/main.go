package main

import (
	"bitbucket.org/dragontailcom/sqlite-exporter"
	"bitbucket.org/dragontailcom/sqlite-exporter/internal/config"
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/database"
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/logging"
	"database/sql"
	"flag"
	"os"
	"time"
)

var DB *sql.DB // set up package level DB reference

const Version = "0.1.1~BTO"

func main() {
	logging.Output.Printf("starting sqlite-exporter %s...", Version)

	// Set up argument parsing
	portArg := flag.String("port", "9001", "a port to listen on")
	dbArg := flag.String("database", "Algo.db", "a sqlite3 database")
	intervalArg := flag.Float64("interval", 60, "per second interval to query the database")
	configArg := flag.String("config", "configuration.yml", "sqlite-exporter configuration file")

	flag.Parse()

	// Open database
	_, err := os.Open(*dbArg)

	if err != nil {
		logging.Error.Fatal("Unable to open database file: ", *dbArg)
	}
	DB, err = sql.Open("sqlite3", *dbArg+"?mode=ro&_busy_timeout=300&cache=shared")
	if err != nil {
		logging.Error.Fatal("Unable to read database")
	}

	// fix issue where multiple threads were locking DB
	DB.SetMaxOpenConns(1)

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
	metricsLoop(*intervalArg, *dbArg)

	logging.Output.Print("Listening on 0.0.0.0:" + *portArg + "...")
	logging.Output.Print("Opened " + *dbArg)

	// finally we can listen on the provided TCP port
	exporter.Listen(*portArg)
}

// metricsLoop iterates over our collection of metrics every n seconds and updates them
func metricsLoop(i float64, dbPath string) {
	go func() {
		for {
			// iterate over Metrics slice and increase()
			for _, m := range exporter.MetricsMap {
				exporter.SetMetric(
					m.Name,
					// todo: start here tomorrow
					//database.QueryMetric(DB, m.Query),
					database.QueryMetricFallback(dbPath, m.Query),
				)
				exporter.UpdateMetric(m)
			}
			d := time.Duration(i) * time.Second
			time.Sleep(d)
		}
	}()
}
