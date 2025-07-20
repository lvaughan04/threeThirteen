package events

import (
	"time"
)

type GameEvent struct {
	Type 	string    `json:"type"`
    GameType string    `json:"game_type"`
	GameID string    `json:"game_id"`
	Data 	interface{} `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}