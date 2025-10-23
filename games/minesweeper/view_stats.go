package minesweeper

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
	totalLines := 40 // Approximate - title + tables + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	title := titleStyle.Width(m.termWidth).Render("🏆 MINESWEEPER LEADERBOARD 🏆")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Easy difficulty table
	b.WriteString(m.renderDifficultyTable("Easy", scores.Minesweeper.Easy))
	b.WriteString("\n\n")

	// Medium difficulty table
	b.WriteString(m.renderDifficultyTable("Medium", scores.Minesweeper.Medium))
	b.WriteString("\n\n")

	// Hard difficulty table
	b.WriteString(m.renderDifficultyTable("Hard", scores.Minesweeper.Hard))
	b.WriteString("\n\n")

	// Help text
	help := "[ESC/S] Back to Menu"
	b.WriteString(helpStyle.Width(m.termWidth).Render(help))

	return b.String()
}

// renderDifficultyTable renders a table for one difficulty level
func (m Model) renderDifficultyTable(difficulty string, diffStats stats.DifficultyStats) string {
	// Table header
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2).
		Width(70)

	// Build table content
	var content strings.Builder
	content.WriteString(fmt.Sprintf("┌─ %s Difficulty ─────────────────────────────────────┐\n", difficulty))
	content.WriteString("│ Period   │ Time      │ Date                            │\n")
	content.WriteString("├──────────┼───────────┼─────────────────────────────────┤\n")

	// All-time
	allTimeStr := formatTimeScore(diffStats.AllTime, "")
	content.WriteString(fmt.Sprintf("│ All-Time │ %-9s │ %-31s │\n", allTimeStr, formatDate(diffStats.AllTime.Date)))

	// Monthly
	monthStr := formatTimeScore(diffStats.Monthly, diffStats.Monthly.Period)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-9s │ %-31s │\n", monthStr, formatDate(diffStats.Monthly.Date)))

	// Weekly
	weekStr := formatTimeScore(diffStats.Weekly, diffStats.Weekly.Period)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-9s │ %-31s │\n", weekStr, formatDate(diffStats.Weekly.Date)))

	content.WriteString("└──────────┴───────────┴─────────────────────────────────┘")

	// Center the table
	tableContent := tableStyle.Render(content.String())
	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableContent)
}

// formatTimeScore formats a time score for display
func formatTimeScore(score stats.TimeScore, period string) string {
	if score.TimeSeconds == 0 {
		return "--:--"
	}
	minutes := score.TimeSeconds / 60
	seconds := score.TimeSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// formatDate formats a date for display
func formatDate(t time.Time) string {
	if t.IsZero() {
		return "---"
	}
	return t.Format("2006-01-02 15:04")
}
