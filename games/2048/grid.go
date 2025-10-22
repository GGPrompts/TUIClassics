package game2048

// move attempts to move tiles in the specified direction
// Returns true if the grid changed
func (m *Model) move(dir Direction) bool {
	moved := false

	switch dir {
	case Left:
		moved = m.moveLeft()
	case Right:
		moved = m.moveRight()
	case Up:
		moved = m.moveUp()
	case Down:
		moved = m.moveDown()
	}

	// If grid changed, spawn new tile and check game state
	if moved {
		m.spawnTile()

		// Check win condition (first time reaching 2048)
		if !m.wonOnce && m.hasValue(2048) {
			m.state = StateWon
			m.wonOnce = true
		}

		// Check game over
		if !m.canMove() {
			m.state = StateGameOver
			if m.score > m.bestScore {
				m.bestScore = m.score
			}
		}
	}

	return moved
}

// moveLeft slides all tiles left
func (m *Model) moveLeft() bool {
	moved := false
	for i := 0; i < 4; i++ {
		if m.slideAndMergeRow(i) {
			moved = true
		}
	}
	return moved
}

// moveRight slides all tiles right
// Strategy: reverse rows, move left, reverse back
func (m *Model) moveRight() bool {
	m.reverseRows()
	moved := m.moveLeft()
	m.reverseRows()
	return moved
}

// moveUp slides all tiles up
// Strategy: transpose, move left, transpose back
func (m *Model) moveUp() bool {
	m.transpose()
	moved := m.moveLeft()
	m.transpose()
	return moved
}

// moveDown slides all tiles down
// Strategy: transpose, reverse, move left, reverse, transpose
func (m *Model) moveDown() bool {
	m.transpose()
	m.reverseRows()
	moved := m.moveLeft()
	m.reverseRows()
	m.transpose()
	return moved
}

// canMove checks if any move is possible
func (m *Model) canMove() bool {
	// Check for empty cells
	if m.hasEmptyCell() {
		return true
	}

	// Check for possible horizontal merges
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			if m.grid[i][j].Value == m.grid[i][j+1].Value {
				return true
			}
		}
	}

	// Check for possible vertical merges
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if m.grid[i][j].Value == m.grid[i+1][j].Value {
				return true
			}
		}
	}

	return false
}
