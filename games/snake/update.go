package snake

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent on each game tick to update the snake's position
type TickMsg time.Time

// tickCmd creates a command that sends a tick message after duration d
func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init initializes the game and requests terminal size
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		tickCmd(m.speed),
	)
}

// Update handles all incoming messages and updates the game state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case TickMsg:
		if m.state == StatePlaying {
			m.moveSnake()

			if m.checkCollision() {
				m.gameOver()
			}

			// Speed up as score increases (max speed: 30ms)
			newSpeed := 100*time.Millisecond - time.Duration(m.score)*2*time.Millisecond
			if newSpeed < 30*time.Millisecond {
				newSpeed = 30 * time.Millisecond
			}
			m.speed = newSpeed
		}

		return m, tickCmd(m.speed)
	}

	return m, nil
}
