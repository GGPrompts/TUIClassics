package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// ESC handling - return from high scores or game
	if msg.String() == "esc" {
		if m.state == StateHighScores {
			m.state = StateMainMenu
			return m, nil
		}
	}

	// Quit keys
	if msg.String() == "q" || msg.String() == "ctrl+c" || msg.String() == "f10" {
		return m, tea.Quit
	}

	// State-specific handling
	if m.state == StateHighScores {
		return m.handleHighScoresKeys(msg)
	}

	// Main menu keys
	switch msg.String() {
	case "up", "k":
		// Move selection up
		if m.landingPage != nil {
			m.landingPage.SelectPrev()
			m.selectedIdx = m.landingPage.selectedBtn
		} else if m.selectedIdx > 0 {
			m.selectedIdx--
		}
		return m, nil

	case "down", "j":
		// Move selection down
		if m.landingPage != nil {
			m.landingPage.SelectNext()
			m.selectedIdx = m.landingPage.selectedBtn
		} else if m.selectedIdx < len(m.games)-1 {
			m.selectedIdx++
		}
		return m, nil

	case "b", "B":
		// Toggle background (starfield/checkered)
		if m.landingPage != nil {
			m.landingPage.ToggleBackground()
		}
		return m, nil

	case "enter", " ":
		// Launch selected game, show high scores, or handle Exit
		if m.landingPage != nil {
			selection := m.landingPage.GetSelectedItem()
			// Check if Exit button selected
			if selection == "Exit 🚪" {
				return m, tea.Quit
			}
			// Check if High Scores selected
			if selection == "High Scores 🏆" {
				m.state = StateHighScores
				m.currentTab = 0 // Start with Minesweeper tab
				return m, nil
			}
			// Otherwise, launch the game at selectedBtn index
			cmd := m.launchGame(m.landingPage.selectedBtn)
			return m, cmd
		} else {
			cmd := m.launchGame(m.selectedIdx)
			return m, cmd
		}
		return m, nil

	default:
		// Check for hotkey matches
		key := msg.String()
		for i, game := range m.games {
			if key == game.Hotkey {
				m.selectedIdx = i
				if m.landingPage != nil {
					m.landingPage.selectedBtn = i
				}
				cmd := m.launchGame(i)
				return m, cmd
			}
		}
	}

	return m, nil
}

// handleHighScoresKeys handles keyboard input in high scores view
func (m Model) handleHighScoresKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "right", "l":
		// Next tab
		m.currentTab = (m.currentTab + 1) % 4 // 4 tabs total
		return m, nil

	case "shift+tab", "left", "h":
		// Previous tab
		m.currentTab--
		if m.currentTab < 0 {
			m.currentTab = 3
		}
		return m, nil

	case "1":
		m.currentTab = 0 // Minesweeper
		return m, nil

	case "2":
		m.currentTab = 1 // 2048
		return m, nil

	case "3":
		m.currentTab = 2 // Solitaire
		return m, nil

	case "4":
		m.currentTab = 3 // Snake
		return m, nil
	}

	return m, nil
}
