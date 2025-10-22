package balatro

import "fmt"

// Model.go - Model Management
// Purpose: Model initialization and layout calculations
// When to extend: Add new initialization logic or layout calculation functions here

// New creates a new Balatro game instance
func New() Model {
	cfg := loadConfig()

	return Model{
		config:           cfg,
		width:            0,
		height:           0,
		state:            StateLanding,
		focusedComponent: "landing",
		statusMsg:        "Balatro TUI - Press ? for help",
		landingPage:      nil, // Will be initialized on first resize
	}
}

// setSize updates the Model dimensions and recalculates layouts
func (m *Model) setSize(width, height int) {
	m.width = width
	m.height = height

	// Initialize or resize landing page
	if m.landingPage == nil {
		m.landingPage = NewLandingPage(width, height)
	} else {
		m.landingPage.Resize(width, height)
	}

	// Recalculate any layout-dependent values here
	// Example:
	// m.viewportHeight = height - 4 // account for title and status bars
	// m.maxVisible = m.viewportHeight - 2
}

// calculateLayout computes layout dimensions based on config
func (m Model) calculateLayout() (int, int) {
	contentWidth := m.width
	contentHeight := m.height

	// Adjust for UI elements
	if m.config.UI.ShowTitle {
		contentHeight -= 2 // title bar height
	}
	if m.config.UI.ShowStatus {
		contentHeight -= 1 // status bar height
	}

	return contentWidth, contentHeight
}

// calculateDualPaneLayout computes left and right pane widths
func (m Model) calculateDualPaneLayout() (int, int) {
	contentWidth, _ := m.calculateLayout()

	dividerWidth := 0
	if m.config.Layout.ShowDivider {
		dividerWidth = 1
	}

	leftWidth := int(float64(contentWidth-dividerWidth) * m.config.Layout.SplitRatio)
	rightWidth := contentWidth - leftWidth - dividerWidth

	return leftWidth, rightWidth
}

// Helper functions for common operations

// getContentArea returns the available content area dimensions
func (m Model) getContentArea() (width, height int) {
	return m.calculateLayout()
}

// isValidSize checks if the terminal size is sufficient
func (m Model) isValidSize() bool {
	return m.width >= 40 && m.height >= 10
}

// newGame initializes a new game state
func (m *Model) newGame() {
	// Create and shuffle deck
	m.deck = NewDeck()
	Shuffle(m.deck)

	// Deal starting hand (8 cards)
	m.hand = make([]Card, 8)
	copy(m.hand, m.deck[:8])
	m.deck = m.deck[8:]

	// Initialize game state
	m.playedCards = make([]Card, 0)
	m.selectedCardIndex = 0
	m.lastScore = ScoreCalculation{}
	m.currentHandInfo = HandInfo{}

	// Phase 3: Initialize with starting joker for testing
	m.jokers = make([]Joker, 0, 5) // Max 5 jokers
	// Give player basic "Joker" (+4 Mult) to test the system
	if basicJoker := GetJokerByID("joker_basic"); basicJoker != nil {
		m.jokers = append(m.jokers, *basicJoker)
	}

	// Phase 2: Initialize roundState with Ante 1, Small Blind
	m.roundState = NewRound(1)
	m.gamePhase = PhaseSelectCards

	// Legacy fields (kept for backward compatibility during transition)
	m.money = m.roundState.Money
	m.currentRound = 1
	m.handsRemaining = m.roundState.HandsRemaining
	m.discardsRemaining = m.roundState.DiscardsRemaining
	m.targetScore = m.roundState.TargetScore
	m.currentScore = m.roundState.CurrentScore

	// Switch to game state
	m.state = StateGame
	m.statusMsg = fmt.Sprintf("Ante 1 - %s: Reach %d chips to win!",
		m.roundState.CurrentBlind.Name,
		m.roundState.TargetScore)
}

// toggleCardForPlay toggles a card to be played
func (m *Model) toggleCardForPlay(index int) {
	if index < 0 || index >= len(m.hand) {
		return
	}

	card := m.hand[index]

	// Check if already in played cards
	for i, pc := range m.playedCards {
		if pc.Suit == card.Suit && pc.Rank == card.Rank {
			// Remove from played cards
			m.playedCards = append(m.playedCards[:i], m.playedCards[i+1:]...)
			m.updateCurrentHandInfo()
			return
		}
	}

	// Add to played cards (max 5)
	if len(m.playedCards) < 5 {
		m.playedCards = append(m.playedCards, card)
		m.updateCurrentHandInfo()
	}
}

// updateCurrentHandInfo updates the hand info based on currently selected cards
func (m *Model) updateCurrentHandInfo() {
	if len(m.playedCards) == 0 {
		m.currentHandInfo = HandInfo{}
		return
	}

	if len(m.playedCards) == 5 {
		m.currentHandInfo = EvaluateHand(m.playedCards)
	} else {
		// Show best possible hand with current selection
		_, info := FindBestPlay(m.playedCards)
		m.currentHandInfo = info
	}
}

// playHand plays the currently selected cards
func (m *Model) playHand() {
	if len(m.playedCards) == 0 {
		m.statusMsg = "Select cards to play!"
		return
	}

	// Phase 2: Check roundState for hands remaining
	if m.roundState.HandsRemaining <= 0 {
		m.statusMsg = "No hands remaining!"
		return
	}

	// Ensure we have exactly 5 cards
	var cardsToScore []Card
	if len(m.playedCards) == 5 {
		cardsToScore = m.playedCards
	} else {
		// Auto-complete to 5 cards if less
		cardsToScore, _ = FindBestPlay(m.hand)
	}

	// Evaluate and score (with joker effects)
	handInfo := EvaluateHand(cardsToScore)
	m.lastScore = CalculateScore(handInfo, cardsToScore, m.jokers)

	// Remove played cards from hand
	m.hand = removeCards(m.hand, cardsToScore)

	// Apply glass destruction
	m.hand = ApplyGlassDestruction(m.hand, cardsToScore)

	// Draw new cards to get back to 8
	cardsNeeded := 8 - len(m.hand)
	if cardsNeeded > 0 && cardsNeeded <= len(m.deck) {
		m.hand = append(m.hand, m.deck[:cardsNeeded]...)
		m.deck = m.deck[cardsNeeded:]
	}

	// Update state
	m.playedCards = make([]Card, 0)
	m.selectedCardIndex = 0

	// Phase 2: Update roundState with score
	blindWon, gameOver := m.roundState.PlayHand(m.lastScore.FinalScore)

	// Sync legacy fields
	m.currentScore = m.roundState.CurrentScore
	m.handsRemaining = m.roundState.HandsRemaining

	// Handle phase transitions
	if blindWon {
		// Blind complete! Show victory screen
		m.gamePhase = PhaseBlindComplete
		m.statusMsg = fmt.Sprintf("Blind beaten! Score: %d/%d", m.roundState.CurrentScore, m.roundState.TargetScore)
	} else if gameOver {
		// Out of hands - game over
		m.gamePhase = PhaseGameOver
		m.statusMsg = "Out of hands - Game Over!"
	} else {
		// Continue playing
		m.statusMsg = fmt.Sprintf("Scored %d points! (%d/%d)",
			m.lastScore.FinalScore,
			m.roundState.CurrentScore,
			m.roundState.TargetScore)
	}
}

// discardCards discards the selected cards
func (m *Model) discardCards() {
	if len(m.playedCards) == 0 {
		m.statusMsg = "Select cards to discard!"
		return
	}

	// Phase 2: Check roundState for discards remaining
	if m.roundState.DiscardsRemaining <= 0 {
		m.statusMsg = "No discards remaining!"
		return
	}

	// Remove selected cards from hand
	m.hand = removeCards(m.hand, m.playedCards)

	// Draw new cards
	cardsNeeded := len(m.playedCards)
	if cardsNeeded <= len(m.deck) {
		m.hand = append(m.hand, m.deck[:cardsNeeded]...)
		m.deck = m.deck[cardsNeeded:]
	}

	// Update state
	m.roundState.Discard()
	m.discardsRemaining = m.roundState.DiscardsRemaining // Sync legacy field
	m.playedCards = make([]Card, 0)
	m.selectedCardIndex = 0
	m.statusMsg = fmt.Sprintf("Discarded %d cards (%d discards left)",
		cardsNeeded,
		m.roundState.DiscardsRemaining)
}

// removeCards removes specified cards from hand
func removeCards(hand []Card, toRemove []Card) []Card {
	result := make([]Card, 0, len(hand))

	for _, card := range hand {
		shouldRemove := false
		for _, remove := range toRemove {
			if card.Suit == remove.Suit && card.Rank == remove.Rank {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			result = append(result, card)
		}
	}

	return result
}
