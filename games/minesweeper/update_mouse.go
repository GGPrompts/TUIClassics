package minesweeper

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseEvent processes mouse events
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Check if smiley was clicked (works in any state except menu)
	if m.state != StateMenu {
		if m.isSmileyClicked(msg.X, msg.Y) && msg.Action == tea.MouseActionPress {
			m.InitGame()
			return m, tick()
		}
	}

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
			// Show surprised face when clicking
			m.smileyState = SmileySurprised
			m.RevealCell(gridX, gridY)
			// If explosion animation started, begin animation ticks
			if m.state == StateExploding {
				return m, animationTick()
			}
		} else if msg.Action == tea.MouseActionRelease {
			// Return to happy if still playing
			if m.state == StatePlaying {
				m.smileyState = SmileyHappy
			}
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
	totalLines := 3 + 2 + 1 + 2 + m.height + 4 + 2 // title + gap + smiley + gap + grid + gap + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Grid starts after: topPadding + title line + "\n\n" (2) + smiley line + "\n\n" (2) + stats line + "\n\n" (2)
	// Line topPadding: title
	// Line topPadding+2: smiley (after \n\n from title)
	// Line topPadding+4: stats (after \n\n from smiley)
	// Line topPadding+6: grid starts (after \n\n from stats)
	gridStartY := topPadding + 6

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

// isSmileyClicked checks if the mouse click was on the smiley button
func (m Model) isSmileyClicked(mouseX, mouseY int) bool {
	// Smiley is rendered at:
	// topPadding + title(1) + \n\n(2) = topPadding + 3

	totalLines := 3 + 2 + m.height + 4 + 2
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	smileyY := topPadding + 3

	// Smiley is "[ :) ]" = 6 chars wide, centered
	smileyWidth := 6
	smileyX := (m.termWidth - smileyWidth) / 2

	// Check if click is within smiley bounds
	return mouseY == smileyY && mouseX >= smileyX && mouseX < smileyX+smileyWidth
}
