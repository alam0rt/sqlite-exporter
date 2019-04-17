package main

import (
	"bitbucket.org/dragontailcom/sqlite-exporter"
	//	"bitbucket.org/dragontailcom/sqlite-exporter/internal/config"
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
		logger = log.New(&buf, "logger: ", log.Lshortfile)
		err    error
	)
	logger.Print("starting sqlite-exporter...")

	// Set up argument parsing
	portArg := flag.String("port", "9001", "a port to listen on")
	dbArg := flag.String("database", "Algo.db", "a sqlite3 database")
	intervalArg := flag.Float64("interval", 15, "per second interval to query the database")
	//	configArg := flag.String("config", "configuration.yml", "sqlite-exporter configuration file")

	flag.Parse()

	//c := config.ProcessConfig(*configArg)
	DB, err = sql.Open("sqlite3", *dbArg)
	if err != nil {
		logger.Fatal("Unable to open database")
	}

	// CREATE
	demoMetric := "sqlite3_random"
	exporter.CreateMetric(
		demoMetric,
		"Example metric which is updated with the output of SELECT random()",
	)
	// CREATE END

	metricsLoop(*intervalArg)
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	fmt.Print(&buf)
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
					database.QueryMetric(DB, "Algo.db", "SELECT random() as metric"),
				)
				exporter.UpdateMetric(m)
			}
			d := time.Duration(i) * time.Second
			time.Sleep(d)
		}
	}()
}
