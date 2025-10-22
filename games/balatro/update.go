package balatro

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// update.go - Main Update Dispatcher
// Purpose: Message dispatching and non-input event handling
// When to extend: Add new message types or top-level event handlers here

// Init is called when the program starts
func (m Model) Init() tea.Cmd {
	// Request window size and start animation
	return tea.Batch(
		tea.WindowSize(),
		tickCmd(),
	)
}

// tickCmd creates a command that sends tick messages for animation
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles all messages and updates the Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Mouse input
	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	// Custom messages
	case errMsg:
		m.err = msg.err
		m.statusMsg = "Error: " + msg.err.Error()
		return m, nil

	case statusMsg:
		m.statusMsg = msg.message
		return m, nil

	// Animation tick for landing page
	case tickMsg:
		if m.state == StateLanding && m.landingPage != nil {
			m.landingPage.Update()
		}
		return m, tickCmd() // Continue animation

	// Add handlers for your custom messages here
	// Example:
	// case itemSelectedMsg:
	//     return m.handleItemSelected(msg)
	//
	// case dataLoadedMsg:
	//     return m.handleDataLoaded(msg)
	}

	return m, nil
}

// Helper functions for message handling

// sendStatus creates a status message command
func sendStatus(message string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{message: message}
	}
}

// sendError creates an error message command
func sendError(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err: err}
	}
}

// isSpecialKey checks if a key is a special key (not printable)
func isSpecialKey(key tea.KeyMsg) bool {
	return key.Type != tea.KeyRunes
}
