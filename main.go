package main

import (
	"encoding/json"
	"fmt"
	"fruitsalad-server/model"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Guess struct {
	GameId int
	Guess  model.Color
}

func main() {
	http.HandleFunc("/game/new", generateRandomGame)
	http.HandleFunc("/game/submit", guessGame)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAuthorizationHeader(req *http.Request) *string {
	val := req.Header.Get("Authorization")

	if val != "" {
		return &val
	}

	return nil
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


// Send an empty string for additionalInfo for default behaviour
func notFound(w http.ResponseWriter, additionalInfo string) {
	w.WriteHeader(http.StatusNotFound)
	if additionalInfo != "" {
		io.WriteString(w, fmt.Sprintf("404 - Not Found: %s", additionalInfo))
	} else {
		io.WriteString(w, "404 - Not Found")
		
	}
}

func generateRandomGame(w http.ResponseWriter, req *http.Request) {
	token := getAuthorizationHeader(req)

	var game *model.Game

	if token == nil {
		newGame, err := model.NewGame(nil)
		if err != nil {
			log.Printf("Error while creating guest-game: %s", err)
			badRequest(w)
			return
		}
		game = newGame
	} else {
		user, err := model.GetUserByToken(*token)
		if err != nil {
			log.Printf("UserByToken (%s) not found: %s", token, err)
			notFound(w, "user")
			return
		}
		newGame, err := model.NewGame(&user.Id)
		if err != nil {
			log.Printf("Error while creating game for user: %+v %s", user, err)
			badRequest(w)
			return
		}
		game = newGame
	}

	jsonData, _ := json.Marshal(game)

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

	game, err := model.GetGameById(guess.GameId)

	if err != nil {
		log.Printf("Error while querying for game: %s", err)
		notFound(w, "Game")
		return
	}

	err = game.SetGuess(guess.Guess)
	 
	if err != nil {
		log.Printf("Error setting the guess: %s", err)
		badRequest(w)
		return
	}

	jsonData, _ := json.Marshal(game)

	io.WriteString(w, string(jsonData))
}
