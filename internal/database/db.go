package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Connect() *sql.DB {
	driverName, dataSourceName := getDataSourceName()
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Print(err.Error())
	}

	return db
}

func getDataSourceName() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver := os.Getenv("DRIVER")
	dbName := os.Getenv("DB_NAME")
	userName := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	result := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, host, port, dbName)

	return driver, result
}
