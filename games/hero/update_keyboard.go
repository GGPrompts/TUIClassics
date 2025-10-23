package hero

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMenuInput handles keyboard input in the menu state
func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		m.songIndex--
		if m.songIndex < 0 {
			m.songIndex = len(getDemoSongs()) - 1
		}

	case "down", "j":
		m.songIndex++
		if m.songIndex >= len(getDemoSongs()) {
			m.songIndex = 0
		}

	case "enter":
		// Start the selected song
		songs := getDemoSongs()
		if m.songIndex >= 0 && m.songIndex < len(songs) {
			m.startSong(&songs[m.songIndex])
		}

	// Number hotkeys for direct selection
	case "1":
		songs := getDemoSongs()
		if len(songs) >= 1 {
			m.startSong(&songs[0])
		}
	case "2":
		songs := getDemoSongs()
		if len(songs) >= 2 {
			m.startSong(&songs[1])
		}
	case "3":
		songs := getDemoSongs()
		if len(songs) >= 3 {
			m.startSong(&songs[2])
		}
	}

	return m, nil
}

// handleGameInput handles keyboard input during gameplay
func (m Model) handleGameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		// Return to menu
		m.state = StateMenu
		return m, nil

	// Lane keys - check for hits
	case "a", "A":
		m.checkHit(0)
	case "s", "S":
		m.checkHit(1)
	case "d", "D":
		m.checkHit(2)
	case "f", "F":
		m.checkHit(3)
	case " ", "j", "J": // Spacebar (primary) or J (alternative) for lane 5
		m.checkHit(4)
	}

	return m, nil
}

// handleFinishedInput handles keyboard input on the results screen
func (m Model) handleFinishedInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "esc":
		// Return to menu
		m.state = StateMenu
		return m, nil

	case "enter":
		// Play again - restart the same song
		if m.currentSong != nil {
			m.startSong(m.currentSong)
		}
	}

	return m, nil
}
