package game2048

import "github.com/charmbracelet/lipgloss"

var (
	// Title and UI styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	scoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Grid styles
	gridStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

// getTileStyle returns the lipgloss style for a tile based on its value
func getTileStyle(value int) lipgloss.Style {
	style := lipgloss.NewStyle().Bold(true)

	switch value {
	case 0:
		return style.Foreground(lipgloss.Color("237")) // Empty
	case 2:
		return style.Foreground(lipgloss.Color("238")) // Dark gray
	case 4:
		return style.Foreground(lipgloss.Color("239"))
	case 8:
		return style.Foreground(lipgloss.Color("208")) // Orange
	case 16:
		return style.Foreground(lipgloss.Color("202"))
	case 32:
		return style.Foreground(lipgloss.Color("196")) // Red
	case 64:
		return style.Foreground(lipgloss.Color("160"))
	case 128:
		return style.Foreground(lipgloss.Color("226")) // Yellow
	case 256:
		return style.Foreground(lipgloss.Color("220"))
	case 512:
		return style.Foreground(lipgloss.Color("214"))
	case 1024:
		return style.Foreground(lipgloss.Color("208"))
	case 2048:
		return style.Foreground(lipgloss.Color("46")) // Green (victory!)
	default:
		return style.Foreground(lipgloss.Color("201")) // Magenta (>2048)
	}
}
