package balatro

// round.go
// Purpose: Round state management for game progression
// Tracks ante, blind progress, scores, and win/loss conditions

// RoundState tracks the current state of game progression
type RoundState struct {
	Ante              int
	CurrentBlind      Blind
	BlindProgress     int // 0=Small, 1=Big, 2=Boss
	CurrentScore      int
	TargetScore       int
	HandsRemaining    int
	DiscardsRemaining int
	Money             int // Player's current money
}

// Starting values for each blind
const (
	StartingHands    = 4
	StartingDiscards = 3
	StartingMoney    = 4 // Start with $4 like Balatro
)

// NewRound creates a new round starting at the given ante
func NewRound(ante int) RoundState {
	blind := GetBlind(ante, SmallBlind)

	return RoundState{
		Ante:              ante,
		CurrentBlind:      blind,
		BlindProgress:     0, // Start at Small Blind
		CurrentScore:      0,
		TargetScore:       blind.TargetScore,
		HandsRemaining:    StartingHands,
		DiscardsRemaining: StartingDiscards,
		Money:             StartingMoney,
	}
}

// PlayHand processes playing a hand and returns whether the blind was won and if game is over
func (r *RoundState) PlayHand(score int) (blindWon bool, gameOver bool) {
	r.CurrentScore += score
	r.HandsRemaining--

	// Check if we've reached the target score
	if r.CurrentScore >= r.TargetScore {
		return true, false // Blind won!
	}

	// Check if we've run out of hands
	if r.HandsRemaining <= 0 {
		return false, true // Game over - lost!
	}

	// Still playing
	return false, false
}

// Discard processes a discard action
func (r *RoundState) Discard() {
	r.DiscardsRemaining--
}

// AdvanceBlind moves to the next blind (Small → Big → Boss → Next Ante)
func (r *RoundState) AdvanceBlind() {
	// Award money for beating the blind
	r.Money += r.CurrentBlind.Reward

	r.BlindProgress++

	// Check if we've completed all three blinds
	if r.BlindProgress > 2 {
		// Move to next ante
		r.Ante++
		r.BlindProgress = 0
	}

	// Get the new blind
	blindType := BlindType(r.BlindProgress)
	r.CurrentBlind = GetBlind(r.Ante, blindType)
	r.TargetScore = r.CurrentBlind.TargetScore

	// Reset for new blind
	r.ResetForNewBlind()
}

// ResetForNewBlind resets hands, discards, and score for a new blind
func (r *RoundState) ResetForNewBlind() {
	r.CurrentScore = 0
	r.HandsRemaining = StartingHands
	r.DiscardsRemaining = StartingDiscards
}

// GetBlindName returns the name of the current blind
func (r *RoundState) GetBlindName() string {
	return r.CurrentBlind.Name
}

// GetProgress returns current blind progress as a string
func (r *RoundState) GetProgress() string {
	switch r.BlindProgress {
	case 0:
		return "Small Blind"
	case 1:
		return "Big Blind"
	case 2:
		return "Boss Blind"
	default:
		return "Unknown"
	}
}
