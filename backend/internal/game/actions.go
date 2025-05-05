package game

import (
	"backend/internal/models"
	"fmt"
)

/*
Functions I need:
	DrawFromDeck()
	DrawFromDiscard()
	Discard(Card card)
	GoOut() -- Need figure out if i want this to always calculate based on hand structure or if they press a button to try to go out
*/

func (g *models.Game) DrawFromDeck(player *models.Player) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if len(g.Deck) == 0 {
		//reshuffle the discard pile into the Deck, leaving the top card of the discard pile
		// need this to work so there is never a delay for the user
	}
	if len(player.Hand) >= g.Round+1 {
		return fmt.Errorf("player cannot draw more cards this round")
	}

	//actual logic
	card := g.Deck[0]
	g.Deck = g.Deck[1:]

	//checks the handsize making sure that players can only add if they have the same
	//amount of cards as the round so they cannot go over
	if len(player.Hand) == (g.Round) {
		player.AddCardToHand(card)
	}
	return nil
}

func (g *models.Game) DrawFromDiscard(player *models.Player) error {
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

func (g *models.Game) DiscardFromHand(player *models.Player, card models.Card) error {
	found := g.IsPlayerInGame(player)
	if !found {
		return fmt.Errorf("player %s not in game", player.PlayerID)
	}

	if player.RemoveCardFromHand(card) != nil {
		g.PushToDiscardPile(card)
	}
	return nil
}
