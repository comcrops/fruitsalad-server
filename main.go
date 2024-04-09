package main

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	. "fruitsalad-server/model"
	"io"
	"log"
	"net/http"
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
	http.HandleFunc("/game/new", generateRandomGame)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

