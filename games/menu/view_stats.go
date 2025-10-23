package menu

import (
	"fmt"
	"strings"
	"time"

	"github.com/GGPrompts/TUIClassics/internal/stats"
	"github.com/charmbracelet/lipgloss"
)

// Tab styles
var (
	tabActiveStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")).
		Background(lipgloss.Color("#404040")).
		Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#7D56F4"))

	tableStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3E4451")).
		Padding(1, 2)
)

// viewHighScores renders the tabbed high scores view
func (m Model) viewHighScores() string {
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

	// Calculate centering
	totalLines := 30 // Approximate
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	title := "🏆 TUI CLASSICS - HIGH SCORES 🏆"
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Width(m.termWidth).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Tab bar
	b.WriteString(m.renderTabBar())
	b.WriteString("\n\n")

	// Tab content based on current tab
	var content string
	switch m.currentTab {
	case 0: // Minesweeper
		content = m.renderMinesweeperStats(scores)
	case 1: // 2048
		content = m.render2048Stats(scores)
	case 2: // Solitaire
		content = m.renderSolitaireStats(scores)
	case 3: // Snake
		content = m.renderSnakeStats(scores)
	default:
		content = m.renderMinesweeperStats(scores)
	}

	b.WriteString(content)
	b.WriteString("\n\n")

	// Help text
	help := "[Tab/Shift+Tab] Switch Tabs | [←/→] Navigate | [ESC] Back to Menu | [Q] Quit"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(m.termWidth).
		Align(lipgloss.Center)
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

// renderTabBar renders the tab navigation bar
func (m Model) renderTabBar() string {
	tabs := []string{"Minesweeper", "2048", "Solitaire", "Snake"}
	var renderedTabs []string

	for i, tab := range tabs {
		if m.currentTab == i {
			renderedTabs = append(renderedTabs, tabActiveStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, tabInactiveStyle.Render(tab))
		}
	}

	tabsRow := strings.Join(renderedTabs, " ")
	centered := lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tabsRow)
	return tabBarStyle.Width(m.termWidth).Render(centered)
}

// renderMinesweeperStats renders Minesweeper leaderboards (all 3 difficulties)
func (m Model) renderMinesweeperStats(scores *stats.ScoreData) string {
	var b strings.Builder

	// Easy
	b.WriteString(m.renderDifficultyTable("Easy", scores.Minesweeper.Easy))
	b.WriteString("\n\n")

	// Medium
	b.WriteString(m.renderDifficultyTable("Medium", scores.Minesweeper.Medium))
	b.WriteString("\n\n")

	// Hard
	b.WriteString(m.renderDifficultyTable("Hard", scores.Minesweeper.Hard))

	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(b.String())
}

// renderDifficultyTable renders a Minesweeper difficulty table
func (m Model) renderDifficultyTable(difficulty string, diffStats stats.DifficultyStats) string {
	var content strings.Builder
	content.WriteString(fmt.Sprintf("┌─ %s ────────────────────────────────────────────────┐\n", difficulty))
	content.WriteString("│ Period   │ Time      │ Date                            │\n")
	content.WriteString("├──────────┼───────────┼─────────────────────────────────┤\n")

	// All-time
	allTimeStr := formatTimeScore(diffStats.AllTime.TimeSeconds)
	content.WriteString(fmt.Sprintf("│ All-Time │ %-9s │ %-31s │\n", allTimeStr, formatDate(diffStats.AllTime.Date)))

	// Monthly
	monthStr := formatTimeScore(diffStats.Monthly.TimeSeconds)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-9s │ %-31s │\n", monthStr, formatDate(diffStats.Monthly.Date)))

	// Weekly
	weekStr := formatTimeScore(diffStats.Weekly.TimeSeconds)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-9s │ %-31s │\n", weekStr, formatDate(diffStats.Weekly.Date)))

	content.WriteString("└──────────┴───────────┴─────────────────────────────────┘")

	return tableStyle.Width(66).Render(content.String())
}

// render2048Stats renders 2048 leaderboards
func (m Model) render2048Stats(scores *stats.ScoreData) string {
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

	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableStyle.Width(70).Render(content.String()))
}

// renderSolitaireStats renders Solitaire leaderboards
func (m Model) renderSolitaireStats(scores *stats.ScoreData) string {
	var content strings.Builder
	content.WriteString("┌─ High Scores ────────────────────────────────────────────────┐\n")
	content.WriteString("│ Period   │ Score  │ Time    │ Moves │ Date                     │\n")
	content.WriteString("├──────────┼────────┼─────────┼───────┼──────────────────────────┤\n")

	// All-time
	allTimeScore := formatScore(scores.Solitaire.AllTime.Score)
	allTimeTime := formatSolitaireTime(scores.Solitaire.AllTime.TimeSeconds)
	allTimeMoves := formatMoves(scores.Solitaire.AllTime.Moves)
	content.WriteString(fmt.Sprintf("│ All-Time │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		allTimeScore, allTimeTime, allTimeMoves, formatDate(scores.Solitaire.AllTime.Date)))

	// Monthly
	monthScore := formatScore(scores.Solitaire.Monthly.Score)
	monthTime := formatSolitaireTime(scores.Solitaire.Monthly.TimeSeconds)
	monthMoves := formatMoves(scores.Solitaire.Monthly.Moves)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		monthScore, monthTime, monthMoves, formatDate(scores.Solitaire.Monthly.Date)))

	// Weekly
	weekScore := formatScore(scores.Solitaire.Weekly.Score)
	weekTime := formatSolitaireTime(scores.Solitaire.Weekly.TimeSeconds)
	weekMoves := formatMoves(scores.Solitaire.Weekly.Moves)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-6s │ %-7s │ %-5s │ %-24s │\n",
		weekScore, weekTime, weekMoves, formatDate(scores.Solitaire.Weekly.Date)))

	content.WriteString("└──────────┴────────┴─────────┴───────┴──────────────────────────┘")

	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableStyle.Width(80).Render(content.String()))
}

// renderSnakeStats renders Snake leaderboards
func (m Model) renderSnakeStats(scores *stats.ScoreData) string {
	var content strings.Builder
	content.WriteString("┌─ High Scores ──────────────────────────────────────────┐\n")
	content.WriteString("│ Period   │ Score     │ Date                            │\n")
	content.WriteString("├──────────┼───────────┼─────────────────────────────────┤\n")

	// All-time
	allTimeStr := formatScore(scores.Snake.AllTime.Score)
	content.WriteString(fmt.Sprintf("│ All-Time │ %-9s │ %-31s │\n", allTimeStr, formatDate(scores.Snake.AllTime.Date)))

	// Monthly
	monthStr := formatScore(scores.Snake.Monthly.Score)
	content.WriteString(fmt.Sprintf("│ Monthly  │ %-9s │ %-31s │\n", monthStr, formatDate(scores.Snake.Monthly.Date)))

	// Weekly
	weekStr := formatScore(scores.Snake.Weekly.Score)
	content.WriteString(fmt.Sprintf("│ Weekly   │ %-9s │ %-31s │\n", weekStr, formatDate(scores.Snake.Weekly.Date)))

	content.WriteString("└──────────┴───────────┴─────────────────────────────────┘")

	return lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(tableStyle.Width(70).Render(content.String()))
}

// Formatting helpers

func formatTimeScore(timeSeconds int) string {
	if timeSeconds == 0 {
		return "--:--"
	}
	minutes := timeSeconds / 60
	seconds := timeSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func formatScore(score int) string {
	if score == 0 {
		return "---"
	}
	return fmt.Sprintf("%d", score)
}

func formatSolitaireTime(timeSeconds int) string {
	if timeSeconds == 0 {
		return "--:--"
	}
	minutes := timeSeconds / 60
	seconds := timeSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func formatMoves(moves int) string {
	if moves == 0 {
		return "---"
	}
	return fmt.Sprintf("%d", moves)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return "---"
	}
	return t.Format("2006-01-02 15:04")
}
