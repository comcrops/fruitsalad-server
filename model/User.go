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
	Color
}

func (game Game) calculateScore() float64 {
	return 1
}

func (user User) generateNewGame() (*Game, error) {
	db := GetDatabaseConnection()
	value := GetRandomRgbValue()

	res := db.QueryRow("INSERT INTO game (user_id, red, green, blue) VALUES ($1, $2, $3, $4)",
		user.Id, value.Red, value.Blue, value.Green)

	var game Game
	err := res.Scan(&game)

	if err != nil {
		return nil, err	
	}
	
	return &game, nil
}
