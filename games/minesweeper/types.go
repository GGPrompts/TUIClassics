package minesweeper

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Difficulty represents game difficulty levels
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyMedium
	DifficultyHard
	DifficultyCustom
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota   // Showing difficulty selection
	StatePlaying                  // Game in progress
	StateWon                      // Player won
	StateLost                     // Player hit a mine
	StatePaused                   // Game paused
)

// Cell represents a single cell in the minesweeper grid
type Cell struct {
	IsMine      bool // Is this cell a mine?
	IsRevealed  bool // Has this cell been revealed?
	IsFlagged   bool // Has user flagged this as mine?
	Adjacent    int  // Number of adjacent mines (0-8)
}

// Model is the main game model
type Model struct {
	// Game state
	state     GameState
	grid      [][]Cell
	width     int      // Grid width
	height    int      // Grid height
	mineCount int      // Total mines

	// Difficulty
	difficulty Difficulty

	// Game progress
	startTime     time.Time
	elapsedTime   time.Duration
	flagsPlaced   int
	cellsRevealed int
	firstClick    bool // Track if this is first click (ensure safe start)

	// UI state
	termWidth  int
	termHeight int
	cursorX    int // For keyboard navigation
	cursorY    int

	// High scores
	bestTime time.Duration
}

// Difficulty configurations
var difficultyConfigs = map[Difficulty]struct {
	width     int
	height    int
	mineCount int
	name      string
}{
	DifficultyEasy:   {width: 8, height: 8, mineCount: 10, name: "Easy"},
	DifficultyMedium: {width: 16, height: 16, mineCount: 40, name: "Medium"},
	DifficultyHard:   {width: 30, height: 16, mineCount: 99, name: "Hard"},
}

// TickMsg is sent every second to update the timer
type TickMsg time.Time

// Timer command that sends TickMsg every second
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
