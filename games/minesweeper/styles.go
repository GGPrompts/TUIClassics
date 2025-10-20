package minesweeper

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette - classic minesweeper colors
var (
	// Number colors (classic Windows minesweeper)
	color1 = lipgloss.Color("#0000FF") // Blue
	color2 = lipgloss.Color("#008000") // Green
	color3 = lipgloss.Color("#FF0000") // Red
	color4 = lipgloss.Color("#000080") // Dark Blue
	color5 = lipgloss.Color("#800000") // Maroon
	color6 = lipgloss.Color("#008080") // Teal
	color7 = lipgloss.Color("#000000") // Black
	color8 = lipgloss.Color("#808080") // Gray

	// UI colors
	colorMine     = lipgloss.Color("#FF0000") // Red for mines
	colorFlag     = lipgloss.Color("#FF6B6B") // Bright red for flags
	colorUnknown  = lipgloss.Color("#5C6370") // Gray for unrevealed
	colorRevealed = lipgloss.Color("#E0E0E0") // Light gray for revealed
	colorCursor   = lipgloss.Color("#61AFEF") // Blue for cursor

	colorBorder     = lipgloss.Color("#3E4451")
	colorTitle      = lipgloss.Color("#61AFEF")
	colorSuccess    = lipgloss.Color("#98C379")
	colorError      = lipgloss.Color("#E06C75")
	colorWarning    = lipgloss.Color("#E5C07B")
	colorDimmed     = lipgloss.Color("#5C6370")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Bold(true).
			Align(lipgloss.Center)

	statsStyle = lipgloss.NewStyle().
			Foreground(colorDimmed)

	gridBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	// Cell styles
	unknownCellStyle = lipgloss.NewStyle().
				Foreground(colorUnknown).
				Bold(true)

	revealedCellStyle = lipgloss.NewStyle().
				Foreground(colorRevealed)

	mineCellStyle = lipgloss.NewStyle().
			Foreground(colorMine).
			Bold(true)

	flagCellStyle = lipgloss.NewStyle().
			Foreground(colorFlag).
			Bold(true)

	cursorCellStyle = lipgloss.NewStyle().
			Foreground(colorCursor).
			Background(lipgloss.Color("#2C323C")).
			Bold(true)

	// Number styles (1-8)
	numberStyles = []lipgloss.Style{
		lipgloss.NewStyle(), // 0 - not used
		lipgloss.NewStyle().Foreground(color1).Bold(true), // 1
		lipgloss.NewStyle().Foreground(color2).Bold(true), // 2
		lipgloss.NewStyle().Foreground(color3).Bold(true), // 3
		lipgloss.NewStyle().Foreground(color4).Bold(true), // 4
		lipgloss.NewStyle().Foreground(color5).Bold(true), // 5
		lipgloss.NewStyle().Foreground(color6).Bold(true), // 6
		lipgloss.NewStyle().Foreground(color7).Bold(true), // 7
		lipgloss.NewStyle().Foreground(color8).Bold(true), // 8
	}

	// Game state messages
	winStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true).
			Align(lipgloss.Center)

	loseStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true).
			Align(lipgloss.Center)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2).
			Align(lipgloss.Center)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDimmed).
			Align(lipgloss.Center)
)

// getNumberStyle returns the style for a number (1-8)
func getNumberStyle(n int) lipgloss.Style {
	if n >= 1 && n <= 8 {
		return numberStyles[n]
	}
	return lipgloss.NewStyle()
}
