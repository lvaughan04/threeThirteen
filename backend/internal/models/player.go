package models

import ("fmt"
	"backend/internal/types"	
)


type Player struct {
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Hand     []types.Card `json:"hand"`
	IsTurn   bool   `json:"is_turn"`
	IsDealer bool   `json:"is_dealer"`
	IsActive bool   `json:"is_active"`
}

func (p *Player) AddCardToHand(card types.Card) {
	p.Hand = append(p.Hand, card)
}

// This function will take a card to be removed and will be
// sent from the client side to a controller method with the card that needs to be removed
func (p *Player) RemoveCardFromHand(target types.Card) error {
	for i, card := range p.Hand {
		if card == target {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return nil // success, no error
		}
	}
	return fmt.Errorf("card not found in hand: %+v", target)
}
