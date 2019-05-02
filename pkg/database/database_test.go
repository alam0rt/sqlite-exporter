package database

import (
	"database/sql"
	"testing"
)

var tests = []testquery{
	{"SELECT count(random())", 1},
	{"SELECT 0", 0},
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
		r := QueryMetric(db, pair.query)
		if r != pair.result {
			t.Error(
				"For", pair.query,
				"Expected", pair.result,
				"Got", r,
			)
		}

	}
}

func BenchmarkOpenDb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = sql.Open("sqlite3", "Test.db")
	}
}

func BenchmarkQueryMetric(b *testing.B) {
	db, _ := sql.Open("sqlite3", "Test.db")
	for _, pair := range tests {
		b.Run(pair.query, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = QueryMetric(db, pair.query)
			}
		})
	}
}
