package game2048

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
