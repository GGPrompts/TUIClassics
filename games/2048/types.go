package game2048

import tea "github.com/charmbracelet/bubbletea"

// Direction represents movement direction
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateWon
	StateGameOver
	StateInstructions
)

// Tile represents a single tile on the grid
type Tile struct {
	Value  int
	Merged bool // Tracks if tile merged this turn (prevent double-merge)
}

// Model holds the game state
type Model struct {
	state         GameState
	previousState GameState // For returning from instructions
	grid          [4][4]Tile
	score         int
	bestScore     int
	wonOnce       bool // True if reached 2048 this game (prevent multiple win screens)

	// Terminal dimensions
	termWidth  int
	termHeight int
}

// Bubbletea interface methods
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return updateGame(m, msg)
}

func (m Model) View() string {
	return renderView(m)
}
