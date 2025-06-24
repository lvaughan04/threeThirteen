package models

import (
	"backend/internal/types"
	"backend/internal/utils"
	"fmt"
	"sort"
)

// ThreeThirteenGame extends the BaseGame with Three Thirteen specific logic
type ThreeThirteenGame struct {
	*BaseGame
	MaxRounds   int         `json:"max_rounds"`   // (3 to 13)
	PointValues map[int]int `json:"point_values"` // Card point values based on card.Value
}

// NewThreeThirteenGame creates a new Three Thirteen game instance
func NewThreeThirteenGame(gameID string) *ThreeThirteenGame {
	game := &ThreeThirteenGame{
		BaseGame: &BaseGame{
			GameID:    gameID,
			GameState: "INITIALIZED",
		},
		MaxRounds: 11,
	}

	// Initialize point values for scoring
	game.initPointValues()

	return game
}

// Initialize card point values for Three Thirteen scoring
func (g *ThreeThirteenGame) initPointValues() {
	g.PointValues = make(map[int]int)

	// All cards are worth their face value
	for i := 1; i <= 13; i++ {
		g.PointValues[i] = i
	}
}

// InitializeGame sets up the game with players and initial state
func (g *ThreeThirteenGame) InitializeGame(players []Player) error {
	if len(players) < 2 || len(players) > 6 {
		return fmt.Errorf("three thirteen requires 2-6 players, got %d", len(players))
	}

	g.Players = players
	g.Round = 3 // Three Thirteen starts with 3 cards
	g.CurrentTurn = 0
	g.GameState = "initialized"

	// Create and shuffle deck
	g.Deck = utils.CreateDeck()
	utils.ShuffleDeck(g.Deck)
	g.DiscardPile = []types.Card{}

	return nil
}

// GetWildCardValue returns the current wild card value based on the round
func (g *ThreeThirteenGame) GetWildCardValue() int {
	return g.Round // The wild card value is the current round number
}

// StartRound begins a new round, dealing cards
func (g *ThreeThirteenGame) StartRound() error {
	if g.GameState != "initialized" && g.GameState != "round_complete" {
		return fmt.Errorf("game must be initialized or round complete to start a new round")
	}

	// Reset hands
	for i := range g.Players {
		g.Players[i].Hand = []types.Card{}
	}

	// If needed, create a new deck
	if len(g.Deck) < len(g.Players)*g.Round {
		g.Deck = utils.CreateDeck()
		utils.ShuffleDeck(g.Deck)
	}

	

	// Place first card in discard pile
	if len(g.Deck) > 0 {
		g.PushToDiscardPile(g.Deck[0])
		g.Deck = g.Deck[1:]
	} else {
		return fmt.Errorf("not enough cards to start discard pile")
	}

	g.GameState = "in_progress"
	return nil
}

// Override DrawFromDeck to add ThreeThirteen-specific rules
func (g *ThreeThirteenGame) DrawFromDeck(player *Player) error {
	// Check if it's player's turn
	if g.Players[g.CurrentTurn].PlayerID != player.PlayerID {
		return fmt.Errorf("not player's turn")
	}

	// ThreeThirteen-specific rule - can only draw at the beginning of turn
	if len(player.Hand) >= g.Round+1 {
		return fmt.Errorf("player cannot draw more cards this round")
	}

	// Call the base implementation
	return g.BaseGame.DrawFromDeck(player)
}

// Override DrawFromDiscard to add ThreeThirteen-specific rules
func (g *ThreeThirteenGame) DrawFromDiscard(player *Player) error {
	// Check if it's player's turn
	if g.Players[g.CurrentTurn].PlayerID != player.PlayerID {
		return fmt.Errorf("not player's turn")
	}

	// ThreeThirteen-specific rule - can only draw at the beginning of turn
	if len(player.Hand) >= g.Round+1 {
		return fmt.Errorf("player cannot draw more cards this round")
	}

	// Call the base implementation
	return g.BaseGame.DrawFromDiscard(player)
}

// DiscardAndEndTurn discards a card and advances to the next player
func (g *ThreeThirteenGame) DiscardAndEndTurn(player *Player, card types.Card) error {
	// Verify it's player's turn
	if g.Players[g.CurrentTurn].PlayerID != player.PlayerID {
		return fmt.Errorf("not player's turn")
	}

	// Must have drawn a card first
	if len(player.Hand) != g.Round+1 {
		return fmt.Errorf("must draw a card before discarding")
	}

	// Discard the card
	err := g.BaseGame.DiscardFromHand(player, card)
	if err != nil {
		return err
	}

	// Move to next player
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)
	return nil
}

// IsSet checks if a group of cards forms a valid set (same rank/value)
func (g *ThreeThirteenGame) IsSet(cards []types.Card) bool {
	if len(cards) < 3 {
		return false
	}

	wildCardValue := g.GetWildCardValue()
	targetValue := 0

	for _, card := range cards {
		// Skip wild cards, they can be anything
		if card.Value == wildCardValue {
			continue
		}

		if targetValue == 0 {
			targetValue = card.Value
		} else if targetValue != card.Value {
			return false
		}
	}

	return true
}

// IsRun checks if a group of cards forms a valid run (sequential, same suit)
func (g *ThreeThirteenGame) IsRun(cards []types.Card) bool {
	if len(cards) < 3 {
		return false
	}

	// Group non-wild cards by suit
	suitedCards := make(map[string][]types.Card)
	wildCards := 0
	wildCardValue := g.GetWildCardValue()

	for _, card := range cards {
		if card.Value == wildCardValue {
			wildCards++
		} else {
			suitedCards[card.Suit] = append(suitedCards[card.Suit], card)
		}
	}

	// Must be same suit (plus wild cards)
	if len(suitedCards) != 1 {
		return false
	}

	// Get the one suit we have
	var cardList []types.Card
	for _, c := range suitedCards {
		cardList = c
		break
	}

	// Sort cards by value
	sort.Slice(cardList, func(i, j int) bool {
		return cardList[i].Value < cardList[j].Value
	})

	// Check if they form a sequence, accounting for wild cards
	gaps := 0
	for i := 1; i < len(cardList); i++ {
		diff := cardList[i].Value - cardList[i-1].Value
		if diff == 1 {
			continue
		}
		gaps += diff - 1
	}

	// Wild cards can fill gaps
	return gaps <= wildCards
}

// CalculateScore calculates a player's score for unmatched cards
func (g *ThreeThirteenGame) CalculateScore(player *Player) int {
	score := 0
	for _, card := range player.Hand {
		score += g.PointValues[card.Value]
	}
	return score
}

// EndRound ends the current round, calculates scores, and prepares for next round
func (g *ThreeThirteenGame) EndRound() {
	// Calculate scores for each player
	for i := range g.Players {
		score := g.CalculateScore(&g.Players[i])
		g.Players[i].Score += score
	}

	// Advance to next round
	g.Round++

	if g.Round > g.MaxRounds {
		g.GameState = "game_over"
	} else {
		g.GameState = "round_complete"
	}
}

// DeclareGoingOut allows a player to go out if all their cards form valid sets or runs
func (g *ThreeThirteenGame) DeclareGoingOut(player *Player, meldGroups [][]types.Card) error {
	// Check if player is in the game
	if !g.IsPlayerInGame(player) {
		return fmt.Errorf("player not in game")
	}

	// Verify it's player's turn
	if g.Players[g.CurrentTurn].PlayerID != player.PlayerID {
		return fmt.Errorf("not player's turn")
	}

	// Check that all card groups are valid
	for _, group := range meldGroups {
		if !g.IsSet(group) && !g.IsRun(group) {
			return fmt.Errorf("invalid card group")
		}
	}

	// Verify all cards in player's hand are included in melds
	// This requires more complex validation that we can add later

	// End the round
	g.EndRound()
	return nil
}

// IsGameOver checks if the game has ended
func (g *ThreeThirteenGame) IsGameOver() bool {
	return g.GameState == "game_over"
}

// GetWinner returns the player with the lowest score (winner)
func (g *ThreeThirteenGame) GetWinner() *Player {
	if !g.IsGameOver() {
		return nil
	}

	lowestScore := g.Players[0].Score
	winnerIndex := 0

	for i, player := range g.Players {
		if player.Score < lowestScore {
			lowestScore = player.Score
			winnerIndex = i
		}
	}

	return &g.Players[winnerIndex]
}


// GO THROUGH THIS CODE AND MAKE SURE IT IS ALL CORRECT TO THE EXPECTATIONS OF THE GAME