package solitaire

import (
	"fmt"

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

// handleMousePress handles mouse button press (start drag)
func (m Model) handleMousePress(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// DEBUG: Print mouse coordinates to help debug hit detection
	// TODO: Remove this debug output once mouse detection is fixed
	fmt.Printf("\nDEBUG: Mouse click at X=%d, Y=%d (Terminal: %dx%d)\n",
		msg.X, msg.Y, m.termWidth, m.termHeight)

	// Determine which pile was clicked
	location := m.getPileAtPosition(msg.X, msg.Y)
	if location == nil {
		fmt.Printf("DEBUG: No pile detected at this position\n")
		return m, nil
	}

	fmt.Printf("DEBUG: Detected pile type=%v, index=%d\n", location.PileType, location.PileIndex)

	// If stock was clicked, draw a card
	if location.PileType == StockPile {
		m.DrawFromStock()
		return m, nil
	}

	// Start dragging a card
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
		m.draggingCard = card
		m.draggingCards = cards
		m.dragFromPile = location
		m.dragFromIndex = cardIndex
	}

	return m, nil
}

// handleMouseRelease handles mouse button release (drop card)
func (m Model) handleMouseRelease(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.draggingCard == nil {
		return m, nil
	}

	// Determine where the card was dropped
	location := m.getPileAtPosition(msg.X, msg.Y)
	if location != nil && m.dragFromPile != nil {
		// Attempt to move the cards
		success := m.MoveCards(*m.dragFromPile, *location, m.draggingCards)

		if success {
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
	// Layout: Title(1) + margin(1) + \n\n(2) + Stats(1) + margin(1) + \n\n(2) = 8 lines
	// Then: TopRow(5) + \n\n(2) = 7 more lines
	// Total before tableau: 15 lines

	// Make hit detection forgiving - allow clicks a bit before/after expected positions
	const topRowStartY = 4   // Start checking earlier
	const topRowHeight = 10  // Wider range to catch clicks
	const tableauStartY = 11 // Start checking earlier for tableau too

	// Top row (stock, waste, foundation) - wider Y range for better detection
	if y >= topRowStartY && y < topRowStartY+topRowHeight {
		// Stock pile: X = 0-7
		if x >= 0 && x <= 7 {
			return &CursorLocation{PileType: StockPile, PileIndex: 0}
		}

		// Waste pile: X = 7-14
		if x >= 7 && x <= 14 {
			return &CursorLocation{PileType: WastePile, PileIndex: 0}
		}

		// Foundation piles: X starts around 17 (with spacing)
		// Each pile is 7 chars wide
		if x >= 17 {
			foundationIndex := (x - 17) / 7
			if foundationIndex >= 0 && foundationIndex < 4 {
				return &CursorLocation{PileType: FoundationPile, PileIndex: foundationIndex}
			}
		}
	}

	// Tableau piles start at tableauStartY
	if y >= tableauStartY {
		// Each tableau pile is 7 chars wide
		tableauIndex := x / 7
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

	const tableauStartY = 15
	pile := m.tableau[pileIndex]

	if len(pile.Cards) == 0 {
		return -1
	}

	// Calculate which card was clicked based on Y position
	cardOffset := y - tableauStartY
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
