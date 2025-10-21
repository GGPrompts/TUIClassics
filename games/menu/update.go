package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If in game, delegate to the game's update
	if m.state == StateInGame && m.currentGame != nil {
		return m.handleGameUpdate(msg)
	}

	// Handle menu updates
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil
	}

	return m, nil
}

// handleGameUpdate delegates messages to the current game
// and handles returning to menu if the game quits
func (m Model) handleGameUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for Esc key to return to menu
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			m.returnToMenu()
			return m, nil
		}
	}

	// Delegate to game
	updatedGame, cmd := m.currentGame.Update(msg)
	m.currentGame = updatedGame

	// Check if game wants to quit (would normally exit the program)
	// We intercept this and return to menu instead
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "q" || keyMsg.String() == "ctrl+c" {
			m.returnToMenu()
			return m, nil
		}
	}

	return m, cmd
}
