package solitaire

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SuitSymbol returns the Unicode symbol for the suit
func (s Suit) SuitSymbol() string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		return "?"
	}
}

// Color returns the color of the suit
func (s Suit) Color() lipgloss.Color {
	if s == Hearts || s == Diamonds {
		return lipgloss.Color("#FF0000") // Red
	}
	return lipgloss.Color("#000000") // Black
}

// IsRed returns true if the suit is red
func (s Suit) IsRed() bool {
	return s == Hearts || s == Diamonds
}

// RankSymbol returns the string representation of the rank
func (r Rank) RankSymbol() string {
	switch r {
	case Ace:
		return "A"
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "?"
	}
}

// Value returns the numeric value of the rank (1-13)
func (r Rank) Value() int {
	return int(r)
}

// NewDeck creates a standard 52-card deck
func NewDeck() []Card {
	deck := make([]Card, 0, 52)

	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{
				Suit:   suit,
				Rank:   rank,
				FaceUp: false,
			})
		}
	}

	return deck
}

// ShuffleDeck shuffles a deck of cards using Fisher-Yates algorithm
func ShuffleDeck(deck []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]Card, len(deck))
	copy(shuffled, deck)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

// CanStackOn checks if this card can be stacked on another in tableau
// (alternating colors, descending rank)
func (c Card) CanStackOn(other Card) bool {
	if !other.FaceUp {
		return false
	}

	// Must be alternating colors
	if c.Suit.IsRed() == other.Suit.IsRed() {
		return false
	}

	// Must be one rank lower
	return c.Rank.Value() == other.Rank.Value()-1
}

// CanMoveToFoundation checks if this card can be moved to a foundation pile
func (c Card) CanMoveToFoundation(foundation []Card) bool {
	if len(foundation) == 0 {
		return c.Rank == Ace
	}

	topCard := foundation[len(foundation)-1]

	// Must be same suit
	if c.Suit != topCard.Suit {
		return false
	}

	// Must be one rank higher
	return c.Rank.Value() == topCard.Rank.Value()+1
}

// IsKing returns true if the card is a King
func (c Card) IsKing() bool {
	return c.Rank == King
}
