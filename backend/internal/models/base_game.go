package models

import (
	"backend/internal/types"
	"fmt"
)

type BaseGame struct {
	GameID      string             `json:"game_id"`
	Players     map[string]*Player `json:"players"`
	Deck        []types.Card       `json:"deck"`
	DiscardPile []types.Card       `json:"discard_pile"`
	CurrentTurn int                `json:"current_turn"`
	GameState   string             `json:"game_state"`
	Round       int                `json:"round"`
}

// Common methods that work the same for all card games
func (g *BaseGame) PushToDiscardPile(card types.Card) {
	g.DiscardPile = append([]types.Card{card}, g.DiscardPile...)
}

func (g *BaseGame) PopFromDiscardPile() (types.Card, error) {
	if len(g.DiscardPile) == 0 {
		return types.Card{}, fmt.Errorf("cannot remove from discard pile when discard pile is empty")
	}
	card := g.DiscardPile[0]
	g.DiscardPile = g.DiscardPile[1:]
	return card, nil
}

func (g *BaseGame) PeekFromDiscardPile() (types.Card, error) {
	if len(g.DiscardPile) == 0 {
		return types.Card{}, fmt.Errorf("cannot remove from discard pile when discard pile is empty")
	}
	return g.DiscardPile[0], nil
}

func (g *BaseGame) IsPlayerInGame(player *Player) bool {
	for _, p := range g.Players {
		if p.PlayerID == player.PlayerID {
			return true
		}
	}
	return false
}

// Common draw functionality (most card games draw in similar ways)
func (g *BaseGame) DrawFromDeck(player *Player) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if len(g.Deck) == 0 {
		if len(g.DiscardPile) <= 1 {
			return fmt.Errorf("no cards left in deck and not enough cards to reshuffle")
		}
		// Basic reshuffling logic - specific games can override this
		topCard := g.DiscardPile[0]
		g.Deck = append([]types.Card{}, g.DiscardPile[1:]...)
		g.DiscardPile = []types.Card{topCard}
	}

	// Basic draw logic - games can override with specific constraints
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	player.AddCardToHand(card)
	return nil
}

func (g *BaseGame) DrawFromDiscard(player *Player) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	card, err := g.PopFromDiscardPile()
	if err != nil {
		return err
	}

	player.AddCardToHand(card)
	return nil
}

func (g *BaseGame) DiscardFromHand(player *Player, card types.Card) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if player.RemoveCardFromHand(card) != nil {
		g.PushToDiscardPile(card)
	}
	return nil
}
