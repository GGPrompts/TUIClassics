package solitaire

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	redColor   = lipgloss.Color("#FF0000")
	blackColor = lipgloss.Color("#000000")
	greenFelt  = lipgloss.Color("#2D5016")
	grayColor  = lipgloss.Color("#808080")
	whiteColor = lipgloss.Color("#FFFFFF")
	goldColor  = lipgloss.Color("#FFD700")
)

// Card styles (5 wide x 3 tall content - better proportions)
var (
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Background(whiteColor).
			Width(5).   // Force consistent width (content area)
			Height(3)   // Force consistent height (content area)

	cardBackStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Background(lipgloss.Color("#0052CC")).
			Foreground(whiteColor).
			Width(5).
			Height(3)

	emptyPileStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(grayColor).
			Width(5).
			Height(3)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(goldColor).
				Background(whiteColor).
				Bold(true).
				Width(5).
				Height(3)

	cursorCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Background(whiteColor).
			Width(5).
			Height(3)
)

// UI styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(goldColor)

	statsStyle = lipgloss.NewStyle().
			Foreground(whiteColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(grayColor).
			MarginTop(1)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(goldColor).
			Padding(1, 2).
			MarginTop(2)

	winStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(goldColor).
			Background(greenFelt).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(goldColor)
)
