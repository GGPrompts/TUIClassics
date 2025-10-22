package balatro

import (
	"sort"
)

// hands.go - Poker Hand Evaluation
// Purpose: Detect poker hands and find best plays

// PokerHand represents the type of poker hand
type PokerHand int

const (
	HighCard PokerHand = iota
	Pair
	TwoPair
	ThreeOfKind
	Straight
	Flush
	FullHouse
	FourOfKind
	StraightFlush
	RoyalFlush
)

// String returns the name of the poker hand
func (ph PokerHand) String() string {
	switch ph {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

// HandInfo stores detected hand details and scoring info
type HandInfo struct {
	Type      PokerHand
	Cards     []Card // The 5 cards that make the hand
	BaseChips int    // Base chips for this hand type
	BaseMult  int    // Base multiplier
	Level     int    // Hand level (from Planet cards, default 1)
}

// Base hand values from Balatro
var baseHandValues = map[PokerHand]HandInfo{
	HighCard:      {BaseChips: 5, BaseMult: 1, Level: 1},
	Pair:          {BaseChips: 10, BaseMult: 2, Level: 1},
	TwoPair:       {BaseChips: 20, BaseMult: 2, Level: 1},
	ThreeOfKind:   {BaseChips: 30, BaseMult: 3, Level: 1},
	Straight:      {BaseChips: 30, BaseMult: 4, Level: 1},
	Flush:         {BaseChips: 35, BaseMult: 4, Level: 1},
	FullHouse:     {BaseChips: 40, BaseMult: 4, Level: 1},
	FourOfKind:    {BaseChips: 60, BaseMult: 7, Level: 1},
	StraightFlush: {BaseChips: 100, BaseMult: 8, Level: 1},
	RoyalFlush:    {BaseChips: 100, BaseMult: 8, Level: 1},
}

// EvaluateHand detects the best poker hand from exactly 5 cards
func EvaluateHand(cards []Card) HandInfo {
	if len(cards) != 5 {
		// Default to high card if not exactly 5 cards
		return HandInfo{
			Type:      HighCard,
			Cards:     cards,
			BaseChips: 5,
			BaseMult:  1,
			Level:     1,
		}
	}

	// Check hands in order from best to worst
	if isRoyalFlush(cards) {
		info := baseHandValues[RoyalFlush]
		info.Type = RoyalFlush
		info.Cards = cards
		return info
	}

	if isStraightFlush(cards) {
		info := baseHandValues[StraightFlush]
		info.Type = StraightFlush
		info.Cards = cards
		return info
	}

	if isFourOfKind(cards) {
		info := baseHandValues[FourOfKind]
		info.Type = FourOfKind
		info.Cards = cards
		return info
	}

	if isFullHouse(cards) {
		info := baseHandValues[FullHouse]
		info.Type = FullHouse
		info.Cards = cards
		return info
	}

	if isFlush(cards) {
		info := baseHandValues[Flush]
		info.Type = Flush
		info.Cards = cards
		return info
	}

	if isStraight(cards) {
		info := baseHandValues[Straight]
		info.Type = Straight
		info.Cards = cards
		return info
	}

	if isThreeOfKind(cards) {
		info := baseHandValues[ThreeOfKind]
		info.Type = ThreeOfKind
		info.Cards = cards
		return info
	}

	if isTwoPair(cards) {
		info := baseHandValues[TwoPair]
		info.Type = TwoPair
		info.Cards = cards
		return info
	}

	if isPair(cards) {
		info := baseHandValues[Pair]
		info.Type = Pair
		info.Cards = cards
		return info
	}

	// High card
	info := baseHandValues[HighCard]
	info.Type = HighCard
	info.Cards = cards
	return info
}

// FindBestPlay finds the best 5-card poker hand from a larger hand
// Returns the best 5 cards and their hand info
func FindBestPlay(hand []Card) ([]Card, HandInfo) {
	if len(hand) <= 5 {
		info := EvaluateHand(hand)
		return hand, info
	}

	// Generate all 5-card combinations
	var bestCards []Card
	bestInfo := HandInfo{Type: HighCard, BaseChips: 0, BaseMult: 0}

	combinations := generateCombinations(hand, 5)
	for _, combo := range combinations {
		info := EvaluateHand(combo)

		// Compare hands (higher type wins, or higher chips if same type)
		if info.Type > bestInfo.Type ||
			(info.Type == bestInfo.Type && info.BaseChips > bestInfo.BaseChips) {
			bestInfo = info
			bestCards = make([]Card, len(combo))
			copy(bestCards, combo)
		}
	}

	return bestCards, bestInfo
}

// Helper functions for hand detection

// isRoyalFlush checks for A-K-Q-J-10 of same suit
func isRoyalFlush(cards []Card) bool {
	if !isStraightFlush(cards) {
		return false
	}

	// Check if it's specifically A-K-Q-J-10
	ranks := getRanks(cards)
	sort.Ints(ranks)

	return ranks[0] == int(Ten) &&
		ranks[1] == int(Jack) &&
		ranks[2] == int(Queen) &&
		ranks[3] == int(King) &&
		ranks[4] == int(Ace)
}

// isStraightFlush checks for straight AND flush
func isStraightFlush(cards []Card) bool {
	return isStraight(cards) && isFlush(cards)
}

// isFourOfKind checks for 4 cards of same rank
func isFourOfKind(cards []Card) bool {
	counts := rankCounts(cards)
	for _, count := range counts {
		if count == 4 {
			return true
		}
	}
	return false
}

// isFullHouse checks for 3 of a kind + pair
func isFullHouse(cards []Card) bool {
	counts := rankCounts(cards)
	hasThree := false
	hasTwo := false

	for _, count := range counts {
		if count == 3 {
			hasThree = true
		}
		if count == 2 {
			hasTwo = true
		}
	}

	return hasThree && hasTwo
}

// isFlush checks if all cards are same suit
func isFlush(cards []Card) bool {
	if len(cards) == 0 {
		return false
	}

	firstSuit := cards[0].Suit
	for _, card := range cards[1:] {
		if card.Suit != firstSuit {
			return false
		}
	}
	return true
}

// isStraight checks for 5 consecutive ranks
func isStraight(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}

	ranks := getRanks(cards)
	sort.Ints(ranks)

	// Check for regular straight
	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]+1 {
			// Not a regular straight, check for A-2-3-4-5 (wheel)
			if ranks[0] == int(Ace) &&
				ranks[1] == int(Two) &&
				ranks[2] == int(Three) &&
				ranks[3] == int(Four) &&
				ranks[4] == int(Five) {
				return true
			}

			// Check for high ace straight (10-J-Q-K-A)
			if ranks[0] == int(Ten) &&
				ranks[1] == int(Jack) &&
				ranks[2] == int(Queen) &&
				ranks[3] == int(King) &&
				ranks[4] == int(Ace) {
				return true
			}

			return false
		}
	}

	return true
}

// isThreeOfKind checks for 3 cards of same rank
func isThreeOfKind(cards []Card) bool {
	counts := rankCounts(cards)
	for _, count := range counts {
		if count == 3 {
			return true
		}
	}
	return false
}

// isTwoPair checks for 2 different pairs
func isTwoPair(cards []Card) bool {
	counts := rankCounts(cards)
	pairCount := 0

	for _, count := range counts {
		if count == 2 {
			pairCount++
		}
	}

	return pairCount == 2
}

// isPair checks for 2 cards of same rank
func isPair(cards []Card) bool {
	counts := rankCounts(cards)
	for _, count := range counts {
		if count == 2 {
			return true
		}
	}
	return false
}

// rankCounts returns a map of rank -> count
func rankCounts(cards []Card) map[Rank]int {
	counts := make(map[Rank]int)
	for _, card := range cards {
		counts[card.Rank]++
	}
	return counts
}

// getRanks extracts just the rank values from cards
func getRanks(cards []Card) []int {
	ranks := make([]int, len(cards))
	for i, card := range cards {
		ranks[i] = int(card.Rank)
	}
	return ranks
}

// generateCombinations generates all k-sized combinations from a slice
func generateCombinations(cards []Card, k int) [][]Card {
	var result [][]Card

	if k == 0 {
		return [][]Card{{}}
	}

	if len(cards) < k {
		return result
	}

	// Include first card
	for _, combo := range generateCombinations(cards[1:], k-1) {
		newCombo := append([]Card{cards[0]}, combo...)
		result = append(result, newCombo)
	}

	// Exclude first card
	result = append(result, generateCombinations(cards[1:], k)...)

	return result
}
