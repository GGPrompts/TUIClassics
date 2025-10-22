package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseEvent handles mouse clicks and scrolling
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle mouse in main menu (not in-game)
	if m.state != StateMainMenu || m.landingPage == nil {
		return m, nil
	}

	switch msg.Type {
	case tea.MouseLeft:
		// Left click - check if button was clicked
		if msg.Action == tea.MouseActionPress {
			buttonIdx := m.landingPage.HandleClick(msg.X, msg.Y)
			if buttonIdx >= 0 {
				m.landingPage.selectedBtn = buttonIdx
				m.selectedIdx = buttonIdx
			}
		} else if msg.Action == tea.MouseActionRelease {
			// On release, launch the game if we're still on the same button
			buttonIdx := m.landingPage.HandleClick(msg.X, msg.Y)
			if buttonIdx >= 0 && buttonIdx == m.landingPage.selectedBtn {
				// Launch the selected item
				selection := m.landingPage.GetSelectedItem()
				switch selection {
				case "Minesweeper 💣":
					cmd := m.launchGame(0)
					return m, cmd
				case "Solitaire 🂡":
					cmd := m.launchGame(1)
					return m, cmd
				case "Exit 🚪":
					return m, tea.Quit
				}
			}
		}

	case tea.MouseWheelUp:
		// Scroll up - previous item
		m.landingPage.SelectPrev()
		m.selectedIdx = m.landingPage.selectedBtn

	case tea.MouseWheelDown:
		// Scroll down - next item
		m.landingPage.SelectNext()
		m.selectedIdx = m.landingPage.selectedBtn
	}

	return m, nil
}
