package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

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

	// Terminal dimensions
	termWidth  int
	termHeight int
}
