package utils

import (
	"backend/internal/types"
	"math/rand"
	
)

func CreateDeck() []types.Card {
	suits := []string{"Hearts", "Diamonds", "Spades", "Clubs"}
	var deck []types.Card
	for _, suit := range suits  { //the suits
		for i := 1 ; i < 14; i++{
			deck = append(deck, types.Card{Suit: suit, Value: i})
		}
	}
	return deck
}

//This function implements the Fischer - Yates shuffle
func ShuffleDeck(deck []types.Card) {
	n := len(deck)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1) // Generate a random index in the range [0, i]
		deck[i], deck[j] = deck[j], deck[i] // Swap elements
	}
}