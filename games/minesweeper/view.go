package minesweeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the game
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.renderMenu()
	case StatePlaying:
		return m.renderGame()
	case StateExploding:
		return m.renderExplosion()
	case StateWon:
		return m.renderWin()
	case StateLost:
		return m.renderLoss()
	default:
		return ""
	}
}

// renderMenu renders the difficulty selection menu
func (m Model) renderMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("* * * MINESWEEPER * * *"))
	b.WriteString("\n\n")

	menu := `
Select Difficulty:

[1] Easy     - 8x8   (10 mines)
[2] Medium   - 16x16 (40 mines)
[3] Hard     - 30x16 (99 mines)

[Q] Quit
`

	b.WriteString(menuStyle.Render(menu))

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
}

// renderGame renders the active game
func (m Model) renderGame() string {
	var b strings.Builder

	// Calculate vertical centering
	totalLines := 3 + 2 + 1 + 2 + m.height + 4 + 2 // title + gap + smiley + gap + grid + gap + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add top padding
	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title
	title := titleStyle.Width(m.termWidth).Render("* * * MINESWEEPER * * *")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Smiley face button (clickable to restart)
	smiley := m.renderSmiley()
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(smiley))
	b.WriteString("\n\n")

	// Stats line
	config := difficultyConfigs[m.difficulty]
	stats := fmt.Sprintf(
		"Mines: %d/%d  |  Time: %s  |  Flags: %d",
		m.mineCount,
		config.mineCount,
		formatDuration(m.elapsedTime),
		m.flagsPlaced,
	)
	b.WriteString(statsStyle.Width(m.termWidth).Align(lipgloss.Center).Render(stats))
	b.WriteString("\n\n")

	// Grid (with calculated horizontal centering)
	grid := m.renderGrid()
	gridWidth := m.width * 2 // Each cell is roughly 2 chars
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Render grid with padding
	gridLines := strings.Split(grid, "\n")
	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// Help text
	help := "Left-click: Reveal  |  Right-click/[F]: Flag  |  Arrow keys: Move  |  [N] New  |  [Q] Quit"
	b.WriteString(helpStyle.Width(m.termWidth).Render(help))

	return b.String()
}

// renderExplosion renders the explosion animation
func (m Model) renderExplosion() string {
	var b strings.Builder

	// Calculate vertical centering
	totalLines := 3 + 2 + 1 + 2 + m.height + 4 + 2 // title + gap + smiley + gap + grid + gap + help
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	// Title with explosion effect
	title := titleStyle.Width(m.termWidth).Render("* * * BOOM! * * *")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Smiley face (dead during explosion)
	smiley := m.renderSmiley()
	b.WriteString(lipgloss.NewStyle().Width(m.termWidth).Align(lipgloss.Center).Render(smiley))
	b.WriteString("\n\n")

	// Stats line
	config := difficultyConfigs[m.difficulty]
	stats := fmt.Sprintf(
		"Mines: %d/%d  |  Time: %s  |  Flags: %d",
		m.mineCount,
		config.mineCount,
		formatDuration(m.elapsedTime),
		m.flagsPlaced,
	)
	b.WriteString(statsStyle.Width(m.termWidth).Align(lipgloss.Center).Render(stats))
	b.WriteString("\n\n")

	// Grid with explosion effect
	grid := m.renderExplosionGrid()
	gridWidth := m.width * 2
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	gridLines := strings.Split(grid, "\n")
	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// Help text
	help := fmt.Sprintf("Explosion expanding... (radius: %d)", m.explosionRadius)
	b.WriteString(helpStyle.Width(m.termWidth).Align(lipgloss.Center).Render(help))

	return b.String()
}

// renderExplosionGrid renders the grid during explosion animation
func (m Model) renderExplosionGrid() string {
	var b strings.Builder

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			cell := m.grid[y][x]
			dist := abs(x-m.explosionCenterX) + abs(y-m.explosionCenterY)

			var symbol string
			var style lipgloss.Style

			// Show explosion wave effect
			if dist == m.explosionRadius && cell.IsMine {
				// Mines at current explosion radius show as exploding
				symbol = "+"
				style = mineCellStyle.Copy().Blink(true)
			} else if dist < m.explosionRadius && cell.IsMine && cell.IsRevealed {
				// Already revealed mines
				symbol = "*"
				style = mineCellStyle
			} else {
				// Regular cell rendering
				symbol, style = m.getCellSymbolAndStyle(cell, false)
			}

			b.WriteString(style.Render(symbol))
			b.WriteString(" ")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// getCellSymbolAndStyle returns the symbol and style for a cell
func (m Model) getCellSymbolAndStyle(cell Cell, isCursor bool) (string, lipgloss.Style) {
	var symbol string
	var style lipgloss.Style

	if cell.IsFlagged && !cell.IsRevealed {
		symbol = "P"
		style = flagCellStyle
	} else if !cell.IsRevealed {
		symbol = "■"
		style = unknownCellStyle
	} else if cell.IsMine {
		symbol = "*"
		style = mineCellStyle
	} else if cell.Adjacent == 0 {
		symbol = "·"
		style = revealedCellStyle
	} else {
		symbol = fmt.Sprintf("%d", cell.Adjacent)
		style = getNumberStyle(cell.Adjacent)
	}

	if isCursor && !cell.IsRevealed {
		style = cursorCellStyle
	}

	return symbol, style
}

// renderGrid renders the minesweeper grid
func (m Model) renderGrid() string {
	var b strings.Builder

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			cell := m.grid[y][x]
			isCursor := (x == m.cursorX && y == m.cursorY)

			cellStr := m.renderCell(cell, isCursor)
			b.WriteString(cellStr)
			b.WriteString(" ") // Space between cells
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderCell renders a single cell
func (m Model) renderCell(cell Cell, isCursor bool) string {
	symbol, style := m.getCellSymbolAndStyle(cell, isCursor)
	return style.Render(symbol)
}

// renderSmiley renders the clickable smiley face button
func (m Model) renderSmiley() string {
	var face string

	switch m.smileyState {
	case SmileyHappy:
		face = "🙂"  // Slightly smiling face
	case SmileySurprised:
		face = "😮"  // Face with open mouth
	case SmileyDead:
		face = "😵"  // Dizzy/dead face
	case SmileyCool:
		face = "😎"  // Smiling face with sunglasses (winner!)
	default:
		face = "🙂"
	}

	return face
}

// renderWin renders the win screen
func (m Model) renderWin() string {
	var b strings.Builder

	// Calculate vertical centering
	totalLines := 1 + 2 + m.height + 2 + 1 + 2 + 1
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	b.WriteString(winStyle.Width(m.termWidth).Render("🎉 YOU WIN! 🎉"))
	b.WriteString("\n\n")

	// Grid with padding
	grid := m.renderGrid()
	gridWidth := m.width * 2
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	gridLines := strings.Split(grid, "\n")
	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	stats := fmt.Sprintf(
		"Time: %s  |  Best: %s",
		formatDuration(m.elapsedTime),
		formatDuration(m.bestTime),
	)
	b.WriteString(statsStyle.Width(m.termWidth).Align(lipgloss.Center).Render(stats))
	b.WriteString("\n\n")

	help := "[R] Restart  |  [N] New Game  |  [Q] Quit"
	b.WriteString(helpStyle.Width(m.termWidth).Render(help))

	return b.String()
}

// renderLoss renders the loss screen
func (m Model) renderLoss() string {
	var b strings.Builder

	// Calculate vertical centering
	totalLines := 1 + 2 + m.height + 2 + 1
	topPadding := (m.termHeight - totalLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	for i := 0; i < topPadding; i++ {
		b.WriteString("\n")
	}

	b.WriteString(loseStyle.Width(m.termWidth).Render("💥 GAME OVER 💥"))
	b.WriteString("\n\n")

	// Grid with padding
	grid := m.renderGrid()
	gridWidth := m.width * 2
	leftPadding := (m.termWidth - gridWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	gridLines := strings.Split(grid, "\n")
	for _, line := range gridLines {
		if line != "" {
			b.WriteString(strings.Repeat(" ", leftPadding))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	help := "[R] Restart  |  [N] New Game  |  [Q] Quit"
	b.WriteString(helpStyle.Width(m.termWidth).Render(help))

	return b.String()
}

// formatDuration formats a duration as MM:SS
func formatDuration(d interface{}) string {
	var seconds int

	switch v := d.(type) {
	case int:
		seconds = v
	default:
		// Assume it's time.Duration
		if dur, ok := d.(time.Duration); ok {
			seconds = int(dur.Seconds())
		}
	}

	minutes := seconds / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
