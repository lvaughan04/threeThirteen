package game

import (
	"time"
)

type GameEvent struct {
	Type      string      `json:"type"`
	GameID    string      `json:"game_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func (g GameEvent) GetType() string {
	return g.Type
}

func (g GameEvent) GetGameID() string {
	return g.GameID
}
