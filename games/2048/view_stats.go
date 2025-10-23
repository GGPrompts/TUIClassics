package game2048

import (
	"fmt"
	"strings"
	"time"

	"github.com/GGPrompts/TUIClassics/internal/stats"
	"github.com/charmbracelet/lipgloss"
)

// renderStats renders the statistics/leaderboard screen
func (m Model) renderStats() string {
	var b strings.Builder

	// Load stats
	scores, err := stats.Load()
	if err != nil {
		return lipgloss.NewStyle().
			Width(m.termWidth).
			Height(m.termHeight).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Error loading stats: " + err.Error())
	}

	// Calculate vertical centering
	totalLines := 20 // Approximate - title + table + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	title := titleStyle.Width(m.termWidth).Render("🏆 2048 LEADERBOARD 🏆")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Stats table
	b.WriteString(m.renderScoreTable(scores))
	b.WriteString("\n\n")

	// Help text
	help := "[ESC/S] Back to Menu"
	b.WriteString(helpStyle.Width(m.termWidth).Render(help))

	return b.String()
}

// renderScoreTable renders the score table
func (m Model) renderScoreTable(scores *stats.ScoreData) string {
	// Table style
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3E4451")).
		Padding(1, 2).
		Width(70)

	// Build table content
	var content strings.Builder
	content.WriteString("┌─ High Scores ──────────────────────────────────────────┐\n")
	content.WriteString("│ Period   │ Score     │ Date                            │\n")
	content.WriteString("├──────────┼───────────┼─────────────────────────────────┤\n")

	// All-time
	allTimeStr := formatScore(scores.Game2048.AllTime.Score)
	content.WriteString(fmt.Sprintf("│ All-Time │ %-9s │ %-31s │\n", allTimeStr, formatDate(scores.Game2048.AllTime.Date)))

	// Monthly
	monthStr := formatScore(scores.Game2048.Monthly.Score)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-9s │ %-31s │\n", monthStr, formatDate(scores.Game2048.Monthly.Date)))

	// Weekly
	weekStr := formatScore(scores.Game2048.Weekly.Score)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-9s │ %-31s │\n", weekStr, formatDate(scores.Game2048.Weekly.Date)))

	content.WriteString("└──────────┴───────────┴─────────────────────────────────┘")

	// Center the table
	tableContent := tableStyle.Render(content.String())
	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableContent)
}

// formatScore formats a score for display
func formatScore(score int) string {
	if score == 0 {
		return "---"
	}
	return fmt.Sprintf("%d", score)
}

// formatDate formats a date for display
func formatDate(t time.Time) string {
	if t.IsZero() {
		return "---"
	}
	return t.Format("2006-01-02 15:04")
}
