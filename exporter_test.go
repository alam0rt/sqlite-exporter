package exporter

import (
	"testing"
)

var tests = []testmetric{
	{"test", "this is a test", 0},
}

type testmetric struct {
	name   string
	help   string
	result float64
}

// Creates a test metric and tests if the value is correct
func TestCreateMetric(t *testing.T) {
	for _, pair := range tests {
		m := CreateMetric(pair.name, pair.help)
		if m.value != pair.result {
			t.Error(
				"For", pair.name,
				"Expected", pair.result,
				"Got", m.value,
			)
		}

	}
}
