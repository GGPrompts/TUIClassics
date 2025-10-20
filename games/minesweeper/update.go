package minesweeper

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case TickMsg:
		if m.state == StatePlaying {
			m.elapsedTime = time.Since(m.startTime)
			return m, tick()
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "n":
		// New game
		m.InitGame()
		return m, tick()

	case "1":
		if m.state == StateMenu {
			m.difficulty = DifficultyEasy
			m.InitGame()
			return m, tick()
		}

	case "2":
		if m.state == StateMenu {
			m.difficulty = DifficultyMedium
			m.InitGame()
			return m, tick()
		}

	case "3":
		if m.state == StateMenu {
			m.difficulty = DifficultyHard
			m.InitGame()
			return m, tick()
		}

	case "r":
		// Restart current difficulty
		if m.state == StateWon || m.state == StateLost {
			m.InitGame()
			return m, tick()
		}

	// Keyboard navigation (for accessibility)
	case "up", "k":
		if m.state == StatePlaying && m.cursorY > 0 {
			m.cursorY--
		}

	case "down", "j":
		if m.state == StatePlaying && m.cursorY < m.height-1 {
			m.cursorY++
		}

	case "left", "h":
		if m.state == StatePlaying && m.cursorX > 0 {
			m.cursorX--
		}

	case "right", "l":
		if m.state == StatePlaying && m.cursorX < m.width-1 {
			m.cursorX++
		}

	case "enter", " ":
		// Reveal cell at cursor
		if m.state == StatePlaying {
			m.RevealCell(m.cursorX, m.cursorY)
		}

	case "f":
		// Flag cell at cursor
		if m.state == StatePlaying {
			m.ToggleFlag(m.cursorX, m.cursorY)
		}
	}

	return m, nil
}
