# sqlite-exporter #
** A Prometheus exporter for exposing sqlite3 queries as metrics **

By default the exporter uses the go-sqlite3 driver at https://github.com/mattn/go-sqlite3. I have encountered issues where an application will lock the sqlite3 database and as a work around have included a function called QueryMetricFallback() which execs() the sqlite3 binary in $PATH.

# Building #

As github.com/mattn/go-sqlite3 uses cgo, the resulting binary will *not* be statically linked and can cause portability issues on systems with an older glibc. To force static compilation I suggest building a Docker image based on golang:stretch and compiling with:
`apt-get update
apt-get install -y musl musl-tools
CC=/usr/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"'  -o sqlite-exporter cmd/exporter/main.go`

# Usage #

To use sqlite-exporter, simply build the binary and run it like:

`  -config string
    	sqlite-exporter configuration file (default "configuration.yml")
  -database string
    	a sqlite3 database (default "database.db")
  -interval float
    	per second interval to query the database (default 60)
  -port string
    	a port to listen on (default "9001")`

The configuration file is in YAML format and looks like:

`sqlite_random:
  query: select random()
  description: returns a random float`
