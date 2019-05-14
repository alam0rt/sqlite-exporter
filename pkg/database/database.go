package database

import (
	"bytes"
	"database/sql"
	"github.com/alam0rt/sqlite-exporter/pkg/logging"
	"github.com/mattn/go-sqlite3"
	"os/exec"
	"strconv"
	"strings"
)

var (
	d *sqlite3.SQLiteDriver
)

// QueryMetric is a function which takes a database object and a query
// and returns a float64. The supplied query must  only return a single
// result and must be a number of some sort.
func QueryMetric(db *sql.DB, query string) float64 {
	var metric float64

	stmt, err := db.Prepare(query)
	if err != nil {
		logging.Error.Printf("%q: %s\n", err, query)
	}
	defer stmt.Close() // close

	rows, err := stmt.Query()
	if err != nil {
		logging.Error.Printf("%q: %s\n", err, query)
	}

	for rows.Next() {
		err := rows.Scan(&metric) // make it so column name aliases are used as metric names
		if err != nil {
			logging.Error.Print(err)
			logging.Error.Fatal("'", query, "' may contain an error and/or is returning more than one result")
		}
	}

	return metric

}

// QueryMetricFallback is an alternative way to query the database and
// return a float64. It uses the sqlite3 binary in $PATH to exec a query
func QueryMetricFallback(db string, query string) float64 {
	var metric float64

	cmd := exec.Command("sqlite3", db, query)
	// create a new bytes buffer
	var out bytes.Buffer
	// set stdout to be the address of said buffer
	cmd.Stdout = &out

	// run the command
	err := cmd.Run()
	if err != nil {
		logging.Error.Fatal(err)
	}

	// convert to a string
	s := out.String()
	// remove trailing newline
	s = strings.Replace(s, "\n", "", 1)
	// convert to a float
	metric, err = strconv.ParseFloat(s, 0)
	if err != nil {
		logging.Error.Print(err)
	}

	return metric

}
