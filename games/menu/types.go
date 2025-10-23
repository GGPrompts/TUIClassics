package menu

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AnimationTickMsg is sent on each animation frame for the landing page
type AnimationTickMsg time.Time

// animationTick returns a command that waits for the next animation frame (60fps)
func animationTick() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return AnimationTickMsg(t)
	})
}

// MenuState represents what the application is currently showing
type MenuState int

const (
	StateMainMenu MenuState = iota // Showing game selection menu
	StateInGame                     // Playing a game
)

// GameInfo represents metadata about an available game
type GameInfo struct {
	Name        string      // Display name
	Description string      // Short description
	Hotkey      string      // Keyboard shortcut (e.g., "m" for minesweeper)
	Model       tea.Model   // The game's model (nil until launched)
	NewFunc     func() tea.Model // Function to create new instance
}

// Model is the main menu model
type Model struct {
	state       MenuState
	games       []GameInfo
	selectedIdx int // Currently selected game in menu

	currentGame tea.Model // The currently running game (if any)

	// Landing page (Windows 95-style launcher)
	landingPage *LandingPage

	// Terminal dimensions
	termWidth  int
	termHeight int

	// Double-click detection
	lastClickTime   time.Time
	lastClickButton int
}
