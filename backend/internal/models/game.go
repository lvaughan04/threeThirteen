package models

type Game struct {
	GameID      string   `json:"game_id"`
	Players     []Player `json:"players"`
	Deck        []Card   `json:"deck"`
	DiscardPile []Card   `json:"discard_pile"` 
	CurrentTurn int      `json:"current_turn"` // Index of player whose turn it is
	GameState   string   `json:"game_state"`
	Round       int      `json:"round"`
}
