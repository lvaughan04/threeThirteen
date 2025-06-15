package models

import (
	"fmt"
	"backend/internal/utils"
	"backend/internal/types"
)
type Game struct {
	GameID      string   `json:"game_id"`
	Players     []Player `json:"players"`
	Deck        []types.Card   `json:"deck"`
	DiscardPile []types.Card   `json:"discard_pile"`
	CurrentTurn int      `json:"current_turn"` // Index of player whose turn it is
	GameState   string   `json:"game_state"`
	Round       int      `json:"round"`
}

func (g *Game) PushToDiscardPile(card types.Card) {
	g.DiscardPile = append([]types.Card{card}, g.DiscardPile...)
}

func (g *Game) PopFromDiscardPile() (types.Card, error) {
	if len(g.DiscardPile) == 0 {
		return types.Card{} , fmt.Errorf("cannot remove from discard pile when discard pile is empty")
	}

	card := g.DiscardPile[0]
	g.DiscardPile = g.DiscardPile[1:]
	return card, nil
}

func (g *Game) PeekFromDiscardPile() (types.Card, error) {
	if len(g.DiscardPile) == 0 {
		return types.Card{} , fmt.Errorf("cannot remove from discard pile when discard pile is empty")
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

func (g *Game) DrawFromDeck(player *Player) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if len(g.Deck) == 0 {
		//reshuffle the discard pile into the Deck, leaving the top card of the discard pile
		// need this to work so there is never a delay for the user
		topCard := g.DiscardPile[0]
		g.Deck = append([]types.Card{}, g.DiscardPile[1:]...) // make a copy
		utils.ShuffleDeck(g.Deck)
		g.DiscardPile = []types.Card{topCard}
	}
	if len(player.Hand) >= g.Round+1 {
		return fmt.Errorf("player cannot draw more cards this round")
	}

	//actual logic
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	player.AddCardToHand(card)

	return nil
}

func (g *Game) DrawFromDiscard(player *Player) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if len(player.Hand) >= g.Round+1 {
		return fmt.Errorf("player cannot draw more cards this round")
	}

	card, err := g.PopFromDiscardPile()
	if err != nil {
		return err // discard pile is empty
	}

	player.AddCardToHand(card)
	return nil
}

func (g *Game) DiscardFromHand(player *Player, card types.Card) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if player.RemoveCardFromHand(card) != nil {
		g.PushToDiscardPile(card)
	}
	return nil
}


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