package game

import (
	"testing"
	"backend/internal/models"
	"reflect"
)

func TestCreateDeck(t *testing.T) {
	deck := CreateDeck()

	if len(deck) != 52{
		t.Errorf("Expected 52 cards, got %d", len(deck))
	} 
	suitCount := make(map[string]int)

	for _, card := range deck {
		suitCount[card.Suit]++
	}

	expectedSuits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}

	// Check that each suit appears exactly 13 times
	for _, suit := range expectedSuits {
		if suitCount[suit] != 13 {
			t.Errorf("Expected 13 cards of suit %s, but got %d", suit, suitCount[suit])
		}
	}
}

func TestShuffleDeck(t *testing.T) {
	
	// Create a deck
	originalDeck := CreateDeck()
	shuffledDeck := make([]models.Card, len(originalDeck))
	copy(shuffledDeck, originalDeck) // Copy to compare after shuffling

	// Shuffle the deck
	ShuffleDeck(shuffledDeck)

	// Test 1: Deck should still have 52 cards
	if len(shuffledDeck) != 52 {
		t.Errorf("Expected 52 cards after shuffling, but got %d", len(shuffledDeck))
	}

	// Test 2: Ensure all original cards are still present (no missing or duplicates)
	cardMap := make(map[models.Card]int)
	for _, card := range originalDeck {
		cardMap[card]++
	}

	for _, card := range shuffledDeck {
		if cardMap[card] == 0 {
			t.Errorf("Card %v is missing after shuffle", card)
		} else {
			cardMap[card]--
		}
	}

	// Ensure no extra or missing cards
	for card, count := range cardMap {
		if count != 0 {
			t.Errorf("Card %v has incorrect count after shuffle", card)
		}
	}

	// Test 3: The order of the deck should be different (most of the time)
	// This test might occasionally fail if the shuffle results in the same order by chance.
	if reflect.DeepEqual(originalDeck, shuffledDeck) {
		t.Errorf("Deck order did not change after shuffling")
	}

}