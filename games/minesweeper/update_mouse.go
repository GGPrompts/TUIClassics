package minesweeper

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseEvent processes mouse events
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.state != StatePlaying {
		// Handle menu clicks when in menu state
		if m.state == StateMenu {
			return m.handleMenuClick(msg)
		}
		return m, nil
	}

	// Convert mouse position to grid coordinates
	gridX, gridY, inBounds := m.mouseToGrid(msg.X, msg.Y)
	if !inBounds {
		return m, nil
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			m.RevealCell(gridX, gridY)
		}

	case tea.MouseButtonRight:
		if msg.Action == tea.MouseActionPress {
			m.ToggleFlag(gridX, gridY)
		}
	}

	return m, nil
}

// handleMenuClick handles clicks on the difficulty menu
func (m Model) handleMenuClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress {
		return m, nil
	}

	// Calculate menu button positions (we'll refine this in view.go)
	// For now, accept clicks anywhere on the screen to start easy mode
	// TODO: Implement proper button hit detection based on rendered positions

	return m, nil
}

// mouseToGrid converts terminal coordinates to grid coordinates
func (m Model) mouseToGrid(mouseX, mouseY int) (gridX, gridY int, inBounds bool) {
	// Match the rendering logic from renderGame()

	// Calculate vertical position
	totalLines := 3 + 2 + m.height + 4 + 2 // title + gap + grid + gap + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Grid starts after: topPadding + title(1) + gap(2) + stats(1) + gap(2)
	gridStartY := topPadding + 1 + 2 + 1 + 2

	// Calculate horizontal position
	gridWidth := m.width * 2 // Each cell is roughly 2 chars
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Convert mouse position to grid position
	relX := mouseX - leftPadding
	relY := mouseY - gridStartY

	// Each cell is 2 characters wide (symbol + space)
	gridX = relX / 2
	gridY = relY

	// Check bounds
	inBounds = gridX >= 0 && gridX < m.width && gridY >= 0 && gridY < m.height

	return gridX, gridY, inBounds
}
