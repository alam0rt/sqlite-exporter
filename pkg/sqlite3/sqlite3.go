package database

import (
	"fmt"
	"github.com/mattn/go-sqlite3"
)

var d *sqlite3.SQLiteDriver

func open(dbFile string) {
	d.Open("Data Source=Algo.db; Version=3; ReadOnly=True;")
}

func main() {
	fmt.Println("debug")
}
