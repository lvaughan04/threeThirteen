package game

import (
	"backend/internal/models"
	"math/rand"
	
)
/*
Functions I need:
	CreateDeck()
	ShuffleDeck()
*/

func CreateDeck() []models.Card {
	suits := []string{"Hearts", "Diamonds", "Spades", "Clubs"}
	var deck []models.Card
	for _, suit := range suits  { //the suits
		for i := 1 ; i < 14; i++{
			deck = append(deck, models.Card{Suit: suit, Value: i})
		}
	}
	return deck
}

//This function implements the Fischer - Yates shuffle
func ShuffleDeck(deck []models.Card) {
	n := len(deck)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1) // Generate a random index in the range [0, i]
		deck[i], deck[j] = deck[j], deck[i] // Swap elements
	}
}