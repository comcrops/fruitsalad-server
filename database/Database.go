package database

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetDatabaseConnection() *sql.DB {
	connection, err := sql.Open("postgres", loadDatabaseUrl())

	if err != nil {
		log.Fatalf("Error while connecting to database")
	}

	return connection
}

func loadDatabaseUrl() string {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file")
	}

	connectionString := os.Getenv("DATABASE_URL")

	if connectionString == "" {
		log.Printf("Connection string wasn't found or is empty!")
	}
	return connectionString
}
