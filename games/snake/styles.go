package snake

import "github.com/charmbracelet/lipgloss"

// Snake styles
var (
	headStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")). // Bright green
		Bold(true)

	bodyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")) // Light green

	foodStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red
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
)
