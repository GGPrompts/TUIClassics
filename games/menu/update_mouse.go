package menu

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const doubleClickThreshold = 500 * time.Millisecond

// handleMouseEvent handles mouse clicks and scrolling
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Handle high scores view separately
	if m.state == StateHighScores {
		return m.handleHighScoresMouseEvent(msg)
	}

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

					// Check if it's the High Scores button
					if selection == "High Scores 🏆" {
						m.state = StateHighScores
						m.currentTab = 0
						m.lastClickTime = time.Time{}
						m.lastClickButton = -1
						return m, nil
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

				// Check if it's the High Scores button
				if selection == "High Scores 🏆" {
					m.state = StateHighScores
					m.currentTab = 0
					return m, nil
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

// handleHighScoresMouseEvent handles mouse clicks in high scores view
func (m Model) handleHighScoresMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Type != tea.MouseLeft || msg.Action != tea.MouseActionPress {
		return m, nil
	}

	// Calculate tab positions (tabs are centered and appear after title)
	// Title takes ~3 lines, tabs appear around line 5-6
	// Each tab is roughly 16 chars wide with 1 space between

	// Approximate tab bar Y position (adjust based on your layout)
	tabBarY := (m.termHeight-30)/2 + 3

	// Check if click is on tab bar row
	if msg.Y >= tabBarY && msg.Y <= tabBarY+1 {
		// Calculate tab positions (tabs are centered)
		// Tab widths: "  Minesweeper  " = ~16, "  2048  " = ~8, etc.
		tabs := []string{"Minesweeper", "2048", "Solitaire", "Snake"}
		totalWidth := 0
		for i, tab := range tabs {
			// Each tab has padding of 2 on each side (4 total)
			tabWidth := len(tab) + 4
			if i < len(tabs)-1 {
				tabWidth += 1 // Space between tabs
			}
			totalWidth += tabWidth
		}

		// Center position
		startX := (m.termWidth - totalWidth) / 2
		currentX := startX

		// Check which tab was clicked
		for i, tab := range tabs {
			tabWidth := len(tab) + 4
			if msg.X >= currentX && msg.X < currentX+tabWidth {
				m.currentTab = i
				return m, nil
			}
			currentX += tabWidth + 1 // +1 for space
		}
	}

	return m, nil
}
