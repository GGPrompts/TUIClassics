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
	case StatePlaying, StatePaused:
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

	// Top border with score
	scoreLine := fmt.Sprintf("  Score: %d    High: %d", m.score, m.highScore)
	padding := m.width*2 - len(scoreLine) + 2
	if padding < 0 {
		padding = 0
	}

	b.WriteString("┌" + strings.Repeat("─", m.width*2) + "┐\n")
	b.WriteString("│" + scoreLine + strings.Repeat(" ", padding) + "│\n")
	b.WriteString("├" + strings.Repeat("─", m.width*2) + "┤\n")

	// Game board
	for y := 0; y < m.height; y++ {
		b.WriteString("│")
		for x := 0; x < m.width; x++ {
			point := Point{x, y}

			if point == m.snake[0] {
				// Snake head
				b.WriteString(headStyle.Render("●●"))
			} else if m.isSnakeBody(point) {
				// Snake body
				b.WriteString(bodyStyle.Render("●●"))
			} else if point == m.food {
				// Food
				b.WriteString(foodStyle.Render("◆◆"))
			} else {
				// Empty space
				b.WriteString("  ")
			}
		}
		b.WriteString("│\n")
	}

	// Bottom border
	b.WriteString("└" + strings.Repeat("─", m.width*2) + "┘\n")

	// Help text
	if m.state == StatePaused {
		b.WriteString("\n")
		b.WriteString(pauseStyle.Render("⏸  PAUSED - Press P to continue"))
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press arrow keys to move | P to pause | Q to quit"))
		b.WriteString("\n")
	}

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
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
