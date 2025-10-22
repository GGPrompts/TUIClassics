package balatro

import (
	tea "github.com/charmbracelet/bubbletea"
)

// update_mouse.go - Mouse Event Handling
// Purpose: All mouse input processing
// When to extend: Add new mouse interactions or clickable elements here

// handleMouseEvent handles mouse input
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.config.UI.MouseEnabled {
		return m, nil
	}

	// Dispatch based on game phase
	switch m.gamePhase {
	case PhaseSelectCards:
		return m.handleGamePhaseMouseEvent(msg)
	case PhaseShop:
		return m.handleShopPhaseMouseEvent(msg)
	default:
		return m, nil
	}
}

// handleGamePhaseMouseEvent handles mouse events during card selection phase
func (m Model) handleGamePhaseMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			m.mousePressX = msg.X
			m.mousePressY = msg.Y
		} else if msg.Action == tea.MouseActionRelease {
			// Check if this was a click (not a drag)
			dx := msg.X - m.mousePressX
			dy := msg.Y - m.mousePressY
			distanceMoved := dx*dx + dy*dy

			if distanceMoved < 4 { // Threshold: less than 2 pixels
				// Handle click on hand cards
				cardIndex := m.getCardAtPosition(msg.X, msg.Y)
				if cardIndex >= 0 && cardIndex < len(m.hand) {
					// Toggle card for play (same as space key)
					m.selectedCardIndex = cardIndex
					m.toggleCardForPlay(cardIndex)
					return m, nil
				}
			}
		}
	}

	return m, nil
}

// handleShopPhaseMouseEvent handles mouse events during shop phase
func (m Model) handleShopPhaseMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			m.mousePressX = msg.X
			m.mousePressY = msg.Y
		} else if msg.Action == tea.MouseActionRelease {
			// Check if this was a click (not a drag)
			dx := msg.X - m.mousePressX
			dy := msg.Y - m.mousePressY
			distanceMoved := dx*dx + dy*dy

			if distanceMoved < 4 { // Threshold: less than 2 pixels
				// Check if clicked on a shop joker
				jokerIndex := m.getShopJokerAtPosition(msg.X, msg.Y)
				if jokerIndex >= 0 && jokerIndex < len(m.shopJokers) {
					m.selectedShopItem = jokerIndex
					return m, nil
				}

				// Check if clicked on "Buy" action (Enter key equivalent)
				// For now, we'll trigger buy when clicking on a selected joker again
				if jokerIndex == m.selectedShopItem && jokerIndex >= 0 {
					return m.handleShopPurchase()
				}
			}
		}
	}

	return m, nil
}

// getCardAtPosition calculates which card in the hand was clicked
func (m Model) getCardAtPosition(x, y int) int {
	if len(m.hand) == 0 {
		return -1
	}

	// Calculate content dimensions and padding
	// This must match the rendering logic in renderGameView()

	// Build sections to calculate content height (must match view.go)
	lineCount := 0

	// Game info (1 line) + blank
	lineCount += 2

	// Jokers (if any): label (1) + cards (varies) + blank
	if len(m.jokers) > 0 {
		lineCount += 1 + 5 + 1 // label + joker height + blank
	}

	// Hand display: cards (3) + numbers (1) + blank
	lineCount += 5 // Includes centering wrapper

	// Selected card detail (if visible): varies
	if m.selectedCardIndex >= 0 && m.selectedCardIndex < len(m.hand) {
		lineCount += 7 + 1 // Card detail panel + blank
	}

	// Current hand info (always present now): 1 line + blank
	lineCount += 2

	// Score breakdown (if visible)
	if m.lastScore.FinalScore > 0 {
		lineCount += 5 + 1 // Score panel + blank
	}

	// Controls: 1 line
	lineCount += 1

	// Calculate padding
	availableHeight := m.height - 2 // -2 for title and status bars
	topPadding := 0
	if lineCount < availableHeight {
		topPadding = (availableHeight - lineCount) / 2
	}

	// Add title bar offset
	if m.config.UI.ShowTitle {
		topPadding += 1
	}

	// Calculate Y position of hand display
	handY := topPadding
	handY += 2 // game info + blank
	if len(m.jokers) > 0 {
		handY += 1 + 5 + 1 // jokers section
	}

	// Check if click is in hand area (cards are 3 lines tall + 1 line for numbers)
	if y < handY || y >= handY+4 {
		return -1
	}

	// Calculate horizontal position
	// Each card is 7 chars wide (5 for card + 2 spacing from join)
	cardWidth := 7
	totalHandWidth := len(m.hand) * cardWidth

	// Hand is centered horizontally
	leftPadding := (m.width - totalHandWidth) / 2

	// Adjust for padding
	relX := x - leftPadding

	if relX < 0 || relX >= totalHandWidth {
		return -1
	}

	// Calculate which card was clicked
	cardIndex := relX / cardWidth

	if cardIndex >= 0 && cardIndex < len(m.hand) {
		return cardIndex
	}

	return -1
}

// getShopJokerAtPosition calculates which shop joker was clicked
func (m Model) getShopJokerAtPosition(x, y int) int {
	if len(m.shopJokers) == 0 {
		return -1
	}

	// Shop screen is centered with lipgloss.Place
	// We need to calculate the position of shop jokers

	// Shop jokers are rendered starting around line 900 in view.go
	// Each joker is 22 chars wide + 2 margin right = 24 chars total
	// Jokers are 6 lines tall

	// Calculate vertical position
	// Shop layout: title(1) + blank(1) + money(1) + blank(1) + next blind(1) + blank(1) + label(1) + blank(1) + jokers(6)
	// = 13 lines to start of jokers

	// The shop is centered on screen
	// We need to calculate content height
	contentHeight := 13 + 6 // Simplified - actual height varies
	topPadding := (m.height - contentHeight) / 2

	jokersStartY := topPadding + 8 // After title, money, next blind, label

	// Check if click is in joker area
	if y < jokersStartY || y >= jokersStartY+6 {
		return -1
	}

	// Calculate horizontal position
	// Each joker is 24 chars wide (22 + 2 margin)
	jokerWidth := 24
	totalWidth := len(m.shopJokers) * jokerWidth

	leftPadding := (m.width - totalWidth) / 2
	relX := x - leftPadding

	if relX < 0 || relX >= totalWidth {
		return -1
	}

	jokerIndex := relX / jokerWidth

	if jokerIndex >= 0 && jokerIndex < len(m.shopJokers) {
		return jokerIndex
	}

	return -1
}

// handleShopPurchase attempts to buy the selected shop joker
func (m Model) handleShopPurchase() (tea.Model, tea.Cmd) {
	if m.selectedShopItem < 0 || m.selectedShopItem >= len(m.shopJokers) {
		return m, nil
	}

	joker := m.shopJokers[m.selectedShopItem]
	cost := joker.GetCost()

	// Check if player can afford it
	if m.roundState.Money < cost {
		m.statusMsg = "Not enough money!"
		return m, nil
	}

	// Check if player has room for more jokers (max 5)
	if len(m.jokers) >= 5 {
		m.statusMsg = "Maximum jokers reached (5/5)"
		return m, nil
	}

	// Purchase the joker
	m.roundState.Money -= cost
	m.jokers = append(m.jokers, joker)

	// Remove from shop
	m.shopJokers = append(m.shopJokers[:m.selectedShopItem], m.shopJokers[m.selectedShopItem+1:]...)

	// Reset selection
	m.selectedShopItem = -1
	if len(m.shopJokers) > 0 {
		m.selectedShopItem = 0
	}

	m.statusMsg = "Purchased joker!"

	return m, nil
}
