package menu

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#00BFFF")  // Deep Sky Blue
	accentColor    = lipgloss.Color("#FFD700")  // Gold
	textColor      = lipgloss.Color("#FFFFFF")  // White
	mutedColor     = lipgloss.Color("#808080")  // Gray
	selectedBg     = lipgloss.Color("#1E3A8A")  // Dark Blue
	disabledColor  = lipgloss.Color("#404040")  // Dark Gray

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			MarginBottom(2)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2).
			PaddingRight(2)

	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Background(selectedBg).
				Bold(true).
				PaddingLeft(2).
				PaddingRight(2)

	disabledMenuItemStyle = lipgloss.NewStyle().
				Foreground(disabledColor).
				PaddingLeft(2).
				PaddingRight(2)

	// Description styles
	descriptionStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				MarginLeft(4)

	// Footer/hint styles
	hintStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(2)
)
