package game2048

import "math/rand"

// spawnTile adds a new tile (2 or 4) to a random empty cell
func (m *Model) spawnTile() {
	// Find all empty cells
	var empty [][2]int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if m.grid[i][j].Value == 0 {
				empty = append(empty, [2]int{i, j})
			}
		}
	}

	if len(empty) == 0 {
		return // No empty cells
	}

	// Pick random empty cell
	pos := empty[rand.Intn(len(empty))]

	// 90% chance of 2, 10% chance of 4
	value := 2
	if rand.Float64() < 0.1 {
		value = 4
	}

	m.grid[pos[0]][pos[1]] = Tile{Value: value}
}

// hasEmptyCell checks if there are any empty cells on the grid
func (m *Model) hasEmptyCell() bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if m.grid[i][j].Value == 0 {
				return true
			}
		}
	}
	return false
}

// hasValue checks if any tile has the target value
func (m *Model) hasValue(target int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if m.grid[i][j].Value == target {
				return true
			}
		}
	}
	return false
}
