package game2048

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderView renders the appropriate view based on game state
func renderView(m Model) string {
	switch m.state {
	case StateMenu:
		return m.renderMenu()
	case StatePlaying:
		return m.renderGame()
	case StateWon:
		return m.renderWin()
	case StateGameOver:
		return m.renderGameOver()
	case StateInstructions:
		return m.renderInstructions()
	}
	return ""
}

// renderMenu renders the starting menu
func (m Model) renderMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("2048"))
	b.WriteString("\n\n")
	b.WriteString("Combine tiles to reach 2048!\n\n")
	b.WriteString(helpStyle.Render("Press ENTER to start"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("I or ? for instructions"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Q to quit"))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderGame renders the game board
func (m Model) renderGame() string {
	var b strings.Builder

	// Title (centered)
	title := titleStyle.Render("2048")
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(title))
	b.WriteString("\n")

	// Score (centered to match grid width)
	scoreText := fmt.Sprintf("Score: %d  Best: %d", m.score, m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(scoreStyle.Render(scoreText)))
	b.WriteString("\n\n")

	// Grid
	b.WriteString(gridStyle.Render("┌────┬────┬────┬────┐"))
	b.WriteString("\n")

	for i := 0; i < 4; i++ {
		b.WriteString(gridStyle.Render("│"))
		for j := 0; j < 4; j++ {
			b.WriteString(m.renderTile(m.grid[i][j]))
			b.WriteString(gridStyle.Render("│"))
		}
		b.WriteString("\n")

		if i < 3 {
			b.WriteString(gridStyle.Render("├────┼────┼────┼────┤"))
			b.WriteString("\n")
		}
	}

	b.WriteString(gridStyle.Render("└────┴────┴────┴────┘"))
	b.WriteString("\n\n")

	// Help text (centered)
	helpText := "Arrow keys/WASD/HJKL to move | I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(75).Align(lipgloss.Center).Render(helpStyle.Render(helpText)))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderTile renders a single tile with centered number
func (m Model) renderTile(tile Tile) string {
	if tile.Value == 0 {
		return "    "
	}

	style := getTileStyle(tile.Value)

	// Center the number within the 4-character tile space
	text := fmt.Sprintf("%d", tile.Value)
	switch len(text) {
	case 1:
		text = " " + text + "  " // " 2  "
	case 2:
		text = " " + text + " "  // " 16 "
	case 3:
		text = text + " "        // "128 "
	default:
		// 4 digits like 2048 fit perfectly
	}

	return style.Render(text)
}

// renderWin renders the win screen (when reaching 2048)
func (m Model) renderWin() string {
	var b strings.Builder

	// Title
	winTitle := titleStyle.Render("🎉 YOU WIN! 🎉")
	b.WriteString(lipgloss.NewStyle().Width(30).Align(lipgloss.Center).Render(winTitle))
	b.WriteString("\n\n")

	// Scores (centered)
	scoreText := fmt.Sprintf("Score: %d", m.score)
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(scoreText))
	b.WriteString("\n")
	bestText := fmt.Sprintf("Best: %d", m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(bestText))
	b.WriteString("\n\n")

	// Show the grid
	b.WriteString(gridStyle.Render("┌────┬────┬────┬────┐"))
	b.WriteString("\n")

	for i := 0; i < 4; i++ {
		b.WriteString(gridStyle.Render("│"))
		for j := 0; j < 4; j++ {
			b.WriteString(m.renderTile(m.grid[i][j]))
			b.WriteString(gridStyle.Render("│"))
		}
		b.WriteString("\n")

		if i < 3 {
			b.WriteString(gridStyle.Render("├────┼────┼────┼────┤"))
			b.WriteString("\n")
		}
	}

	b.WriteString(gridStyle.Render("└────┴────┴────┴────┘"))
	b.WriteString("\n\n")

	// Help text (centered)
	winHelpText := "C to continue | I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(60).Align(lipgloss.Center).Render(helpStyle.Render(winHelpText)))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderGameOver renders the game over screen
func (m Model) renderGameOver() string {
	var b strings.Builder

	// Title
	gameOverTitle := titleStyle.Render("GAME OVER")
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(gameOverTitle))
	b.WriteString("\n\n")

	// Scores (centered)
	scoreText := fmt.Sprintf("Score: %d", m.score)
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(scoreText))
	b.WriteString("\n")
	bestText := fmt.Sprintf("Best: %d", m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(21).Align(lipgloss.Center).Render(bestText))
	b.WriteString("\n\n")

	// Show the grid
	b.WriteString(gridStyle.Render("┌────┬────┬────┬────┐"))
	b.WriteString("\n")

	for i := 0; i < 4; i++ {
		b.WriteString(gridStyle.Render("│"))
		for j := 0; j < 4; j++ {
			b.WriteString(m.renderTile(m.grid[i][j]))
			b.WriteString(gridStyle.Render("│"))
		}
		b.WriteString("\n")

		if i < 3 {
			b.WriteString(gridStyle.Render("├────┼────┼────┼────┤"))
			b.WriteString("\n")
		}
	}

	b.WriteString(gridStyle.Render("└────┴────┴────┴────┘"))
	b.WriteString("\n\n")

	// Help text (centered)
	gameOverHelpText := "I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(45).Align(lipgloss.Center).Render(helpStyle.Render(gameOverHelpText)))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderInstructions renders the instructions screen
func (m Model) renderInstructions() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("HOW TO PLAY 2048")
	b.WriteString(lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(title))
	b.WriteString("\n\n")

	// Instructions content
	instructionsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)

	b.WriteString(instructionsStyle.Render("OBJECTIVE:"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  Combine numbered tiles to reach "))
	b.WriteString(highlightStyle.Render("2048"))
	b.WriteString(instructionsStyle.Render("!"))
	b.WriteString("\n\n")

	b.WriteString(instructionsStyle.Render("HOW TO PLAY:"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  • Use "))
	b.WriteString(highlightStyle.Render("arrow keys"))
	b.WriteString(instructionsStyle.Render(" to slide ALL tiles"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  • When two tiles with the "))
	b.WriteString(highlightStyle.Render("same number"))
	b.WriteString(instructionsStyle.Render(" touch, they merge!"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  • After each move, a new tile ("))
	b.WriteString(highlightStyle.Render("2"))
	b.WriteString(instructionsStyle.Render(" or "))
	b.WriteString(highlightStyle.Render("4"))
	b.WriteString(instructionsStyle.Render(") appears"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  • Keep merging to build bigger numbers!"))
	b.WriteString("\n\n")

	b.WriteString(instructionsStyle.Render("CONTROLS:"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(highlightStyle.Render("↑ ↓ ← →"))
	b.WriteString(instructionsStyle.Render("  or  "))
	b.WriteString(highlightStyle.Render("W A S D"))
	b.WriteString(instructionsStyle.Render("  or  "))
	b.WriteString(highlightStyle.Render("H J K L"))
	b.WriteString(instructionsStyle.Render("  - Move tiles"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(highlightStyle.Render("R"))
	b.WriteString(instructionsStyle.Render("                     - Restart game"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(highlightStyle.Render("I"))
	b.WriteString(instructionsStyle.Render(" or "))
	b.WriteString(highlightStyle.Render("?"))
	b.WriteString(instructionsStyle.Render("                - Show instructions (this screen)"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(highlightStyle.Render("Q"))
	b.WriteString(instructionsStyle.Render("                     - Quit game"))
	b.WriteString("\n\n")

	b.WriteString(instructionsStyle.Render("EXAMPLE:"))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(getTileStyle(2).Render(" 2  "))
	b.WriteString(instructionsStyle.Render(" + "))
	b.WriteString(getTileStyle(2).Render(" 2  "))
	b.WriteString(instructionsStyle.Render("  =  "))
	b.WriteString(getTileStyle(4).Render(" 4  "))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(getTileStyle(4).Render(" 4  "))
	b.WriteString(instructionsStyle.Render(" + "))
	b.WriteString(getTileStyle(4).Render(" 4  "))
	b.WriteString(instructionsStyle.Render("  =  "))
	b.WriteString(getTileStyle(8).Render(" 8  "))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  "))
	b.WriteString(getTileStyle(8).Render(" 8  "))
	b.WriteString(instructionsStyle.Render(" + "))
	b.WriteString(getTileStyle(8).Render(" 8  "))
	b.WriteString(instructionsStyle.Render("  =  "))
	b.WriteString(getTileStyle(16).Render(" 16 "))
	b.WriteString("\n")
	b.WriteString(instructionsStyle.Render("  ... and so on!"))
	b.WriteString("\n\n")

	b.WriteString(instructionsStyle.Render("TIP: "))
	b.WriteString(instructionsStyle.Render("Keep your highest tile in a corner!"))
	b.WriteString("\n\n")

	// Return instruction
	b.WriteString(helpStyle.Render("Press ENTER, ESC, I, or ? to return"))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}
