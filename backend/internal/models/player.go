package models

type Player struct {
	PlayerID   string `json:"player_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Hand     []Card `json:"hand"`
	IsTurn   bool   `json:"is_turn"`
	IsDealer bool   `json:"is_dealer"`
	IsActive bool   `json:"is_active"`
}
