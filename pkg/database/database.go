package database

import (
	"bitbucket.org/dragontailcom/sqlite-exporter/pkg/logging"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

var (
	d *sqlite3.SQLiteDriver
)

func QueryMetric(db *sql.DB, query string) float64 {
	var metric float64

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("%q: %s\n", err, query)
	}
	defer stmt.Close() // close

	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("%q: %s\n", err, query)
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
