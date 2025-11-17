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

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		// Resize landing page
		if m.landingPage != nil {
			m.landingPage.Resize(msg.Width, msg.Height)
		}
		return m, nil

	case AnimationTickMsg:
		// Update landing page animation
		if m.landingPage != nil {
			m.landingPage.Update()
		}
		return m, animationTick() // Continue animation loop
	}

	return m, nil
}

// handleGameUpdate delegates messages to the current game
// and handles returning to menu if the game quits
func (m Model) handleGameUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle animation tick to keep landing page animation alive
	// This ensures the animation doesn't freeze when we return to menu
	if _, ok := msg.(AnimationTickMsg); ok {
		if m.landingPage != nil {
			m.landingPage.Update()
		}
		// Continue animation loop AND delegate to game
		updatedGame, gameCmd := m.currentGame.Update(msg)
		m.currentGame = updatedGame
		return m, tea.Batch(animationTick(), gameCmd)
	}

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
