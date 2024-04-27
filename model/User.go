package model

import "fruitsalad-server/database"

type User struct {
	Id       int
	Username string
	password string
	token    string
}

func NewUser(username, password string) {
	// db := database.GetDatabaseConnection()
}

func (user User) GenerateNewGame() (*Game, error) {
	db := database.GetDatabaseConnection()
	defer db.Close()

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

func GetUserByToken(token string) (*User, error) {
	return nil, nil
}
