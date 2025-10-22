package game2048

import tea "github.com/charmbracelet/bubbletea"

// updateGame is the main update function for the game
func updateGame(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit key
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle game-specific keys
		return handleKeyPress(m, msg)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil
	}

	return m, nil
}
