package model

import (
	"errors"
	"fruitsalad-server/database"
	"math"
)

const maxPoints = 5000

type Game struct {
	Id     int
	UserId int
	Value  Color
	Guess  *Color
}

func (game Game) CalculateScore() (int, error) {
	if game.Guess == nil {
		return 0, errors.New("Game has to be guessed in order to calculate a score")
	}

	percentage := 0.0

	percentage += float64(game.Guess.Red) / float64(game.Value.Red)
	percentage += float64(game.Guess.Green) / float64(game.Value.Green)
	percentage += float64(game.Guess.Blue) / float64(game.Value.Blue)

	percentage /= 3

	return int(math.Round(maxPoints*percentage)), nil
}

func (game *Game) SetGuess(guess Color) error {
	if game.Guess != nil {
		return errors.New("Game was already guessed")
	}

	db := database.GetDatabaseConnection()

	_, err := db.Exec("UPDATE game SET guess_red=$1, guess_green=$2, guess_blue=$3 WHERE id=$4", guess.Red, guess.Green, guess.Blue, game.Id)

	if err != nil {
		return err
	}

	game.Guess = &guess

	return nil
}

func (user User) GenerateNewGame() (*Game, error) {
	db := database.GetDatabaseConnection()
	value := GetRandomRgbValue()

	res := db.QueryRow("INSERT INTO game (user_id, red, green, blue) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Id, value.Red, value.Blue, value.Green)

	var gameId int
	err := res.Scan(&gameId)

	if err != nil {
		return nil, err
	}

	return GetGameById(gameId)
}

func GetGameById(id int) (*Game, error) {
	db := database.GetDatabaseConnection()
	res := db.QueryRow("SELECT * FROM game WHERE id=$1", id)

	game := new(Game)

	err := res.Scan(&game.Id, &game.UserId, &game.Value.Red, &game.Value.Green, &game.Value.Blue, &game.Guess.Red, &game.Guess.Green, &game.Guess.Blue)

	if err != nil {
		return nil, err
	}

	return game, nil
}
