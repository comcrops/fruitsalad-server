package main

import (
	"encoding/json"
	. "fruitsalad-server/model"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Guess struct {
	GameId int
	Guess  Color
}

func main() {
	http.HandleFunc("/game/new", generateRandomGame)
	http.HandleFunc("/game/submit", guessGame)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	log.Printf("Method %s is not allowed", req.Method)
	w.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(w, "405 - Method Not Allowed")
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, "400 - Bad Request")
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "404 - Not Found")
}

func generateRandomGame(w http.ResponseWriter, req *http.Request) {
	value := GetRandomRgbValue()
	jsonData, _ := json.Marshal(value)

	io.WriteString(w, string(jsonData))
}

func guessGame(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		methodNotAllowed(w, req)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var guess Guess
	err := decoder.Decode(&guess)

	if err != nil {
		log.Printf("Error while decoding guess: %s", err)
		badRequest(w)
		return
	}

	game, err := GetGameById(guess.GameId)

	if err != nil {
		log.Printf("Error while querying for game: %s", err)
		notFound(w)
		return
	}

	game.SetGuess(guess.Guess)
}
