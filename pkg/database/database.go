package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
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
		err := rows.Scan(&metric) // make it so column names are used as metric names
		if err != nil {
			log.Print(err)
		}
	}

	return metric

}
