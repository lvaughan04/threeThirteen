package ui

const (
	TT_PLAYER_DREW_CARD      = "TT_PLAYER_DREW_CARD"
	TT_PLAYER_DISCARDED_CARD = "TT_PLAYER_DISCARDED_CARD"
	TT_PLAYER_WENT_OUT       = "TT_PLAYER_WENT_OUT"
)

type PlayerDrewCardIntent struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
	Source   string `json:"source"`
}

type PlayerDiscardedCardIntent struct {
	GameID   string      `json:"game_id"`
	PlayerID string      `json:"player_id"`
	Card     interface{} `json:"card"`
}

type PlayerWentOutIntent struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}