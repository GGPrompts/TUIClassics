package menu

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const doubleClickThreshold = 500 * time.Millisecond

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
			// Easter egg: Check if X button was clicked
			if m.landingPage.HandleCloseButtonClick(msg.X, msg.Y) {
				return m, tea.Quit
			}

			buttonIdx := m.landingPage.HandleClick(msg.X, msg.Y)
			if buttonIdx >= 0 {
				m.landingPage.selectedBtn = buttonIdx
				m.selectedIdx = buttonIdx

				// Double-click detection
				now := time.Now()
				timeSinceLastClick := now.Sub(m.lastClickTime)

				// If clicked same button within threshold, it's a double-click
				if buttonIdx == m.lastClickButton && timeSinceLastClick < doubleClickThreshold {
					// Double-click detected! Launch immediately
					selection := m.landingPage.GetSelectedItem()

					// Check if it's the Exit button
					if selection == "Exit 🚪" {
						return m, tea.Quit
					}

					// For all games, use the button index to launch
					if buttonIdx >= 0 && buttonIdx < len(m.games) {
						cmd := m.launchGame(buttonIdx)
						// Reset click tracking after launch
						m.lastClickTime = time.Time{}
						m.lastClickButton = -1
						return m, cmd
					}
				}

				// Record this click for double-click detection
				m.lastClickTime = now
				m.lastClickButton = buttonIdx
			}
		} else if msg.Action == tea.MouseActionRelease {
			// On release, launch the game if we're still on the same button
			// (This preserves the original single-click behavior)
			buttonIdx := m.landingPage.HandleClick(msg.X, msg.Y)
			if buttonIdx >= 0 && buttonIdx == m.landingPage.selectedBtn {
				// Launch the selected item
				selection := m.landingPage.GetSelectedItem()

				// Check if it's the Exit button
				if selection == "Exit 🚪" {
					return m, tea.Quit
				}

				// For all games, use the button index to launch
				if buttonIdx < len(m.games) {
					cmd := m.launchGame(buttonIdx)
					return m, cmd
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
