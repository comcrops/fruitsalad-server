package database

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

// Hopefully somewhat of a singleton, dk if that actually works
func GetDatabaseConnection() *sql.DB {
	if db != nil {
		return db
	}

	connection, err := sql.Open("postgres", loadDatabaseUrl())

	if err != nil {
		log.Fatalf("Error while connecting to database")
	}

	db = connection
	defer db.Close()

	return db
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
