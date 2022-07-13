package database

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "mydb"
)

var db *sql.DB

func init() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	postgresDb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected!")
	}

	db = postgresDb

}

func Database() *sql.DB {
	return db
}
