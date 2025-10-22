package snake

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent on each game tick to update the snake's position
type TickMsg time.Time

// CrashDelayMsg is sent after crash animation delay
type CrashDelayMsg time.Time

// tickCmd creates a command that sends a tick message after duration d
func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// crashDelayCmd waits 1 second then transitions to game over
func crashDelayCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return CrashDelayMsg(t)
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
				m.crash()
				return m, crashDelayCmd() // Show crash animation for 1 second
			}

			// Speed up gradually as score increases based on difficulty
			initialSpeed, minSpeed, speedDecrease := m.getDifficultySettings()
			newSpeed := initialSpeed - time.Duration(m.score)*time.Duration(speedDecrease)*time.Millisecond
			if newSpeed < minSpeed {
				newSpeed = minSpeed
			}
			m.speed = newSpeed
		}

		return m, tickCmd(m.speed)

	case CrashDelayMsg:
		// Transition from crash animation to game over screen
		if m.state == StateCrashed {
			m.gameOver()
		}
		return m, nil
	}

	return m, nil
}
