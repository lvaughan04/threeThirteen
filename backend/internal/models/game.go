package models

import "fmt"

type Game struct {
	GameID      string   `json:"game_id"`
	Players     []Player `json:"players"`
	Deck        []Card   `json:"deck"`
	DiscardPile []Card   `json:"discard_pile"`
	CurrentTurn int      `json:"current_turn"` // Index of player whose turn it is
	GameState   string   `json:"game_state"`
	Round       int      `json:"round"`
}

func (g *Game) PushToDiscardPile(card Card) {
	g.DiscardPile = append([]Card{card}, g.DiscardPile...)
}

func (g *Game) PopFromDiscardPile() (Card, error) {
	if len(g.DiscardPile) == 0 {
		return Card{} , fmt.Errorf("cannot remove from discard pile when discard pile is empty")
	}

	card := g.DiscardPile[0]
	g.DiscardPile = g.DiscardPile[1:]
	return card, nil
}

func (g *Game) PeekFromDiscardPile() (Card, error) {
	if len(g.DiscardPile) == 0 {
		return Card{} , fmt.Errorf("cannot remove from discard pile when discard pile is empty")
	}

	return g.DiscardPile[0], nil
}


func (g *Game) IsPlayerInGame(player *Player) bool {
	for _, p := range g.Players {
		if p.PlayerID == player.PlayerID {
			return true
		}
	}
	return false
}


/**
 * Player actions
 */

//GoOut 

/*
 * Game Maintenance and Round logic 
	InitilizeGame()
	StartGame()
	EndRound()
	StartRound()
	EndGame()
*/
/*
 * Given An array of players, and sets all of the game information to its starting values
 */
func (g *Game) InitilizeGame(players *[]Player) {

}
