package game



const (
	TT_GAME_INITIALIZED  = "TT_GAME_INITIALIZED"
    TT_ROUND_STARTED     = "TT_ROUND_STARTED"
    TT_CARD_DRAWN        = "TT_CARD_DRAWN"
    TT_CARD_DISCARDED    = "TT_CARD_DISCARDED"
    TT_PLAYER_WENT_OUT   = "TT_PLAYER_WENT_OUT"
    TT_ROUND_ENDED       = "TT_ROUND_ENDED"
    TT_GAME_OVER         = "TT_GAME_OVER"
)

type TTCardDrawnData struct {
    PlayerID string      `json:"player_id"`
    Source   string      `json:"source"`
    Card     interface{} `json:"card,omitempty"`
}

type TTCardDiscardedData struct {
    PlayerID string      `json:"player_id"`
    Card     interface{} `json:"card"`
}

type TTRoundStartedData struct {
    Round         int `json:"round"`
    CardsPerPlayer int `json:"cards_per_player"`
    WildCardValue int `json:"wild_card_value"`
}
