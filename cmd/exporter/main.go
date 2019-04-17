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
	//	configArg := flag.String("config", "configuration.yml", "sqlite-exporter configuration file")

	flag.Parse()

	//c := config.ProcessConfig(*configArg)
	db, err := sql.Open("sqlite3", *dbArg)
	if err != nil {
		logger.Fatal("Unable to open database")
	}

	// CREATE
	demoMetric := "sqlite3_random"
	exporter.CreateMetric(
		demoMetric,
		"Example metric which is updated with the output of SELECT random()",
	)
	// todo: i need to make sure the value is obtained and the metric is updated in one fell swoop
	go func() {
		for {
			value := database.QueryMetric(db, "Algo.db", "SELECT random() as metric")
			exporter.SetMetric(demoMetric, value)
			time.Sleep(2 * time.Second)
		}
	}()
	// CREATE END

	go exporter.RecordMetrics()
	logger.Print("Listening on port " + *portArg + "...")
	logger.Print("Opening " + *dbArg)
	fmt.Print(&buf)
	exporter.Listen(*portArg)
}
