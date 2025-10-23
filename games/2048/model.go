package game2048

import "github.com/GGPrompts/TUIClassics/internal/stats"

// New creates a new 2048 game model
func New() Model {
	return Model{
		state:     StateMenu,
		bestScore: 0,
	}
}

// startGame initializes a new game
func (m *Model) startGame() {
	// Clear grid
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			m.grid[i][j] = Tile{Value: 0}
		}
	}

	// Reset game state
	m.score = 0
	m.wonOnce = false
	m.state = StatePlaying

	// Spawn two initial tiles
	m.spawnTile()
	m.spawnTile()
}

// recordStats records the game result in the stats database
func (m *Model) recordStats() {
	// Load current stats
	scores, err := stats.Load()
	if err != nil {
		return // Silently fail if stats can't be loaded
	}

	// Update stats and get achievements
	achievements := stats.Update2048Score(scores, m.score)

	// Save updated stats
	if err := stats.Save(scores); err != nil {
		return // Silently fail if stats can't be saved
	}

	// Store achievements for display
	m.achievements = make([]string, len(achievements))
	for i, ach := range achievements {
		m.achievements[i] = ach.String()
	}
}
