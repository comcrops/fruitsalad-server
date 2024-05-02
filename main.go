package main

import (
	"encoding/json"
	"fmt"
	fsHttp "fruitsalad-server/http"
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


func main() {
	http.HandleFunc("/game/new", generateRandomGame)
	http.HandleFunc("/game/addGuess", guessGame)

	log.Fatal(http.ListenAndServe(":8080", nil))
}


func generateRandomGame(w http.ResponseWriter, req *http.Request) {
	token := fsHttp.GetAuthorizationHeader(req)

	var game *model.Game

	if token == nil {
		newGame, err := model.NewGame(nil)
		if err != nil {
			slog.Error(fmt.Sprintf("Error while creating guest-game: %s", err))
			fsHttp.BadRequest(w, err)
			return
		}
		game = newGame
	} else {
		user, err := model.GetUserByToken(*token)
		if err != nil {
			slog.Error(fmt.Sprintf("UserByToken (%s) not found: %s", *token, err))
			fsHttp.NotFound(w, "user")
			return
		}
		newGame, err := model.NewGame(&user.Id)
		if err != nil {
			slog.Error(fmt.Sprintf("Error while creating game for user: %+v %s", user, err))
			fsHttp.BadRequest(w, err)
			return
		}
		game = newGame
	}

	jsonData, _ := json.Marshal(game)

	fsHttp.JsonHeaders(w)
	io.WriteString(w, string(jsonData))
}

func guessGame(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fsHttp.MethodNotAllowed(w, req)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var guess model.GameGuess
	err := decoder.Decode(&guess)

	if err != nil {
		slog.Error(fmt.Sprintf("Error while decoding guess: %s", err))
		fsHttp.BadRequest(w, err)
		return
	}

	game, err := model.GetGameById(guess.GameId)

	if err != nil {
		slog.Error(fmt.Sprintf("Error while querying for game: %s", err))
		fsHttp.NotFound(w, "Game")
		return
	}

	err = game.AddGuess(guess.Color)

	if err != nil {
		slog.Error(fmt.Sprintf("Error setting the guess: %s", err))
		fsHttp.BadRequest(w, err)
		return
	}

	value, err := game.CalculateScore()

	if err != nil {
		slog.Error(fmt.Sprintf("Calculating the score for the current guess went wrong: %s", err))
		fsHttp.BadRequest(w, err)
		return
	}

	score := &Score{
		Score: value,
	}

	jsonData, _ := json.Marshal(score)

	fsHttp.JsonHeaders(w)
	io.WriteString(w, string(jsonData))
}
