package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fruitsalad-server/model"
	"io"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

type Score struct {
	Score int
}

type HttpError struct {
	Detail string `json:"detail"`
}

func main() {
	http.HandleFunc("/game/new", generateRandomGame)
	http.HandleFunc("/game/addGuess", guessGame)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAuthorizationHeader(req *http.Request) *string {
	val := req.Header.Get("Authorization")

	if val != "" {
		return &val
	}

	return nil
}

func jsonHeaders(w http.ResponseWriter){
	w.Header().Add("Content-Type", "application/json")
}

func methodNotAllowed(w http.ResponseWriter, req *http.Request) {
	errorMessage := fmt.Sprintf("Method %s is not allowed", req.Method)
	slog.Error(errorMessage)
	errorRequest(w, http.StatusMethodNotAllowed, errors.New(errorMessage))
}

//Send an empty error string for a default bad request return
func badRequest(w http.ResponseWriter, err error) {
	if err.Error() == "" {
		errorRequest(w, http.StatusBadRequest, errors.New("400 - Bad Request"))
	} else {
		errorRequest(w, http.StatusBadRequest, err)
	}
}

func errorRequest(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(&HttpError{
		Detail: err.Error(),
	})
	io.WriteString(w, string(data))
}

// Send an empty string for additionalInfo for default behaviour
func notFound(w http.ResponseWriter, additionalInfo string) {
	if additionalInfo != "" {
		additionalInfo = fmt.Sprintf("404 - Not Found: %s", additionalInfo)
	} else {
		additionalInfo = "404 - Not found"
	}

	errorRequest(w, http.StatusNotFound, errors.New(additionalInfo))
}

func generateRandomGame(w http.ResponseWriter, req *http.Request) {
	token := getAuthorizationHeader(req)

	var game *model.Game

	if token == nil {
		newGame, err := model.NewGame(nil)
		if err != nil {
			log.Printf("Error while creating guest-game: %s", err)
			badRequest(w, err)
			return
		}
		game = newGame
	} else {
		user, err := model.GetUserByToken(*token)
		if err != nil {
			log.Printf("UserByToken (%s) not found: %s", *token, err)
			notFound(w, "user")
			return
		}
		newGame, err := model.NewGame(&user.Id)
		if err != nil {
			log.Printf("Error while creating game for user: %+v %s", user, err)
			badRequest(w, err)
			return
		}
		game = newGame
	}

	jsonData, _ := json.Marshal(game)

	jsonHeaders(w)
	io.WriteString(w, string(jsonData))
}

func guessGame(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		methodNotAllowed(w, req)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var guess model.GameGuess
	err := decoder.Decode(&guess)

	if err != nil {
		log.Printf("Error while decoding guess: %s", err)
		badRequest(w, err)
		return
	}

	game, err := model.GetGameById(guess.GameId)

	if err != nil {
		log.Printf("Error while querying for game: %s", err)
		notFound(w, "Game")
		return
	}

	err = game.AddGuess(guess.Color)

	if err != nil {
		log.Printf("Error setting the guess: %s", err)
		badRequest(w, err)
		return
	}

	value, err := game.CalculateScore()

	if err != nil {
		log.Printf("Calculating the score for the current guess went wrong: %s", err)
		badRequest(w, err)
		return
	}

	score := &Score{
		Score: value,
	}

	jsonData, _ := json.Marshal(score)

	jsonHeaders(w)
	io.WriteString(w, string(jsonData))
}
