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

	dbSessionsCleaned = time.Now()

	var err error
	POSTGRES_CONNECTION := "user=mlpsuser password=m!psusER19! dbname=ml host=/cloudsql/medaliving:europe-west3:pgsqldev"

	DB, err = sql.Open("postgres", POSTGRES_CONNECTION)
	if err != nil {
		panic(err)
	}

	/*
		if err = DB.Ping(); err != nil {
			panic(err)
		}

	*/
	fmt.Println("You connected to your ml database.")

}
