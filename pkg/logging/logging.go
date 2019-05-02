package logging

import (
	"log"
	"os"
)

// This package sets up log handlers
var (
	Error  = log.New(os.Stderr, "sqlite-exporter: ", log.Ltime|log.Ldate)
	Output = log.New(os.Stdout, "sqlite-exporter: ", log.Ltime|log.Ldate)
	err    error
)
