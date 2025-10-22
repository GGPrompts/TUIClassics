package snake

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current game state
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.renderMenu()
	case StatePlaying, StatePaused, StateCrashed:
		return m.renderGame()
	case StateGameOver:
		return m.renderGameOver()
	}
	return ""
}

// renderMenu displays the welcome screen
func (m Model) renderMenu() string {
	var b strings.Builder

	title := titleStyle.Render("🐍 SNAKE")
	b.WriteString(title)
	b.WriteString("\n\n")

	instructions := []string{
		"Classic Snake game - eat food, grow longer, don't hit walls or yourself!",
		"",
		"Controls:",
		"  Arrow Keys  or  WASD  or  HJKL  -  Move",
		"  P / Space  -  Pause",
		"  Q  -  Quit",
		"",
		"Press ENTER to start",
	}

	for _, line := range instructions {
		b.WriteString(menuStyle.Render(line))
		b.WriteString("\n")
	}

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderGame displays the playing field
func (m Model) renderGame() string {
	var b strings.Builder

	// Score header
	scoreLine := fmt.Sprintf("Score: %d    High: %d", m.score, m.highScore)
	b.WriteString(scoreHeaderStyle.Render(scoreLine))
	b.WriteString("\n")

	// Game board
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			point := Point{x, y}

			if point == m.snake[0] {
				// Snake head - different expressions for different states
				if m.state == StateCrashed {
					b.WriteString(headCrashedStyle.Render("😵"))
				} else if m.justAte {
					b.WriteString(headEatingStyle.Render("😮"))
				} else {
					b.WriteString(headStyle.Render("😊"))
				}
			} else if m.isSnakeBody(point) {
				// Snake body - bright green circles
				b.WriteString(bodyStyle.Render("●●"))
			} else if point == m.food {
				// Food - apple emoji
				b.WriteString(foodStyle.Render("🍎"))
			} else {
				// Empty space
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	}

	// Wrap in border
	bordered := gameBorderStyle.Render(b.String())

	// Help text
	var help string
	if m.state == StatePaused {
		help = pauseStyle.Render("⏸  PAUSED - Press P to continue")
	} else {
		help = helpStyle.Render("Press arrow keys to move | P to pause | Q to quit")
	}

	content := lipgloss.JoinVertical(lipgloss.Center, bordered, "", help)

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, content)
}

// renderGameOver displays the game over screen
func (m Model) renderGameOver() string {
	var b strings.Builder

	title := gameOverStyle.Render("GAME OVER!")
	b.WriteString(title)
	b.WriteString("\n\n")

	score := scoreStyle.Render(fmt.Sprintf("Final Score: %d", m.score))
	b.WriteString(score)
	b.WriteString("\n")

	if m.score == m.highScore && m.score > 0 {
		newHighScore := highScoreStyle.Render("🏆 NEW HIGH SCORE! 🏆")
		b.WriteString(newHighScore)
		b.WriteString("\n")
	} else {
		highScore := fmt.Sprintf("High Score: %d", m.highScore)
		b.WriteString(highScore)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("Press R to restart\n")
	b.WriteString("Press M for menu\n")
	b.WriteString("Press Q to quit\n")

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// isSnakeBody checks if a point is part of the snake
func (m Model) isSnakeBody(point Point) bool {
	for _, segment := range m.snake {
		if point == segment {
			return true
		}
	}
	return false
}
