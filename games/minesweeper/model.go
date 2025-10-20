package minesweeper

import (
	"math/rand"
	"time"
)

// New creates a new minesweeper game
func New() Model {
	return Model{
		state:      StateMenu,
		difficulty: DifficultyEasy,
		firstClick: true,
	}
}

// InitGame initializes a new game with the selected difficulty
func (m *Model) InitGame() {
	config := difficultyConfigs[m.difficulty]
	m.width = config.width
	m.height = config.height
	m.mineCount = config.mineCount

	// Initialize grid
	m.grid = make([][]Cell, m.height)
	for i := range m.grid {
		m.grid[i] = make([]Cell, m.width)
	}

	m.state = StatePlaying
	m.startTime = time.Now()
	m.elapsedTime = 0
	m.flagsPlaced = 0
	m.cellsRevealed = 0
	m.firstClick = true
	m.cursorX = 0
	m.cursorY = 0
}

// PlaceMines randomly places mines on the grid, avoiding the first clicked cell
func (m *Model) PlaceMines(avoidX, avoidY int) {
	placed := 0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for placed < m.mineCount {
		x := rng.Intn(m.width)
		y := rng.Intn(m.height)

		// Don't place mine on first click or if already has mine
		if (x == avoidX && y == avoidY) || m.grid[y][x].IsMine {
			continue
		}

		m.grid[y][x].IsMine = true
		placed++
	}

	// Calculate adjacent mine counts for all cells
	m.calculateAdjacentMines()
}

// calculateAdjacentMines counts mines adjacent to each cell
func (m *Model) calculateAdjacentMines() {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if !m.grid[y][x].IsMine {
				m.grid[y][x].Adjacent = m.countAdjacentMines(x, y)
			}
		}
	}
}

// countAdjacentMines counts mines in the 8 surrounding cells
func (m *Model) countAdjacentMines(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx, ny := x+dx, y+dy
			if m.isValidCell(nx, ny) && m.grid[ny][nx].IsMine {
				count++
			}
		}
	}
	return count
}

// isValidCell checks if coordinates are within grid bounds
func (m *Model) isValidCell(x, y int) bool {
	return x >= 0 && x < m.width && y >= 0 && y < m.height
}

// RevealCell reveals a cell and potentially cascades to adjacent cells
func (m *Model) RevealCell(x, y int) {
	if !m.isValidCell(x, y) {
		return
	}

	cell := &m.grid[y][x]

	// Can't reveal flagged or already revealed cells
	if cell.IsFlagged || cell.IsRevealed {
		return
	}

	// First click - ensure it's safe
	if m.firstClick {
		m.PlaceMines(x, y)
		m.firstClick = false
	}

	// Reveal the cell
	cell.IsRevealed = true
	m.cellsRevealed++

	// Hit a mine - start explosion animation
	if cell.IsMine {
		m.state = StateExploding
		m.explosionCenterX = x
		m.explosionCenterY = y
		m.explosionRadius = 0

		// Calculate max steps based on grid size (max distance from any corner)
		maxDist := 0
		corners := [][2]int{{0, 0}, {m.width - 1, 0}, {0, m.height - 1}, {m.width - 1, m.height - 1}}
		for _, corner := range corners {
			dist := abs(corner[0]-x) + abs(corner[1]-y) // Manhattan distance
			if dist > maxDist {
				maxDist = dist
			}
		}
		m.explosionMaxSteps = maxDist + 3 // A few extra frames for effect
		return
	}

	// If no adjacent mines, cascade reveal
	if cell.Adjacent == 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				m.RevealCell(x+dx, y+dy)
			}
		}
	}

	// Check win condition
	m.checkWinCondition()
}

// ToggleFlag toggles flag on a cell
func (m *Model) ToggleFlag(x, y int) {
	if !m.isValidCell(x, y) {
		return
	}

	cell := &m.grid[y][x]

	// Can't flag revealed cells
	if cell.IsRevealed {
		return
	}

	cell.IsFlagged = !cell.IsFlagged

	if cell.IsFlagged {
		m.flagsPlaced++
	} else {
		m.flagsPlaced--
	}
}

// revealAllMines reveals all mines when game is lost
func (m *Model) revealAllMines() {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if m.grid[y][x].IsMine {
				m.grid[y][x].IsRevealed = true
			}
		}
	}
}

// checkWinCondition checks if player has won
func (m *Model) checkWinCondition() {
	totalCells := m.width * m.height
	safeCells := totalCells - m.mineCount

	if m.cellsRevealed == safeCells {
		m.state = StateWon
		m.elapsedTime = time.Since(m.startTime)

		// Update best time if this is better
		if m.bestTime == 0 || m.elapsedTime < m.bestTime {
			m.bestTime = m.elapsedTime
		}
	}
}

// progressExplosion advances the explosion animation by one frame
func (m *Model) progressExplosion() {
	m.explosionRadius++

	// Reveal mines within the current explosion radius
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			if m.grid[y][x].IsMine {
				// Calculate Manhattan distance from explosion center
				dist := abs(x-m.explosionCenterX) + abs(y-m.explosionCenterY)

				// Reveal if within current radius
				if dist <= m.explosionRadius {
					m.grid[y][x].IsRevealed = true
				}
			}
		}
	}

	// End animation when complete
	if m.explosionRadius >= m.explosionMaxSteps {
		m.state = StateLost
		m.revealAllMines() // Ensure all mines are revealed
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
