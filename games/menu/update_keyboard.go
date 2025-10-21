package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		// Quit the application
		return m, tea.Quit

	case "up", "k":
		// Move selection up
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
		return m, nil

	case "down", "j":
		// Move selection down
		if m.selectedIdx < len(m.games)-1 {
			m.selectedIdx++
		}
		return m, nil

	case "enter", " ":
		// Launch selected game
		cmd := m.launchGame(m.selectedIdx)
		return m, cmd

	default:
		// Check for hotkey matches
		key := msg.String()
		for i, game := range m.games {
			if key == game.Hotkey {
				m.selectedIdx = i
				cmd := m.launchGame(i)
				return m, cmd
			}
		}
	}

	return m, nil
}
