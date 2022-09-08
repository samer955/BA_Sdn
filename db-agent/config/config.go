package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Connection  *sql.DB
	TableSchema string
}

var connection *sql.DB
var tableSchema string

func GetConfig() Config {
	return Config{
		Connection:  connection,
		TableSchema: tableSchema,
	}
}

func init() {
	err := godotenv.Load("db.env")

	if err != nil {
		log.Println("Error loading db.env file")
		os.Exit(1)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_DRIVER := os.Getenv("DB_DRIVER")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open(DB_DRIVER, psqlInfo)

	if err != nil {
		log.Println("cannot connect to Database")
		panic(err)
	}
	log.Println("connected to Database")
	connection = db

	sqlSchema, err := ioutil.ReadFile("repository/migrations/000001_init_schema.up.sql")
	if err != nil {
		log.Println("unable to read schema.sql")
		panic(err)
	}
	tableSchema = string(sqlSchema)
}
