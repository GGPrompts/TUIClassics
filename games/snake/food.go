package snake

import "math/rand"

// spawnFood places food at a random empty location on the board
func (m *Model) spawnFood() {
	// Keep trying random positions until we find an empty spot
	for {
		food := Point{
			X: rand.Intn(m.width),
			Y: rand.Intn(m.height),
		}

		// Check if food spawned on snake
		onSnake := false
		for _, segment := range m.snake {
			if food == segment {
				onSnake = true
				break
			}
		}

		if !onSnake {
			m.food = food
			return
		}
	}
}
