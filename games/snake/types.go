package snake

import "time"

// Direction represents the snake's movement direction
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Point represents a coordinate on the game board
type Point struct {
	X int
	Y int
}

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StateDifficultySelect
	StatePlaying
	StatePaused
	StateCrashed  // Brief pause showing crash animation
	StateGameOver
	StateStats
)

// Difficulty represents game difficulty level
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

// Model holds the game state
type Model struct {
	state              GameState
	previousState      GameState  // For returning from stats view
	difficulty         Difficulty
	selectedDifficulty Difficulty // Cursor position in difficulty select
	snake              []Point    // Head is at index 0
	direction          Direction
	nextDir            Direction // Buffer for next move to prevent missed inputs
	food               Point
	score              int
	highScore          int
	justAte            bool // True for one tick after eating food

	// Stats/Achievements
	achievements []string // Latest achievements to display

	// Game board dimensions
	width  int
	height int
	speed  time.Duration // Time between moves

	// Terminal dimensions
	termWidth  int
	termHeight int
}
