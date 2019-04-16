package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
)

var (
	d      *sqlite3.SQLiteDriver
	result int64 // test
)

func QueryMetric(dbFile string, query string) float64 {
	db := open(dbFile)
	defer db.Close()

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
		err := rows.Scan(&metric)
		if err != nil {
			log.Print(err)
		}
		log.Println(metric)
	}

	return metric

}

func open(dbFile string) *sql.DB {
	log.Print("Opening " + dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Print(err)
	}
	return db
}
