package model

import (
	"errors"
	"fruitsalad-server/database"
	"math"
)

const maxPoints = 5000

type Game struct {
	Id         int
	UserId     *int
	Value      Color
	isFinished bool
	Guesses    []Color
}

type GameGuess struct {
	Id     int
	GameId int
	Color
}

func (game Game) hasGuesses() bool {
	return game.Guesses != nil && len(game.Guesses) != 0
}

func (game Game) CalculateScore() (int, error) {
	if !game.hasGuesses() {
		return 0, errors.New("Game has to be guessed in order to calculate a score")
	}

	latestGuess := game.Guesses[len(game.Guesses)-1]

	percentage := float64(0)

	percentage += 1 - float64(latestGuess.Red-game.Value.Red)/float64(255)
	percentage += 1 - float64(latestGuess.Green-game.Value.Green)/float64(255)
	percentage += 1 - float64(latestGuess.Blue-game.Value.Blue)/float64(255)

	percentage /= 3
	score := int(math.Round(maxPoints * percentage))

	if score == 5000 {
		err := game.finish()
		if err != nil {
			return 0, err
		}
	}

	return score, nil
}

func (game *Game) finish() error {
	db := database.GetDatabaseConnection()
	_, err := db.Exec("UPDATE game SET is_finished=true WHERE id=$1", game.Id)

	if err != nil {
		return errors.Join(errors.New("Couldn't finish game"), err)
	}

	game.isFinished = true

	return nil
}

func (game *Game) AddGuess(guess Color) error {
	if game.Guesses == nil {
		return errors.New("Guesses array is Nil, did you instatiate right?")
	}

	if game.isFinished {
		return errors.New("Game is already finished therefore no more scores can be added")
	}

	if len(game.Guesses) >= 5 {
		return errors.New("Game must not have more than 5 guesses")
	}

	db := database.GetDatabaseConnection()
	defer db.Close()

	_, err := db.Exec("INSERT INTO game_guess (game_id, red, green, blue) VALUES ($1, $2, $3, $4)", game.Id, guess.Red, guess.Green, guess.Blue)

	if err != nil {
		return err
	}

	game.Guesses = append(game.Guesses, guess)

	if len(game.Guesses) >= 5 {
		game.finish()
	}

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
	err := res.Scan(&game.Id, &game.UserId, &game.Value.Red, &game.Value.Green, &game.Value.Blue, &game.isFinished)
	if err != nil {
		return nil, err
	}

	game.Guesses = make([]Color, 0)

	guesses, err := db.Query("SELECT red, green, blue FROM game_guess WHERE game_id=$1", game.Id)

	if err != nil {
		return nil, err
	}

	for guesses.Next() {
		guess := new(GameGuess)
		err := guesses.Scan(&guess.Red, &guess.Green, &guess.Blue)
		if err != nil {
			return nil, errors.Join(errors.New("Reading the guesses failed"), err)
		}

		game.Guesses = append(game.Guesses, guess.Color)
	}

	return game, nil
}
