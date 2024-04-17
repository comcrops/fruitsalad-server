package model

// import "fruitsalad-server/database"

type User struct {
	Id       int
	Username string
	password string
	token    string
}

func NewUser(username, password string) {
	// db := database.GetDatabaseConnection()
	
}

