package snake

import tea "github.com/charmbracelet/bubbletea"

// handleKeyPress processes keyboard input based on current game state
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys (work in all states)
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "s", "S":
		// Toggle stats view (from menu or game over state)
		if m.state == StateMenu || m.state == StateGameOver {
			m.previousState = m.state
			m.state = StateStats
			return m, nil
		} else if m.state == StateStats {
			// Return from stats view
			m.state = m.previousState
			return m, nil
		}

	case "esc":
		// ESC returns from stats view
		if m.state == StateStats {
			m.state = m.previousState
			return m, nil
		}
	}

	switch m.state {
	case StateMenu:
		switch msg.String() {
		case "enter":
			m.state = StateDifficultySelect
		}

	case StateDifficultySelect:
		switch msg.String() {
		case "up", "k", "w":
			if m.selectedDifficulty > Easy {
				m.selectedDifficulty--
			}
		case "down", "j", "s":
			if m.selectedDifficulty < Hard {
				m.selectedDifficulty++
			}
		case "enter":
			m.startGame()
			return m, tickCmd(m.speed)
		case "esc":
			m.state = StateMenu
		}

	case StatePlaying:
		switch msg.String() {
		case "up", "k", "w":
			m.nextDir = Up
		case "down", "j", "s":
			m.nextDir = Down
		case "left", "h", "a":
			m.nextDir = Left
		case "right", "l", "d":
			m.nextDir = Right
		case "p", " ":
			m.state = StatePaused
		}

	case StatePaused:
		switch msg.String() {
		case "p", " ":
			m.state = StatePlaying
		}

	case StateGameOver:
		switch msg.String() {
		case "r":
			m.startGame()
			return m, tickCmd(m.speed)
		case "m":
			m.state = StateMenu
		}
	}

	return m, nil
}
