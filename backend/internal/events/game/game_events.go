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
