package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var dbSessionsCleaned time.Time

func init() {
	var err error

	dbSessionsCleaned = time.Now()

	DB, err = sql.Open("postgres", "postgres://ml:ml@localhost/ml?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your ml database.")
}
