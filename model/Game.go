package model

import (
	"database/sql"
	"errors"
	"fruitsalad-server/database"
	"math"
)

const maxPoints = 5000

type Game struct {
	Id      int
	UserId  sql.NullInt64
	Value   Color
	Guesses []Color
}

type GameGuess struct {
	Id     int
	GameId int
	Color
}

func (game Game) hasGuesses() bool {
	return game.Guesses == nil || len(game.Guesses) == 0
}

func (game Game) CalculateScore() (int, error) {
	if game.hasGuesses() {
		return 0, errors.New("Game has to be guessed in order to calculate a score")
	}

	latestGuess := game.Guesses[len(game.Guesses)-1]

	percentage := 0.0

	percentage += float64(latestGuess.Red) / float64(game.Value.Red)
	percentage += float64(latestGuess.Green) / float64(game.Value.Green)
	percentage += float64(latestGuess.Blue) / float64(game.Value.Blue)

	percentage /= 3

	return int(math.Round(maxPoints * percentage)), nil
}

func (game *Game) SetGuess(guess Color) error {
	if game.hasGuesses() {
		return errors.New("Game was already guessed")
	}

	db := database.GetDatabaseConnection()
	defer db.Close()

	_, err := db.Exec("UPDATE game SET guess_red=$1, guess_green=$2, guess_blue=$3 WHERE id=$4", guess.Red, guess.Green, guess.Blue, game.Id)

	if err != nil {
		return err
	}

	game.Guesses = append(game.Guesses, guess)

	return nil
}

// Creates a new game.
// If userId is nil then it will be a guest game
func NewGame(userId *int) (*Game, error) {
	db := database.GetDatabaseConnection()
	defer db.Close()

	value := GetRandomRgbValue()

	res := db.QueryRow("INSERT INTO game (user_id, red, green, blue) VALUES ($1, $2, $3, $4) RETURNING id",
		userId, value.Red, value.Blue, value.Green)

	var gameId int
	err := res.Scan(&gameId)

	if err != nil {
		return nil, err
	}

	return GetGameById(gameId)
}

func GetGameById(id int) (*Game, error) {
	db := database.GetDatabaseConnection()
	defer db.Close()
	res := db.QueryRow("SELECT * FROM game WHERE id=$1", id)

	game := new(Game)
	err := res.Scan(&game.Id, &game.UserId, &game.Value.Red, &game.Value.Green, &game.Value.Blue)
	if err != nil {
		return nil, err
	}

	game.Guesses = make([]Color, 0)

	guesses, err := db.Query("SELECT * FROM game_guess WHERE game_id=$1", game.Id)

	if err != nil {
		return nil, err
	}

	for guesses.Next() {
		guess := new(GameGuess)
		err := guesses.Scan(&guess)
		if err != nil {
			return nil, err
		}

		game.Guesses = append(game.Guesses, guess.Color)
	}

	return game, nil
}
