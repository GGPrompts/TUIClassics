package snake

import "github.com/charmbracelet/lipgloss"

// Snake styles
var (
	headStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Bright green
		Bold(true)

	headEatingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow when eating
		Bold(true)

	headCrashedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red when crashed
		Bold(true)

	bodyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")) // Bright green (same as head)

	foodStyle = lipgloss.NewStyle().
		Bold(true)
)

// UI styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Green
		Bold(true).
		Align(lipgloss.Center)

	menuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")) // Light gray

	scoreStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true)

	scoreHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true).
		Align(lipgloss.Center)

	highScoreStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")). // Gold
		Bold(true)

	gameOverStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red
		Bold(true).
		Align(lipgloss.Center)

	pauseStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")) // Dark gray

	// Difficulty selection styles
	difficultyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")) // Light gray

	difficultyDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")) // Dark gray

	selectedDifficultyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Green
		Bold(true)

	selectedDifficultyDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")) // Light green

	// Separator line below score
	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")) // Green (matches border)

	// Game border - uses Lipgloss rounded border
	gameBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")). // Green border
		Padding(0, 1)
)
