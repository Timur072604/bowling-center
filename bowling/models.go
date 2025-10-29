package bowling

import "time"

type Client struct {
	ID          int       `json:"ID"`
	GameString  string    `json:"GameString"`
	ArrivalTime time.Time `json:"-"`
}

type GameStatus string

const (
	StatusPlayed GameStatus = "Сыграл"
	StatusLeft   GameStatus = "Ушел"
)

type GameResult struct {
	Client Client     `json:"Client"`
	Score  int        `json:"Score"`
	Status GameStatus `json:"Status"`
	Err    error      `json:"Err,omitempty"`
}
