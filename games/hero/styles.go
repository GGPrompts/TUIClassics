package hero

import "github.com/charmbracelet/lipgloss"

// Lane colors (neon style for each lane)
var laneColors = []lipgloss.Color{
	lipgloss.Color("46"),  // Green (A)
	lipgloss.Color("196"), // Red (S)
	lipgloss.Color("226"), // Yellow (D)
	lipgloss.Color("21"),  // Blue (F)
	lipgloss.Color("201"), // Magenta (J)
}

// Hit result colors
var (
	perfectColor = lipgloss.Color("220") // Gold
	goodColor    = lipgloss.Color("46")  // Green
	okColor      = lipgloss.Color("226") // Yellow
	missColor    = lipgloss.Color("196") // Red
)

// UI colors
var (
	titleColor = lipgloss.Color("201") // Magenta
	textColor  = lipgloss.Color("255") // White
	dimColor   = lipgloss.Color("240") // Gray
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			Align(lipgloss.Center)

	laneHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center)

	noteStyle = lipgloss.NewStyle().
			Bold(true)

	hitZoneStyle = lipgloss.NewStyle().
			Bold(true)

	scoreStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			Align(lipgloss.Center)

	comboStyle = lipgloss.NewStyle().
			Foreground(perfectColor).
			Bold(true).
			Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center)

	// Base lane style (will be customized per lane with different border colors)
	baseLaneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, true, true, true)
)

// getLaneColor returns the color for a specific lane
func getLaneColor(lane int) lipgloss.Color {
	if lane >= 0 && lane < len(laneColors) {
		return laneColors[lane]
	}
	return textColor
}

// getLaneBorderStyle returns a styled border for a specific lane
func getLaneBorderStyle(lane int) lipgloss.Style {
	return baseLaneStyle.Copy().
		BorderForeground(getLaneColor(lane))
}

// getHitResultColor returns the color for a hit result
func getHitResultColor(result HitResult) lipgloss.Color {
	switch result {
	case HitPerfect:
		return perfectColor
	case HitGood:
		return goodColor
	case HitOK:
		return okColor
	default:
		return missColor
	}
}

// getHitResultText returns the text for a hit result
func getHitResultText(result HitResult) string {
	switch result {
	case HitPerfect:
		return " PERFECT! "
	case HitGood:
		return "GOOD!"
	case HitOK:
		return "OK"
	case HitMiss:
		return "MISS"
	default:
		return ""
	}
}
