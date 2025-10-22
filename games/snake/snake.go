package snake

// moveSnake updates the snake's position based on current direction
func (m *Model) moveSnake() {
	// Reset eating animation flag
	m.justAte = false

	// Update direction (prevent 180-degree turns)
	if !m.isOpposite(m.direction, m.nextDir) {
		m.direction = m.nextDir
	}

	// Calculate new head position
	head := m.snake[0]
	var newHead Point

	switch m.direction {
	case Up:
		newHead = Point{head.X, head.Y - 1}
	case Down:
		newHead = Point{head.X, head.Y + 1}
	case Left:
		newHead = Point{head.X - 1, head.Y}
	case Right:
		newHead = Point{head.X + 1, head.Y}
	}

	// Add new head to front of snake
	m.snake = append([]Point{newHead}, m.snake...)

	// Check if snake ate food
	if newHead == m.food {
		m.score++
		m.justAte = true // Show eating animation
		m.spawnFood()
		// Snake grows (don't remove tail)
	} else {
		// Remove tail (snake moves forward without growing)
		m.snake = m.snake[:len(m.snake)-1]
	}
}

// isOpposite checks if two directions are opposite to each other
func (m *Model) isOpposite(d1, d2 Direction) bool {
	return (d1 == Up && d2 == Down) ||
		(d1 == Down && d2 == Up) ||
		(d1 == Left && d2 == Right) ||
		(d1 == Right && d2 == Left)
}
