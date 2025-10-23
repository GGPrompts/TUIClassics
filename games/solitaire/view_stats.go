package solitaire

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
	title := titleStyle.Width(m.termWidth).Render("🏆 SOLITAIRE LEADERBOARD 🏆")
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
		Width(80)

	// Build table content
	var content strings.Builder
	content.WriteString("┌─ High Scores ─────────────────────────────────────────────────────┐\n")
	content.WriteString("│ Period   │ Score  │ Time    │ Moves │ Date                     │\n")
	content.WriteString("├──────────┼────────┼─────────┼───────┼──────────────────────────┤\n")

	// All-time
	allTimeScore := formatSolitaireScore(scores.Solitaire.AllTime.Score)
	allTimeTime := formatSolitaireTime(scores.Solitaire.AllTime.TimeSeconds)
	allTimeMoves := formatSolitaireMoves(scores.Solitaire.AllTime.Moves)
	content.WriteString(fmt.Sprintf("│ All-Time │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		allTimeScore, allTimeTime, allTimeMoves, formatDate(scores.Solitaire.AllTime.Date)))

	// Monthly
	monthScore := formatSolitaireScore(scores.Solitaire.Monthly.Score)
	monthTime := formatSolitaireTime(scores.Solitaire.Monthly.TimeSeconds)
	monthMoves := formatSolitaireMoves(scores.Solitaire.Monthly.Moves)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		monthScore, monthTime, monthMoves, formatDate(scores.Solitaire.Monthly.Date)))

	// Weekly
	weekScore := formatSolitaireScore(scores.Solitaire.Weekly.Score)
	weekTime := formatSolitaireTime(scores.Solitaire.Weekly.TimeSeconds)
	weekMoves := formatSolitaireMoves(scores.Solitaire.Weekly.Moves)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		weekScore, weekTime, weekMoves, formatDate(scores.Solitaire.Weekly.Date)))

	content.WriteString("└──────────┴────────┴─────────┴───────┴──────────────────────────┘")

	// Center the table
	tableContent := tableStyle.Render(content.String())
	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableContent)
}

// formatSolitaireScore formats a solitaire score for display
func formatSolitaireScore(score int) string {
	if score == 0 {
		return "---"
	}
	return fmt.Sprintf("%d", score)
}

// formatSolitaireTime formats time in seconds to MM:SS
func formatSolitaireTime(timeSeconds int) string {
	if timeSeconds == 0 {
		return "--:--"
	}
	minutes := timeSeconds / 60
	seconds := timeSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// formatSolitaireMoves formats moves count for display
func formatSolitaireMoves(moves int) string {
	if moves == 0 {
		return "---"
	}
	return fmt.Sprintf("%d", moves)
}

// formatDate formats a date for display
func formatDate(t time.Time) string {
	if t.IsZero() {
		return "---"
	}
	return t.Format("2006-01-02 15:04")
}
