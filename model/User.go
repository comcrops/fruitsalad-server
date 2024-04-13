package model

import (
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
	value  Color
	guess  Color
}

func (game Game) calculateScore() float64 {
	return 1
}

func (user User) generateNewGame() (*Game, error) {
	db := GetDatabaseConnection()
	value := GetRandomRgbValue()

	res := db.QueryRow("INSERT INTO game (user_id, red, green, blue) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Id, value.Red, value.Blue, value.Green)

	var gameId int
	err := res.Scan(&gameId)

	if err != nil {
		return nil, err
	}

	return getGameById(gameId)
}

func getGameById(id int) (*Game, error) {
	db := GetDatabaseConnection()
	res := db.QueryRow("SELECT * FROM game WHERE id=$1", id)

	game := new(Game)

	err := res.Scan(&game.Id, &game.UserId, &game.value.Red, &game.value.Green, &game.value.Blue, &game.guess.Red, &game.guess.Green, &game.guess.Blue)

	if err != nil {
		return nil, err
	}

	return game, nil
}
