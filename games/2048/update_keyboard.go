package game2048

import tea "github.com/charmbracelet/bubbletea"

// handleKeyPress processes keyboard input based on game state
func handleKeyPress(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	// Global instructions hotkey (works from any state except already in instructions)
	if (msg.String() == "i" || msg.String() == "?") && m.state != StateInstructions {
		m.previousState = m.state
		m.state = StateInstructions
		return m, nil
	}

	switch m.state {
	case StateMenu:
		if msg.String() == "enter" || msg.String() == " " {
			m.startGame()
		}

	case StatePlaying:
		switch msg.String() {
		case "up", "k", "w":
			m.move(Up)
		case "down", "j", "s":
			m.move(Down)
		case "left", "h", "a":
			m.move(Left)
		case "right", "l", "d":
			m.move(Right)
		case "r":
			m.startGame() // Restart
		}

	case StateWon:
		switch msg.String() {
		case "c":
			// Continue playing after reaching 2048
			m.state = StatePlaying
		case "r":
			m.startGame()
		}

	case StateGameOver:
		if msg.String() == "r" {
			m.startGame()
		}

	case StateInstructions:
		// Any key returns from instructions
		if msg.String() == "escape" || msg.String() == "i" || msg.String() == "?" || msg.String() == "enter" || msg.String() == " " {
			m.state = m.previousState
		}
	}

	return m, nil
}
