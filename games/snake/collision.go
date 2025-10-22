package snake

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
}
