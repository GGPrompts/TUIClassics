package balatro

import (
	"fmt"
	"math/rand"
)

// scoring.go - Scoring System
// Purpose: Calculate scores from poker hands with chips and multipliers

// ScoreCalculation stores the breakdown of a score calculation
type ScoreCalculation struct {
	// Chips
	BaseChips     int
	ChipModifiers []string // Human-readable modifiers like "+30 (Bonus)"
	TotalChips    int

	// Mult
	BaseMult      int
	MultModifiers []string // Human-readable modifiers like "+4 (Mult card)"
	TotalMult     int

	// Final
	FinalScore int // Chips × Mult
}

// CalculateScore calculates the score for a played hand
// This is the core scoring engine - chips × mult
func CalculateScore(handInfo HandInfo, playedCards []Card) ScoreCalculation {
	calc := ScoreCalculation{
		BaseChips:     handInfo.BaseChips,
		BaseMult:      handInfo.BaseMult,
		ChipModifiers: make([]string, 0),
		MultModifiers: make([]string, 0),
	}

	// Start with base values
	calc.TotalChips = calc.BaseChips
	calc.TotalMult = calc.BaseMult

	// Add chips from played cards
	for _, card := range playedCards {
		// Base rank chips
		rankChips := card.Rank.ChipValue()
		calc.TotalChips += rankChips
		calc.ChipModifiers = append(calc.ChipModifiers,
			fmt.Sprintf("+%d (%s)", rankChips, card.Rank.ShortString()))

		// Enhancement chip bonuses
		switch card.Enhancement {
		case Bonus:
			calc.TotalChips += 30
			calc.ChipModifiers = append(calc.ChipModifiers, "+30 (Bonus)")

		case Steel:
			// Steel gives +50 chips (while in hand, so also when played)
			calc.TotalChips += 50
			calc.ChipModifiers = append(calc.ChipModifiers, "+50 (Steel)")

		case Stone:
			// Stone gives +50 chips but no rank
			calc.TotalChips += 50
			calc.ChipModifiers = append(calc.ChipModifiers, "+50 (Stone)")
			// Remove the rank chips we added earlier
			calc.TotalChips -= rankChips
			// Update modifier to show it replaced rank chips
			calc.ChipModifiers[len(calc.ChipModifiers)-2] = fmt.Sprintf("+0 (%s - Stone)", card.Rank.ShortString())
		}

		// Edition chip bonuses
		if card.Edition == Foil {
			calc.TotalChips += 50
			calc.ChipModifiers = append(calc.ChipModifiers, "+50 (Foil)")
		}

		// Enhancement mult bonuses
		switch card.Enhancement {
		case Mult:
			calc.TotalMult += 4
			calc.MultModifiers = append(calc.MultModifiers, "+4 (Mult)")

		case Lucky:
			// 1 in 5 chance for +20 mult
			if rand.Intn(5) == 0 {
				calc.TotalMult += 20
				calc.MultModifiers = append(calc.MultModifiers, "+20 (Lucky!)")
			}
		}

		// Edition mult bonuses
		if card.Edition == Holographic {
			calc.TotalMult += 10
			calc.MultModifiers = append(calc.MultModifiers, "+10 (Holo)")
		}
	}

	// Apply multiplicative modifiers (Glass, Polychrome)
	multMultiplier := 1.0
	for _, card := range playedCards {
		if card.Enhancement == Glass {
			multMultiplier *= 2.0
			calc.MultModifiers = append(calc.MultModifiers, "x2 (Glass)")
		}

		if card.Edition == Polychrome {
			multMultiplier *= 1.5
			calc.MultModifiers = append(calc.MultModifiers, "x1.5 (Poly)")
		}
	}

	// Apply mult multiplier
	if multMultiplier != 1.0 {
		calc.TotalMult = int(float64(calc.TotalMult) * multMultiplier)
	}

	// Final score = Chips × Mult
	calc.FinalScore = calc.TotalChips * calc.TotalMult

	return calc
}

// FormatScore returns a formatted score breakdown string
func (sc ScoreCalculation) FormatScore() string {
	result := fmt.Sprintf("%d chips × %d mult = %d\n\n",
		sc.TotalChips, sc.TotalMult, sc.FinalScore)

	result += "Chips:\n"
	result += fmt.Sprintf("  Base: %d\n", sc.BaseChips)
	for _, mod := range sc.ChipModifiers {
		result += fmt.Sprintf("  %s\n", mod)
	}
	result += fmt.Sprintf("  Total: %d\n\n", sc.TotalChips)

	result += "Mult:\n"
	result += fmt.Sprintf("  Base: %d\n", sc.BaseMult)
	for _, mod := range sc.MultModifiers {
		result += fmt.Sprintf("  %s\n", mod)
	}
	result += fmt.Sprintf("  Total: %d\n", sc.TotalMult)

	return result
}

// ApplyGlassDestruction removes Glass cards from the hand after scoring
// Glass cards are destroyed when scored (x2 mult but fragile)
func ApplyGlassDestruction(hand []Card, playedCards []Card) []Card {
	result := make([]Card, 0, len(hand))

	// Create set of played glass cards to remove
	glassToRemove := make(map[int]bool)
	for i, card := range playedCards {
		if card.Enhancement == Glass {
			glassToRemove[i] = true
		}
	}

	// Keep only non-glass or non-played cards
	playedIndex := 0
	for _, card := range hand {
		// Check if this card was played and is glass
		isPlayedGlass := false
		if playedIndex < len(playedCards) {
			if card.Suit == playedCards[playedIndex].Suit &&
				card.Rank == playedCards[playedIndex].Rank &&
				card.Enhancement == Glass {
				isPlayedGlass = true
				playedIndex++
			}
		}

		if !isPlayedGlass {
			result = append(result, card)
		}
	}

	return result
}

// CalculateGoldEarnings calculates money earned from Gold cards
// Gold cards give $3 at end of round if held
func CalculateGoldEarnings(hand []Card) int {
	money := 0
	for _, card := range hand {
		if card.Enhancement == Gold {
			money += 3
		}
	}
	return money
}
