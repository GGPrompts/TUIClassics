package balatro

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
)

// jokers.go - Joker System
// Purpose: Joker cards that modify scoring with various effects

// JokerEffect represents the type of effect a joker has
type JokerEffect int

const (
	NoEffect JokerEffect = iota
	AddChips              // +N chips
	AddMult               // +N mult
	MultIfPair            // +N mult if hand contains pair or better
	MultIfSuit            // +N mult per card of specific suit
	ChipsPerCard          // +N chips per card of rank
	RetriggerFirst        // Retrigger first card scored
	MultPercentage        // Multiply mult by percentage (e.g., x1.5)
)

// Rarity represents joker rarity levels
type Rarity int

const (
	Common Rarity = iota
	Uncommon
	Rare
	Legendary
)

// String returns rarity name
func (r Rarity) String() string {
	switch r {
	case Common:
		return "Common"
	case Uncommon:
		return "Uncommon"
	case Rare:
		return "Rare"
	case Legendary:
		return "Legendary"
	default:
		return "Unknown"
	}
}

// Color returns display color for rarity
func (r Rarity) Color() lipgloss.Color {
	switch r {
	case Common:
		return lipgloss.Color("250") // Light Gray
	case Uncommon:
		return lipgloss.Color("46")  // Green
	case Rare:
		return lipgloss.Color("51")  // Cyan
	case Legendary:
		return lipgloss.Color("226") // Gold
	default:
		return lipgloss.Color("250")
	}
}

// Joker represents a joker card that modifies gameplay
type Joker struct {
	ID          string
	Name        string
	Description string
	Rarity      Rarity
	Effect      JokerEffect
	Value       int // Effect magnitude (e.g., +4 for AddMult)

	// Conditional parameters
	TargetSuit Suit // For MultIfSuit effect
	TargetRank Rank // For ChipsPerCard effect

	// State (for stateful jokers)
	Counter int
}

// Apply applies the joker's effect to a score calculation
// This is called during scoring, modifying the calculation in-place
func (j *Joker) Apply(calc *ScoreCalculation, handInfo HandInfo, playedCards []Card) {
	switch j.Effect {
	case AddChips:
		// Simple chip bonus
		calc.TotalChips += j.Value
		calc.ChipModifiers = append(calc.ChipModifiers,
			fmt.Sprintf("+%d (%s)", j.Value, j.Name))

	case AddMult:
		// Simple mult bonus
		calc.TotalMult += j.Value
		calc.MultModifiers = append(calc.MultModifiers,
			fmt.Sprintf("+%d (%s)", j.Value, j.Name))

	case MultIfPair:
		// Conditional mult if hand contains pair or better
		if handInfo.Type >= Pair {
			calc.TotalMult += j.Value
			calc.MultModifiers = append(calc.MultModifiers,
				fmt.Sprintf("+%d (%s)", j.Value, j.Name))
		}

	case MultIfSuit:
		// Add mult for each played card of specific suit
		count := 0
		for _, card := range playedCards {
			if card.Suit == j.TargetSuit {
				count++
			}
		}
		if count > 0 {
			totalBonus := j.Value * count
			calc.TotalMult += totalBonus
			calc.MultModifiers = append(calc.MultModifiers,
				fmt.Sprintf("+%d (%s - %d cards)", totalBonus, j.Name, count))
		}

	case ChipsPerCard:
		// Add chips for each played card of specific rank
		count := 0
		for _, card := range playedCards {
			if card.Rank == j.TargetRank {
				count++
			}
		}
		if count > 0 {
			totalBonus := j.Value * count
			calc.TotalChips += totalBonus
			calc.ChipModifiers = append(calc.ChipModifiers,
				fmt.Sprintf("+%d (%s - %d cards)", totalBonus, j.Name, count))
		}

	case RetriggerFirst:
		// Retrigger first played card (add its chip/mult contributions again)
		if len(playedCards) > 0 {
			card := playedCards[0]
			// Add chip value again
			chips := card.GetChipValue()
			calc.TotalChips += chips
			calc.ChipModifiers = append(calc.ChipModifiers,
				fmt.Sprintf("+%d (%s retrigger)", chips, j.Name))

			// Apply mult effects again if card has them
			if card.Enhancement == Mult {
				calc.TotalMult += 4
				calc.MultModifiers = append(calc.MultModifiers,
					fmt.Sprintf("+4 (%s retrigger)", j.Name))
			}
		}

	case MultPercentage:
		// Multiply total mult by percentage (e.g., Value=150 means x1.5)
		if j.Value > 100 {
			oldMult := calc.TotalMult
			calc.TotalMult = (calc.TotalMult * j.Value) / 100
			calc.MultModifiers = append(calc.MultModifiers,
				fmt.Sprintf("x%.2f (%s)", float64(j.Value)/100.0, j.Name))
			_ = oldMult // For debugging
		}
	}
}

// GetCost returns the purchase cost for this joker
func (j Joker) GetCost() int {
	switch j.Rarity {
	case Common:
		return 5
	case Uncommon:
		return 7
	case Rare:
		return 10
	case Legendary:
		return 20
	default:
		return 5
	}
}

// GetSellValue returns the sell value (half of cost)
func (j Joker) GetSellValue() int {
	return j.GetCost() / 2
}

// Starting jokers - simple ones to get the system working
var StartingJokers = []Joker{
	{
		ID:          "joker_basic",
		Name:        "Joker",
		Description: "+4 Mult",
		Rarity:      Common,
		Effect:      AddMult,
		Value:       4,
	},
	{
		ID:          "greedy_joker",
		Name:        "Greedy Joker",
		Description: "Played Diamond cards give +3 Mult each",
		Rarity:      Common,
		Effect:      MultIfSuit,
		Value:       3,
		TargetSuit:  Diamonds,
	},
	{
		ID:          "lusty_joker",
		Name:        "Lusty Joker",
		Description: "Played Heart cards give +3 Mult each",
		Rarity:      Common,
		Effect:      MultIfSuit,
		Value:       3,
		TargetSuit:  Hearts,
	},
	{
		ID:          "wrathful_joker",
		Name:        "Wrathful Joker",
		Description: "Played Spade cards give +3 Mult each",
		Rarity:      Common,
		Effect:      MultIfSuit,
		Value:       3,
		TargetSuit:  Spades,
	},
	{
		ID:          "gluttonous_joker",
		Name:        "Gluttonous Joker",
		Description: "Played Club cards give +3 Mult each",
		Rarity:      Common,
		Effect:      MultIfSuit,
		Value:       3,
		TargetSuit:  Clubs,
	},
}

// GetJokerByID retrieves a joker by its ID
func GetJokerByID(id string) *Joker {
	for i := range StartingJokers {
		if StartingJokers[i].ID == id {
			// Return a copy
			joker := StartingJokers[i]
			return &joker
		}
	}
	return nil
}

// GetRandomJoker returns a random joker of the specified max rarity
func GetRandomJoker(maxRarity Rarity) *Joker {
	// Filter jokers by rarity
	available := make([]Joker, 0)
	for _, joker := range StartingJokers {
		if joker.Rarity <= maxRarity {
			available = append(available, joker)
		}
	}

	if len(available) == 0 {
		return nil
	}

	// Return a copy of a random joker
	joker := available[randomInt(len(available))]
	return &joker
}

// Helper function for random int
func randomInt(n int) int {
	if n <= 0 {
		return 0
	}
	// This will use the global rand which should be seeded elsewhere
	return int(rand.Int63n(int64(n)))
}
