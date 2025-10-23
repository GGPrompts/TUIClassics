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

	// Calculate vertical centering
	// title(1) + blank(1) + subtitle(1) + blank(1) + enter(1) + instructions(1) + quit(1) = 7 lines
	totalLines := 7
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title (centered horizontally)
	b.WriteString(titleStyle.Width(m.termWidth).Align(lipgloss.Center).Render("2048"))
	b.WriteString("\n\n")

	// Subtitle
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render("Combine tiles to reach 2048!"))
	b.WriteString("\n\n")

	// Help text
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render("Press ENTER to start")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render("I or ? for instructions")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render("Q to quit")))

	return b.String()
}

// renderGame renders the game board
func (m Model) renderGame() string {
	var b strings.Builder

	// Calculate vertical centering
	// title(1) + score(1) + blank(1) + grid_top(1) + grid_rows(4) + grid_separators(3) + grid_bottom(1) + blank(1) + help(1) = 14 lines
	totalLines := 14
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title (centered)
	title := titleStyle.Render("2048")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(title))
	b.WriteString("\n")

	// Score (centered)
	scoreText := fmt.Sprintf("Score: %d  Best: %d", m.score, m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(scoreStyle.Render(scoreText)))
	b.WriteString("\n\n")

	// Grid (centered)
	gridContent := m.renderGridContent()
	gridLines := strings.Split(gridContent, "\n")
	gridWidth := 21 // ┌────┬────┬────┬────┐ = 21 chars
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// Help text (centered)
	helpText := "Arrow keys/WASD/HJKL to move | I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render(helpText)))

	return b.String()
}

// renderGridContent renders just the grid without padding
func (m Model) renderGridContent() string {
	var b strings.Builder

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

	return b.String()
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

	// Calculate vertical centering
	// title(1) + blank(1) + score(1) + best(1) + blank(1) + grid(9) + blank(1) + help(1) = 16 lines
	totalLines := 16
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	winTitle := titleStyle.Render("🎉 YOU WIN! 🎉")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(winTitle))
	b.WriteString("\n\n")

	// Scores (centered)
	scoreText := fmt.Sprintf("Score: %d", m.score)
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(scoreText))
	b.WriteString("\n")
	bestText := fmt.Sprintf("Best: %d", m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(bestText))
	b.WriteString("\n\n")

	// Grid (centered)
	gridContent := m.renderGridContent()
	gridLines := strings.Split(gridContent, "\n")
	gridWidth := 21
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// Help text (centered)
	winHelpText := "C to continue | I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render(winHelpText)))

	return b.String()
}

// renderGameOver renders the game over screen
func (m Model) renderGameOver() string {
	var b strings.Builder

	// Calculate vertical centering
	// title(1) + blank(1) + score(1) + best(1) + blank(1) + grid(9) + blank(1) + help(1) = 16 lines
	totalLines := 16
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	gameOverTitle := titleStyle.Render("GAME OVER")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(gameOverTitle))
	b.WriteString("\n\n")

	// Scores (centered)
	scoreText := fmt.Sprintf("Score: %d", m.score)
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(scoreText))
	b.WriteString("\n")
	bestText := fmt.Sprintf("Best: %d", m.bestScore)
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(bestText))
	b.WriteString("\n\n")

	// Grid (centered)
	gridContent := m.renderGridContent()
	gridLines := strings.Split(gridContent, "\n")
	gridWidth := 21
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// Help text (centered)
	gameOverHelpText := "I/? for help | R to restart | Q to quit"
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render(gameOverHelpText)))

	return b.String()
}

// renderInstructions renders the instructions screen
func (m Model) renderInstructions() string {
	var b strings.Builder

	// Calculate vertical centering
	// Count all lines in the instructions
	// title(1) + blank(1) + objective_header(1) + objective_content(1) + blank(1) +
	// howto_header(1) + howto_content(4) + blank(1) + controls_header(1) + controls_content(4) + blank(1) +
	// example_header(1) + example_content(4) + blank(1) + tip(1) + blank(1) + return(1) = 26 lines
	totalLines := 26
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	title := titleStyle.Render("HOW TO PLAY 2048")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(title))
	b.WriteString("\n\n")

	// Instructions content
	instructionsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)

	// Center all content
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("OBJECTIVE:")))
	b.WriteString("\n")
	objectiveText := instructionsStyle.Render("Combine numbered tiles to reach ") + highlightStyle.Render("2048") + instructionsStyle.Render("!")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(objectiveText))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("HOW TO PLAY:")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("• Use ") + highlightStyle.Render("arrow keys") + instructionsStyle.Render(" to slide ALL tiles")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("• When two tiles with the ") + highlightStyle.Render("same number") + instructionsStyle.Render(" touch, they merge!")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("• After each move, a new tile (") + highlightStyle.Render("2") + instructionsStyle.Render(" or ") + highlightStyle.Render("4") + instructionsStyle.Render(") appears")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("• Keep merging to build bigger numbers!")))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("CONTROLS:")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(highlightStyle.Render("↑ ↓ ← →") + instructionsStyle.Render("  or  ") + highlightStyle.Render("W A S D") + instructionsStyle.Render("  or  ") + highlightStyle.Render("H J K L") + instructionsStyle.Render("  - Move tiles")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(highlightStyle.Render("R") + instructionsStyle.Render(" - Restart game")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(highlightStyle.Render("I") + instructionsStyle.Render(" or ") + highlightStyle.Render("?") + instructionsStyle.Render(" - Show instructions (this screen)")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(highlightStyle.Render("Q") + instructionsStyle.Render(" - Quit game")))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("EXAMPLE:")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(getTileStyle(2).Render(" 2  ") + instructionsStyle.Render(" + ") + getTileStyle(2).Render(" 2  ") + instructionsStyle.Render("  =  ") + getTileStyle(4).Render(" 4  ")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(getTileStyle(4).Render(" 4  ") + instructionsStyle.Render(" + ") + getTileStyle(4).Render(" 4  ") + instructionsStyle.Render("  =  ") + getTileStyle(8).Render(" 8  ")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(getTileStyle(8).Render(" 8  ") + instructionsStyle.Render(" + ") + getTileStyle(8).Render(" 8  ") + instructionsStyle.Render("  =  ") + getTileStyle(16).Render(" 16 ")))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("... and so on!")))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(instructionsStyle.Render("TIP: Keep your highest tile in a corner!")))
	b.WriteString("\n\n")

	// Return instruction
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(helpStyle.Render("Press ENTER, ESC, I, or ? to return")))

	return b.String()
}
