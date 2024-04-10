package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	. "fruitsalad-server/model"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type MyStruct struct {
	Beidl int8
}

type User struct {
	Username string
	password string
	token string
}

type Guess struct {
	value RgbValue
	Guess RgbValue
	Player User
}

func main() {
	db, err := sql.Open("postgres", loadDatabaseUrl())

	if err != nil {
		log.Fatalf("Error while connecting to db")
	}
	defer db.Close()

	http.HandleFunc("/game/new", generateRandomGame)
	log.Fatal(http.ListenAndServe(":8080", nil))
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


func generateRandomGame(w http.ResponseWriter, req *http.Request) {
	value := GetRandomRgbValue()
	jsonData, err := json.Marshal(value)
	fmt.Printf("%s", string(jsonData))

	if err != nil {
		log.Printf("There was an error converting the RgbValue to JSON")
	}
	io.WriteString(w, string(jsonData))
	
}

