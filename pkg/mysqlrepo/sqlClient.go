package mysqlrepo

import (
	"database/sql"
	"fmt"
	"os"
)

type DatabaseClient struct {
	db *sql.DB
}

var dbClient *DatabaseClient

func ConnectDB() *DatabaseClient {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	dbClient.db = db

	fmt.Println("Connected to the MySQL database!")
	return dbClient
}
