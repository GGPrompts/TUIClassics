package snake

import (
	"time"

	"github.com/GGPrompts/TUIClassics/internal/stats"
)

// New creates a new Snake game model
func New() Model {
	return Model{
		state:              StateMenu,
		difficulty:         Medium,
		selectedDifficulty: Medium,
		width:              24, // Reduced from 30 for better aspect ratio
		height:             20,
		highScore:          0,
	}
}

// getDifficultySettings returns initial speed and speed progression for difficulty
func (m *Model) getDifficultySettings() (initialSpeed, minSpeed time.Duration, speedDecrease int) {
	switch m.difficulty {
	case Easy:
		// Responsive start, slower max speed, gentle progression
		return 180 * time.Millisecond, 100 * time.Millisecond, 2
	case Medium:
		return 200 * time.Millisecond, 80 * time.Millisecond, 4
	case Hard:
		return 150 * time.Millisecond, 50 * time.Millisecond, 5
	default:
		return 200 * time.Millisecond, 80 * time.Millisecond, 4
	}
}

// startGame initializes a new game
func (m *Model) startGame() {
	// Set difficulty from selection
	m.difficulty = m.selectedDifficulty

	// Set speed based on difficulty
	initialSpeed, _, _ := m.getDifficultySettings()
	m.speed = initialSpeed

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

// recordStats records the game result in the stats database
func (m *Model) recordStats() {
	// Load current stats
	scores, err := stats.Load()
	if err != nil {
		return // Silently fail if stats can't be loaded
	}

	// Update stats and get achievements
	achievements := stats.UpdateSnakeScore(scores, m.score)

	// Save updated stats
	if err := stats.Save(scores); err != nil {
		return // Silently fail if save fails
	}

	// Store achievements for display
	m.achievements = make([]string, len(achievements))
	for i, ach := range achievements {
		m.achievements[i] = string(ach)
	}
}
