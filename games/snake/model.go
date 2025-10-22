package snake

import "time"

// New creates a new Snake game model
func New() Model {
	return Model{
		state:     StateMenu,
		width:     24, // Reduced from 30 for better aspect ratio
		height:    20,
		speed:     200 * time.Millisecond, // Slower starting speed
		highScore: 0,
	}
}

// startGame initializes a new game
func (m *Model) startGame() {
	// Place snake in the middle of the board
	centerX := m.width / 2
	centerY := m.height / 2

	// Create initial snake with 3 segments
	m.snake = []Point{
		{centerX, centerY},
		{centerX - 1, centerY},
		{centerX - 2, centerY},
	}

	m.direction = Right
	m.nextDir = Right
	m.score = 0
	m.spawnFood()
	m.state = StatePlaying
}
