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
	StatePlaying
	StatePaused
	StateGameOver
)

// Model holds the game state
type Model struct {
	state     GameState
	snake     []Point   // Head is at index 0
	direction Direction
	nextDir   Direction // Buffer for next move to prevent missed inputs
	food      Point
	score     int
	highScore int
	justAte   bool // True for one tick after eating food

	// Game board dimensions
	width  int
	height int
	speed  time.Duration // Time between moves

	// Terminal dimensions
	termWidth  int
	termHeight int
}
