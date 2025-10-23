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
	StateExploding                // Explosion animation in progress
)

// SmileyState represents the smiley face button state
type SmileyState int

const (
	SmileyHappy SmileyState = iota     // :) - Normal playing
	SmileySurprised                    // :O - Mouse down on cell
	SmileyDead                         // X_X - Hit a mine
	SmileyCool                         // B) - Won the game
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
	mouseMode  bool // True when using mouse (hides cursor), false when using keyboard

	// Grid rendering boundaries (calculated during View())
	gridStartX int
	gridStartY int

	// Animation state
	explosionCenterX  int // X coordinate of clicked mine
	explosionCenterY  int // Y coordinate of clicked mine
	explosionRadius   int // Current explosion radius
	explosionMaxSteps int // Total animation steps

	// UI state
	smileyState SmileyState // Current smiley face state

	// Double-click detection (for touch-screen flag support)
	lastClickX    int
	lastClickY    int
	lastClickTime time.Time

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

// AnimationTickMsg is sent for explosion animation frames
type AnimationTickMsg time.Time

// Timer command that sends TickMsg every second
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Animation command that sends AnimationTickMsg for explosion animation
func animationTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return AnimationTickMsg(t)
	})
}
