package database

import (
	"database/sql"
	"testing"
)

var tests = []testquery{
	{"SELECT count(random())", 1},
}

type testquery struct {
	query  string
	result float64
}

// Runs through a list of query / result pairs and ensures we get the expected result
func TestQuery(t *testing.T) {
	db, err := sql.Open("sqlite3", "Test.db")
	if err != nil {
		t.Error(
			"For", "opening database in memory",
			"Expected", db,
			"Got", err,
		)
	}
	for _, pair := range tests {
		r := QueryMetric(db, "Test.db", pair.query)
		if r != pair.result {
			t.Error(
				"For", pair.query,
				"Expected", pair.result,
				"Got", r,
			)
		}

	}
}
