package types

type Card struct {
	Suit  string `json:"suit"` // Needs to be in the format of Heart, Diamond, Spade, or Club
	Value int    `json:"value"`
}