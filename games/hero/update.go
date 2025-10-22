package hero

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the game
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.WindowSize(), // Request terminal dimensions
		tickCmd(),
	)
}

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case TickMsg:
		if m.state == StatePlaying {
			m.currentTime = time.Now()
			m.updateNotePositions()
			m.checkMissedNotes()
			m.checkSongEnd()

			// Hide hit feedback after 500ms
			if m.showHitFeedback && time.Since(m.hitFeedbackTime) > 500*time.Millisecond {
				m.showHitFeedback = false
			}
		}
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress delegates to update_keyboard.go
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateMenu:
		return m.handleMenuInput(msg)
	case StatePlaying:
		return m.handleGameInput(msg)
	case StateFinished:
		return m.handleFinishedInput(msg)
	default:
		return m, nil
	}
}
