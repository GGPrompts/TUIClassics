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

	b.WriteString(titleStyle.Render("💣 MINESWEEPER 💣"))
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

	// Title
	title := titleStyle.Render("💣 MINESWEEPER 💣")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Stats line
	config := difficultyConfigs[m.difficulty]
	stats := fmt.Sprintf(
		"💣 Mines: %d/%d  |  ⏱️  %s  |  🚩 Flags: %d",
		m.mineCount,
		config.mineCount,
		formatDuration(m.elapsedTime),
		m.flagsPlaced,
	)
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	// Grid
	grid := m.renderGrid()
	b.WriteString(gridBorderStyle.Render(grid))
	b.WriteString("\n\n")

	// Help text
	help := "Left-click: Reveal  |  Right-click: Flag  |  [N] New Game  |  [Q] Quit"
	b.WriteString(helpStyle.Render(help))

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
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
	var symbol string
	var style lipgloss.Style

	if cell.IsFlagged && !cell.IsRevealed {
		symbol = "🚩"
		style = flagCellStyle
	} else if !cell.IsRevealed {
		symbol = "□"
		style = unknownCellStyle
	} else if cell.IsMine {
		symbol = "💣"
		style = mineCellStyle
	} else if cell.Adjacent == 0 {
		symbol = " "
		style = revealedCellStyle
	} else {
		symbol = fmt.Sprintf("%d", cell.Adjacent)
		style = getNumberStyle(cell.Adjacent)
	}

	// Apply cursor highlight
	if isCursor && !cell.IsRevealed {
		style = cursorCellStyle
	}

	return style.Render(symbol)
}

// renderWin renders the win screen
func (m Model) renderWin() string {
	var b strings.Builder

	b.WriteString(winStyle.Render("🎉 YOU WIN! 🎉"))
	b.WriteString("\n\n")

	grid := m.renderGrid()
	b.WriteString(gridBorderStyle.Render(grid))
	b.WriteString("\n\n")

	stats := fmt.Sprintf(
		"Time: %s  |  Best: %s",
		formatDuration(m.elapsedTime),
		formatDuration(m.bestTime),
	)
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	help := "[R] Restart  |  [N] New Game  |  [Q] Quit"
	b.WriteString(helpStyle.Render(help))

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
}

// renderLoss renders the loss screen
func (m Model) renderLoss() string {
	var b strings.Builder

	b.WriteString(loseStyle.Render("💥 GAME OVER 💥"))
	b.WriteString("\n\n")

	grid := m.renderGrid()
	b.WriteString(gridBorderStyle.Render(grid))
	b.WriteString("\n\n")

	help := "[R] Restart  |  [N] New Game  |  [Q] Quit"
	b.WriteString(helpStyle.Render(help))

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
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
