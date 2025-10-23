package snake

// wouldCollide checks if moving to a specific position would cause collision
func (m *Model) wouldCollide(pos Point) bool {
	// Wall collision
	if pos.X < 0 || pos.X >= m.width ||
		pos.Y < 0 || pos.Y >= m.height {
		return true
	}

	// Self collision (check if position hits any body segment)
	// Note: We check all segments because the tail will move away
	for i := 0; i < len(m.snake)-1; i++ {
		if pos == m.snake[i] {
			return true
		}
	}

	return false
}

// checkCollision detects if the snake hit a wall or itself
func (m *Model) checkCollision() bool {
	head := m.snake[0]

	// Wall collision
	if head.X < 0 || head.X >= m.width ||
		head.Y < 0 || head.Y >= m.height {
		return true
	}

	// Self collision (check if head hits body)
	for i := 1; i < len(m.snake); i++ {
		if head == m.snake[i] {
			return true
		}
	}

	return false
}

// crash triggers the crash animation state
func (m *Model) crash() {
	m.state = StateCrashed
}

// gameOver transitions from crash to game over screen
func (m *Model) gameOver() {
	m.state = StateGameOver
	if m.score > m.highScore {
		m.highScore = m.score
	}
	m.recordStats() // Record stats and achievements
}
