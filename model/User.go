package model

import (
	"errors"
	. "fruitsalad-server/database"
)

type User struct {
	Id       int
	Username string
	password string
	token    string
}

type Game struct {
	Id     int
	UserId int
	Value  Color
	Guess  *Color
}

func (game Game) calculateScore() float64 {
	return 1
}

func (game *Game) setGuess(guess Color) error {
	if game.Guess != nil {
		return errors.New("Game was already guessed")
	}

	db := GetDatabaseConnection()

	_, err := db.Exec("UPDATE game SET guess_red=$1, guess_green=$2, guess_blue=$3 WHERE id=$4", guess.Red, guess.Green, guess.Blue, game.Id)

	if err != nil {
		return err
	}

	game.Guess = &guess

	return nil
}

func (user User) GenerateNewGame() (*Game, error) {
	db := GetDatabaseConnection()
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
	db := GetDatabaseConnection()
	res := db.QueryRow("SELECT * FROM game WHERE id=$1", id)

	game := new(Game)

	err := res.Scan(&game.Id, &game.UserId, &game.Value.Red, &game.Value.Green, &game.Value.Blue, &game.Guess.Red, &game.Guess.Green, &game.Guess.Blue)

	if err != nil {
		return nil, err
	}

	return game, nil
}
