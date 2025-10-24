package balatro

import (
	"math/rand"

	"github.com/charmbracelet/lipgloss"
)

// cards.go - Card System
// Purpose: Card types, deck management, and rendering helpers

// Suit represents the four card suits
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

// Symbol returns the Unicode symbol for the suit
func (s Suit) Symbol() string {
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

// String returns the full name of the suit
func (s Suit) String() string {
	switch s {
	case Spades:
		return "SPADES"
	case Hearts:
		return "HEARTS"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	default:
		return "UNKNOWN"
	}
}

// Color returns the display color for the suit
func (s Suit) Color() lipgloss.Color {
	switch s {
	case Spades:
		return lipgloss.Color("51") // Cyan
	case Hearts:
		return lipgloss.Color("201") // Magenta
	case Diamonds:
		return lipgloss.Color("214") // Orange
	case Clubs:
		return lipgloss.Color("46") // Green
	default:
		return lipgloss.Color("15") // White
	}
}

// Rank represents card ranks (Ace through King)
type Rank int

const (
	Ace Rank = 1 + iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// ShortString returns the short representation (A, 2, 3, ..., K)
func (r Rank) ShortString() string {
	switch r {
	case Ace:
		return "A"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		// For 2-9
		return string(rune('0' + r))
	}
}

// String returns the full name
func (r Rank) String() string {
	switch r {
	case Ace:
		return "ACE"
	case Ten:
		return "TEN"
	case Jack:
		return "JACK"
	case Queen:
		return "QUEEN"
	case King:
		return "KING"
	default:
		// For 2-9
		return r.ShortString()
	}
}

// ChipValue returns the base chip value for this rank
func (r Rank) ChipValue() int {
	switch r {
	case Ace:
		return 11
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten:
		return 10
	case Jack:
		return 10
	case Queen:
		return 10
	case King:
		return 10
	default:
		return 0
	}
}

// Enhancement represents card enhancements that modify scoring
type Enhancement int

const (
	NoEnhancement Enhancement = iota
	Bonus                     // +30 chips
	Mult                      // +4 mult
	Glass                     // x2 mult, destroys when scored
	Steel                     // +50 chips while in hand
	Stone                     // +50 chips, no rank
	Gold                      // $3 at end of round
	Lucky                     // 1 in 5 chance +20 mult
)

// ShortString returns abbreviated enhancement name
func (e Enhancement) ShortString() string {
	switch e {
	case NoEnhancement:
		return ""
	case Bonus:
		return "+30"
	case Mult:
		return "+4M"
	case Glass:
		return "GLS"
	case Steel:
		return "STL"
	case Stone:
		return "STN"
	case Gold:
		return "GLD"
	case Lucky:
		return "LCK"
	default:
		return "?"
	}
}

// String returns full enhancement name
func (e Enhancement) String() string {
	switch e {
	case NoEnhancement:
		return "None"
	case Bonus:
		return "Bonus"
	case Mult:
		return "Mult"
	case Glass:
		return "Glass"
	case Steel:
		return "Steel"
	case Stone:
		return "Stone"
	case Gold:
		return "Gold"
	case Lucky:
		return "Lucky"
	default:
		return "Unknown"
	}
}

// Description returns what the enhancement does
func (e Enhancement) Description() string {
	switch e {
	case NoEnhancement:
		return "No enhancement"
	case Bonus:
		return "+30 chips when scored"
	case Mult:
		return "+4 mult when scored"
	case Glass:
		return "x2 mult, destroys when scored"
	case Steel:
		return "+50 chips while in hand"
	case Stone:
		return "+50 chips, no rank/suit"
	case Gold:
		return "$3 at end of round"
	case Lucky:
		return "1 in 5 chance for +20 mult"
	default:
		return "Unknown"
	}
}

// BorderColor returns the border color for this enhancement
func (e Enhancement) BorderColor() lipgloss.Color {
	switch e {
	case NoEnhancement:
		return lipgloss.Color("240") // Gray
	case Bonus:
		return lipgloss.Color("226") // Gold
	case Mult:
		return lipgloss.Color("201") // Magenta
	case Glass:
		return lipgloss.Color("51") // Cyan
	case Steel:
		return lipgloss.Color("245") // Silver
	case Stone:
		return lipgloss.Color("130") // Brown
	case Gold:
		return lipgloss.Color("220") // Bright Gold
	case Lucky:
		return lipgloss.Color("46") // Green
	default:
		return lipgloss.Color("240")
	}
}

// TextColor returns the text color for enhancement labels
func (e Enhancement) TextColor() lipgloss.Color {
	switch e {
	case Bonus:
		return lipgloss.Color("226")
	case Mult:
		return lipgloss.Color("201")
	case Glass:
		return lipgloss.Color("51")
	case Steel:
		return lipgloss.Color("250")
	case Stone:
		return lipgloss.Color("130")
	case Gold:
		return lipgloss.Color("220")
	case Lucky:
		return lipgloss.Color("46")
	default:
		return lipgloss.Color("250")
	}
}

// Edition represents special card editions
type Edition int

const (
	NoEdition Edition = iota
	Foil           // +50 chips
	Holographic    // +10 mult
	Polychrome     // x1.5 mult
)

// String returns edition name
func (e Edition) String() string {
	switch e {
	case NoEdition:
		return "None"
	case Foil:
		return "Foil"
	case Holographic:
		return "Holographic"
	case Polychrome:
		return "Polychrome"
	default:
		return "Unknown"
	}
}

// ShortString returns abbreviated edition
func (e Edition) ShortString() string {
	switch e {
	case NoEdition:
		return ""
	case Foil:
		return "F"
	case Holographic:
		return "H"
	case Polychrome:
		return "P"
	default:
		return "?"
	}
}

// Seal represents card seals
type Seal int

const (
	NoSeal Seal = iota
	RedSeal    // Retrigger card
	BlueSeal   // Create Planet card if held
	GoldSeal   // $3 when played
	PurpleSeal // Create Tarot when discarded
)

// String returns seal name
func (s Seal) String() string {
	switch s {
	case NoSeal:
		return "None"
	case RedSeal:
		return "Red"
	case BlueSeal:
		return "Blue"
	case GoldSeal:
		return "Gold"
	case PurpleSeal:
		return "Purple"
	default:
		return "Unknown"
	}
}

// Symbol returns seal emoji/symbol
func (s Seal) Symbol() string {
	switch s {
	case NoSeal:
		return ""
	case RedSeal:
		return "🔴"
	case BlueSeal:
		return "🔵"
	case GoldSeal:
		return "🟡"
	case PurpleSeal:
		return "🟣"
	default:
		return "?"
	}
}

// Card represents a single playing card
type Card struct {
	Suit        Suit
	Rank        Rank
	Enhancement Enhancement
	Edition     Edition
	Seal        Seal
	Selected    bool // For UI selection
	Played      bool // Marked for playing
}

// NewDeck creates a standard 52-card deck
func NewDeck() []Card {
	deck := make([]Card, 0, 52)

	for suit := Clubs; suit <= Spades; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{
				Suit:        suit,
				Rank:        rank,
				Enhancement: NoEnhancement,
				Edition:     NoEdition,
				Seal:        NoSeal,
				Selected:    false,
				Played:      false,
			})
		}
	}

	return deck
}

// Shuffle shuffles a deck using Fisher-Yates algorithm
func Shuffle(deck []Card) {
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// GetChipValue returns the total chip value of a card
// Includes rank value + enhancement bonuses
func (c Card) GetChipValue() int {
	chips := c.Rank.ChipValue()

	// Add enhancement chips
	switch c.Enhancement {
	case Bonus:
		chips += 30
	case Steel:
		chips += 50
	case Stone:
		chips += 50
	}

	// Add edition chips
	if c.Edition == Foil {
		chips += 50
	}

	return chips
}

// ShortDisplay returns a compact string representation
func (c Card) ShortDisplay() string {
	return c.Rank.ShortString() + c.Suit.Symbol()
}

// Copy creates a deep copy of the card
func (c Card) Copy() Card {
	return Card{
		Suit:        c.Suit,
		Rank:        c.Rank,
		Enhancement: c.Enhancement,
		Edition:     c.Edition,
		Seal:        c.Seal,
		Selected:    c.Selected,
		Played:      c.Played,
	}
}

// SortCardsBySuit sorts cards by suit (Clubs, Diamonds, Hearts, Spades), then by rank within each suit
func SortCardsBySuit(cards []Card) []Card {
	sorted := make([]Card, len(cards))
	copy(sorted, cards)

	// Sort using a simple bubble sort for clarity
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			// Compare by suit first
			if sorted[j].Suit > sorted[j+1].Suit {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			} else if sorted[j].Suit == sorted[j+1].Suit {
				// Within same suit, sort by rank
				if sorted[j].Rank > sorted[j+1].Rank {
					sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
				}
			}
		}
	}

	return sorted
}

// SortCardsByRank sorts cards by rank (A, K, Q, J, 10, ..., 2), then by suit within each rank
func SortCardsByRank(cards []Card) []Card {
	sorted := make([]Card, len(cards))
	copy(sorted, cards)

	// Sort using a simple bubble sort for clarity
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			// Compare by rank first (descending - high cards first)
			// Special case: Ace (1) should be treated as highest
			rankJ := sorted[j].Rank
			rankJ1 := sorted[j+1].Rank
			if rankJ == Ace {
				rankJ = King + 1 // Treat Ace as higher than King
			}
			if rankJ1 == Ace {
				rankJ1 = King + 1
			}

			if rankJ < rankJ1 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			} else if rankJ == rankJ1 {
				// Within same rank, sort by suit
				if sorted[j].Suit > sorted[j+1].Suit {
					sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
				}
			}
		}
	}

	return sorted
}
