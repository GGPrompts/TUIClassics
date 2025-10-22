package game2048

// slideAndMergeRow slides and merges tiles in a single row (left direction)
// Returns true if the row changed
func (m *Model) slideAndMergeRow(row int) bool {
	// Extract non-zero values from the row
	values := []int{}
	for j := 0; j < 4; j++ {
		if m.grid[row][j].Value != 0 {
			values = append(values, m.grid[row][j].Value)
		}
	}

	// Merge adjacent equal values
	merged := []int{}
	i := 0
	for i < len(values) {
		if i+1 < len(values) && values[i] == values[i+1] {
			// Merge the two tiles
			newValue := values[i] * 2
			merged = append(merged, newValue)
			m.score += newValue
			i += 2 // Skip both merged tiles
		} else {
			// No merge, just keep the value
			merged = append(merged, values[i])
			i++
		}
	}

	// Check if row changed
	changed := false
	for j := 0; j < 4; j++ {
		newValue := 0
		if j < len(merged) {
			newValue = merged[j]
		}

		if m.grid[row][j].Value != newValue {
			changed = true
		}

		m.grid[row][j] = Tile{Value: newValue}
	}

	return changed
}

// transpose swaps rows and columns (for vertical movement)
func (m *Model) transpose() {
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			m.grid[i][j], m.grid[j][i] = m.grid[j][i], m.grid[i][j]
		}
	}
}

// reverseRows reverses each row (for right movement)
func (m *Model) reverseRows() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			m.grid[i][j], m.grid[i][3-j] = m.grid[i][3-j], m.grid[i][j]
		}
	}
}
