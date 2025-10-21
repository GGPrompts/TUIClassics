package solitaire

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseEvent handles mouse events
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.state != StatePlaying {
		return m, nil
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			return m.handleMousePress(msg)
		} else if msg.Action == tea.MouseActionRelease {
			return m.handleMouseRelease(msg)
		}

	case tea.MouseButtonRight:
		// Right-click for auto-move to foundation
		if msg.Action == tea.MouseActionPress {
			return m.handleRightClick(msg)
		}
	}

	return m, nil
}

// handleMousePress handles mouse button press (start drag or select)
func (m Model) handleMousePress(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Store press position to detect click vs drag
	m.mousePressX = msg.X
	m.mousePressY = msg.Y

	// Determine which pile was clicked
	location := m.getPileAtPosition(msg.X, msg.Y)
	if location == nil {
		return m, nil
	}

	// If stock was clicked, draw a card
	if location.PileType == StockPile {
		m.DrawFromStock()
		return m, nil
	}

	// Prepare card for potential drag or selection
	var card *Card
	var cards []Card
	var cardIndex int

	switch location.PileType {
	case TableauPile:
		pile := m.tableau[location.PileIndex]
		if len(pile.Cards) > 0 {
			// Check which card in the tableau was clicked (for sequences)
			cardIndex = m.getTableauCardIndex(location.PileIndex, msg.Y)
			if cardIndex >= 0 && cardIndex < len(pile.Cards) {
				if pile.Cards[cardIndex].FaceUp {
					sequence := m.GetMovableSequence(location.PileIndex, cardIndex)
					if sequence != nil {
						card = &pile.Cards[cardIndex]
						cards = sequence
					}
				}
			}
		}

	case WastePile:
		if len(m.waste.Cards) > 0 {
			topCard := m.waste.Cards[len(m.waste.Cards)-1]
			card = &topCard
			cards = []Card{topCard}
			cardIndex = len(m.waste.Cards) - 1
		}

	case FoundationPile:
		pile := m.foundation[location.PileIndex]
		if len(pile.Cards) > 0 {
			topCard := pile.Cards[len(pile.Cards)-1]
			card = &topCard
			cards = []Card{topCard}
			cardIndex = len(pile.Cards) - 1
		}
	}

	if card != nil {
		// Set up drag state (will be used if mouse moves)
		m.draggingCard = card
		m.draggingCards = cards
		m.dragFromPile = location
		m.dragFromIndex = cardIndex
	}

	return m, nil
}

// handleMouseRelease handles mouse button release (drop card or select)
func (m Model) handleMouseRelease(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.draggingCard == nil {
		return m, nil
	}

	// Calculate distance moved since press
	dx := msg.X - m.mousePressX
	dy := msg.Y - m.mousePressY
	distanceMoved := dx*dx + dy*dy

	// If mouse didn't move much, treat as click (select for keyboard movement)
	// Otherwise treat as drag-and-drop
	if distanceMoved < 4 { // Threshold: less than 2 pixels in any direction
		// Click to select - set keyboard selection state
		if m.dragFromPile != nil {
			m.selectedPile = m.dragFromPile
			m.selectedCard = m.draggingCard
			m.selectedIndex = m.dragFromIndex
			// Update cursor to match selection
			m.cursor = *m.dragFromPile
		}
	} else {
		// Drag-and-drop - attempt to move the cards
		location := m.getPileAtPosition(msg.X, msg.Y)
		if location != nil && m.dragFromPile != nil {
			success := m.MoveCards(*m.dragFromPile, *location, m.draggingCards)

			if success {
				// Clear selection after successful move
				m.selectedPile = nil
				m.selectedCard = nil

				// Check win condition
				if m.CheckWin() {
					m.state = StateWon
					m.elapsedTime = m.elapsedTime
					m.StartWaterfallAnimation()
					m.draggingCard = nil
					m.draggingCards = nil
					m.dragFromPile = nil
					return m, animationTick()
				}
			}
		}
	}

	// Clear drag state
	m.draggingCard = nil
	m.draggingCards = nil
	m.dragFromPile = nil

	return m, nil
}

// handleRightClick handles right-click for auto-move to foundation
func (m Model) handleRightClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	location := m.getPileAtPosition(msg.X, msg.Y)
	if location == nil {
		return m, nil
	}

	var card *Card

	switch location.PileType {
	case TableauPile:
		pile := m.tableau[location.PileIndex]
		if len(pile.Cards) > 0 && pile.Cards[len(pile.Cards)-1].FaceUp {
			card = &pile.Cards[len(pile.Cards)-1]
		}

	case WastePile:
		if len(m.waste.Cards) > 0 {
			card = &m.waste.Cards[len(m.waste.Cards)-1]
		}
	}

	if card != nil {
		m.AutoMoveToFoundation(*card, *location)

		// Check win condition
		if m.CheckWin() {
			m.state = StateWon
			m.elapsedTime = m.elapsedTime
			m.StartWaterfallAnimation()
			return m, animationTick()
		}
	}

	return m, nil
}

// getPileAtPosition determines which pile is at the given screen position
func (m *Model) getPileAtPosition(x, y int) *CursorLocation {
	// Calculate padding (must match viewGame() rendering logic)
	maxTableauHeight := 0
	for _, pile := range m.tableau {
		height := 0
		if len(pile.Cards) > 0 {
			height = (len(pile.Cards)-1)*2 + 5
		} else {
			height = 5
		}
		if height > maxTableauHeight {
			maxTableauHeight = height
		}
	}

	// Total height calculation (must match view.go exactly)
	totalLines := 11 + maxTableauHeight
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	contentWidth := 7 * 7
	leftPadding := (m.termWidth - contentWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Adjust mouse coordinates for padding
	relX := x - leftPadding
	relY := y - topPadding

	// Match rendering logic from viewGame():
	// Line 0: Title
	// Line 1: blank (from "\n\n")
	// Line 2: Stats
	// Line 3: blank (from "\n\n")
	// Lines 4-8: TopRow (cards are 5 lines: border + 3 content + border)
	// Line 9: blank (from "\n\n")
	// Line 10+: Tableau

	const topRowStartY = 4   // Top row starts at line 4 (relative to content)
	const topRowHeight = 5   // Cards are 5 lines tall
	const tableauStartY = 10 // Tableau starts at line 10 (relative to content)

	// Top row (stock, waste, foundation) - wider Y range for better detection
	if relY >= topRowStartY && relY < topRowStartY+topRowHeight {
		// Stock pile: X = 0-7
		if relX >= 0 && relX <= 7 {
			return &CursorLocation{PileType: StockPile, PileIndex: 0}
		}

		// Waste pile: X = 7-14
		if relX >= 7 && relX <= 14 {
			return &CursorLocation{PileType: WastePile, PileIndex: 0}
		}

		// Foundation piles: X starts around 17 (with spacing)
		// Each pile is 7 chars wide
		if relX >= 17 {
			foundationIndex := (relX - 17) / 7
			if foundationIndex >= 0 && foundationIndex < 4 {
				return &CursorLocation{PileType: FoundationPile, PileIndex: foundationIndex}
			}
		}
	}

	// Tableau piles start at tableauStartY
	if relY >= tableauStartY {
		// Each tableau pile is 7 chars wide
		tableauIndex := relX / 7
		if tableauIndex >= 0 && tableauIndex < 7 {
			return &CursorLocation{PileType: TableauPile, PileIndex: tableauIndex}
		}
	}

	return nil
}

// getTableauCardIndex determines which card in a tableau pile was clicked
func (m *Model) getTableauCardIndex(pileIndex, y int) int {
	// With stacking: each card except the last shows 2 lines (top border + first content line)
	// The last card shows full (5 lines total)

	// Calculate padding to adjust Y coordinate
	maxTableauHeight := 0
	for _, pile := range m.tableau {
		height := 0
		if len(pile.Cards) > 0 {
			height = (len(pile.Cards)-1)*2 + 5
		} else {
			height = 5
		}
		if height > maxTableauHeight {
			maxTableauHeight = height
		}
	}

	// Total height calculation (must match view.go exactly)
	totalLines := 11 + maxTableauHeight
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	const tableauStartY = 10 // Must match getPileAtPosition (relative to content)
	pile := m.tableau[pileIndex]

	if len(pile.Cards) == 0 {
		return -1
	}

	// Calculate which card was clicked based on Y position (adjusted for padding)
	relY := y - topPadding
	cardOffset := relY - tableauStartY
	if cardOffset < 0 {
		return -1
	}

	// Each stacked card takes 2 lines, last card takes 5 lines (Height(3) + 2 borders)
	numCards := len(pile.Cards)
	if numCards == 1 {
		// Only one card, click anywhere on it
		return 0
	}

	// Calculate total height of stacked area
	stackedHeight := (numCards - 1) * 2

	// If clicked in the stacked area
	if cardOffset < stackedHeight {
		return cardOffset / 2 // Each stacked card is 2 lines
	}

	// Otherwise clicked on the last (full) card
	return numCards - 1
}
