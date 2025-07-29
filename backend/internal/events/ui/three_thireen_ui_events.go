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

func (p PlayerDrewCardIntent) GetType() string {
	return TT_PLAYER_DREW_CARD
}
func (p PlayerDrewCardIntent) GetGameID() string {
	return p.GameID
}

type PlayerDiscardedCardIntent struct {
	GameID   string      `json:"game_id"`
	PlayerID string      `json:"player_id"`
	Card     interface{} `json:"card"`
}

func (p PlayerDiscardedCardIntent) GetType() string {
	return TT_PLAYER_DISCARDED_CARD
}

func (p PlayerDiscardedCardIntent) GetGameID() string {
	return p.GameID
}

type PlayerWentOutIntent struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}
func (p PlayerWentOutIntent) GetType() string {
	return TT_PLAYER_WENT_OUT
}
func (p PlayerWentOutIntent) GetGameID() string {
	return p.GameID
}