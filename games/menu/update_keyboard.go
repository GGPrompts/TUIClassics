package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "f10":
		// Quit the application
		return m, tea.Quit

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
		// Launch selected game or handle Exit
		if m.landingPage != nil {
			selection := m.landingPage.GetSelectedItem()
			// Check if Exit button selected
			if selection == "Exit 🚪" {
				return m, tea.Quit
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
