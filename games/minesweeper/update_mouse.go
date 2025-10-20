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
	// Calculate the grid's screen position
	// Each cell is 2 characters wide (to make it more square-looking)
	// The grid is centered on screen

	cellWidth := 2
	cellHeight := 1

	// Calculate grid offset to center it
	gridPixelWidth := m.width * cellWidth
	gridPixelHeight := m.height * cellHeight

	offsetX := (m.termWidth - gridPixelWidth) / 2
	offsetY := (m.termHeight - gridPixelHeight) / 2

	// Account for header (title + stats line = ~6 lines)
	offsetY += 6

	// Convert mouse position to grid position
	relX := mouseX - offsetX
	relY := mouseY - offsetY

	gridX = relX / cellWidth
	gridY = relY / cellHeight

	// Check bounds
	inBounds = gridX >= 0 && gridX < m.width && gridY >= 0 && gridY < m.height

	return gridX, gridY, inBounds
}
