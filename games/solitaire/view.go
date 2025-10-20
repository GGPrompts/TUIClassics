package solitaire

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.viewMenu()
	case StatePlaying:
		if m.animating {
			return m.viewWaterfall()
		}
		return m.viewGame()
	case StateWon:
		return m.viewWin()
	default:
		return m.viewGame()
	}
}

// viewMenu renders the main menu
func (m Model) viewMenu() string {
	title := titleStyle.Render("♠ ♥ SOLITAIRE ♦ ♣")

	menu := menuStyle.Render(
		"Welcome to Klondike Solitaire!\n\n" +
			"Press [N] to start a New Game\n" +
			"Press [Q] to Quit",
	)

	content := lipgloss.JoinVertical(lipgloss.Center, title, menu)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// viewGame renders the main game view
func (m Model) viewGame() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("♠ ♥ SOLITAIRE ♦ ♣")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Stats
	elapsed := m.elapsedTime
	if m.state == StatePlaying {
		elapsed = time.Since(m.startTime)
	}
	stats := fmt.Sprintf("Score: %d  |  Moves: %d  |  Time: %s",
		m.score, m.moves, formatDuration(elapsed))
	b.WriteString(statsStyle.Render(stats))
	b.WriteString("\n\n")

	// Top row: Stock, Waste, and Foundation piles
	topRow := m.renderTopRow()
	b.WriteString(topRow)
	b.WriteString("\n\n")

	// Tableau piles
	tableau := m.renderTableau()
	b.WriteString(tableau)
	b.WriteString("\n")

	// Help text
	help := helpStyle.Render("Arrows: Navigate | Enter: Select/Move | Space: Draw | N: New Game | Q: Quit")
	b.WriteString(help)

	return b.String()
}

// renderTopRow renders the stock, waste, and foundation piles
func (m Model) renderTopRow() string {
	var cards []string

	// Stock pile
	stockCard := m.renderStockPile()
	cards = append(cards, stockCard)

	// Waste pile
	wasteCard := m.renderWastePile()
	cards = append(cards, wasteCard)

	// Space between waste and foundation
	cards = append(cards, strings.Repeat(" ", 3))

	// Foundation piles (4 piles)
	for i := 0; i < 4; i++ {
		foundationCard := m.renderFoundationPile(i)
		cards = append(cards, foundationCard)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...)
}

// renderStockPile renders the stock pile
func (m Model) renderStockPile() string {
	if len(m.stock.Cards) == 0 {
		// Empty stock - show recycling indicator
		content := lipgloss.JoinVertical(lipgloss.Center,
			"",
			centerText("↻", 5),
			"",
		)
		empty := emptyPileStyle.Render(content)
		return empty
	}

	// Show card back
	isCursor := m.cursor.PileType == StockPile
	return m.renderCardBack(isCursor)
}

// renderWastePile renders the waste pile
func (m Model) renderWastePile() string {
	if len(m.waste.Cards) == 0 {
		content := lipgloss.JoinVertical(lipgloss.Center,
			"",
			centerText("", 5),
			"",
		)
		empty := emptyPileStyle.Render(content)
		return empty
	}

	// Show top card
	topCard := m.waste.Cards[len(m.waste.Cards)-1]
	isCursor := m.cursor.PileType == WastePile
	isSelected := m.selectedPile != nil && m.selectedPile.PileType == WastePile
	return m.renderCard(topCard, isCursor, isSelected)
}

// renderFoundationPile renders a foundation pile
func (m Model) renderFoundationPile(index int) string {
	pile := m.foundation[index]

	if len(pile.Cards) == 0 {
		// Empty foundation - show suit symbol
		var suit string
		switch index {
		case 0:
			suit = Spades.SuitSymbol()
		case 1:
			suit = Hearts.SuitSymbol()
		case 2:
			suit = Diamonds.SuitSymbol()
		case 3:
			suit = Clubs.SuitSymbol()
		}
		content := lipgloss.JoinVertical(lipgloss.Center,
			"",
			centerText(suit, 5),
			"",
		)
		empty := emptyPileStyle.Render(content)
		return empty
	}

	// Show top card
	topCard := pile.Cards[len(pile.Cards)-1]
	isCursor := m.cursor.PileType == FoundationPile && m.cursor.PileIndex == index
	isSelected := m.selectedPile != nil && m.selectedPile.PileType == FoundationPile && m.selectedPile.PileIndex == index
	return m.renderCard(topCard, isCursor, isSelected)
}

// renderTableau renders all tableau piles with stacking
func (m Model) renderTableau() string {
	var columns []string

	for col := 0; col < 7; col++ {
		pile := m.tableau[col]

		if len(pile.Cards) == 0 {
			// Empty pile - show King placeholder
			content := lipgloss.JoinVertical(lipgloss.Center,
				"",
				centerText("K", 5),
				"",
			)
			empty := emptyPileStyle.Render(content)
			columns = append(columns, empty)
		} else {
			// Build this column from top to bottom
			var columnContent strings.Builder

			for cardIdx := 0; cardIdx < len(pile.Cards); cardIdx++ {
				card := pile.Cards[cardIdx]
				isLast := cardIdx == len(pile.Cards)-1
				isCursor := m.cursor.PileType == TableauPile && m.cursor.PileIndex == col && isLast
				isSelected := m.selectedPile != nil && m.selectedPile.PileType == TableauPile &&
					m.selectedPile.PileIndex == col && m.selectedIndex == cardIdx

				if isLast {
					// Last card in pile - show full card
					if card.FaceUp {
						columnContent.WriteString(m.renderCard(card, isCursor, isSelected))
					} else {
						columnContent.WriteString(m.renderCardBack(false))
					}
				} else {
					// Not last card - show only top lines for stacking effect
					if card.FaceUp {
						columnContent.WriteString(m.renderCardTopLine(card, isSelected))
						columnContent.WriteString("\n")
					} else {
						columnContent.WriteString(m.renderCardBackTopLine())
						columnContent.WriteString("\n")
					}
				}
			}

			columns = append(columns, columnContent.String())
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

// renderCard renders a face-up card
func (m Model) renderCard(card Card, isCursor, isSelected bool) string {
	suit := card.Suit.SuitSymbol()
	rank := card.Rank.RankSymbol()
	color := card.Suit.Color()

	// Match the placeholder format exactly: 5 chars wide, centered content
	line1 := fmt.Sprintf("%s %s", rank, suit)
	line3 := fmt.Sprintf("%s %s", suit, rank)

	// Use the same structure as empty placeholders: Center alignment, 3 lines
	content := lipgloss.JoinVertical(lipgloss.Center,
		centerText(line1, 5),
		"",
		centerText(line3, 5),
	)

	style := cardStyle.Copy().Foreground(color)
	if isSelected {
		style = selectedCardStyle.Copy().Foreground(color)
	} else if isCursor {
		style = cursorCardStyle.Copy().Foreground(color)
	}

	return style.Render(content)
}

// renderCardBack renders a face-down card
func (m Model) renderCardBack(isCursor bool) string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		centerText("░░░", 5),
		"",
	)

	style := cardBackStyle
	if isCursor {
		style = cardBackStyle.Copy().BorderForeground(lipgloss.Color("#00FF00"))
	}
	return style.Render(content)
}

// renderCardTopLine renders just the top 2 lines of a face-up card (for stacking)
// This creates visual overlap so you see the rounded top of each card
func (m Model) renderCardTopLine(card Card, isSelected bool) string {
	// Render a full card using the same logic as renderCard, then extract top 2 lines
	suit := card.Suit.SuitSymbol()
	rank := card.Rank.RankSymbol()
	color := card.Suit.Color()

	// Build content exactly like full cards
	line1 := fmt.Sprintf("%s %s", rank, suit)
	line3 := fmt.Sprintf("%s %s", suit, rank)

	content := lipgloss.JoinVertical(lipgloss.Center,
		centerText(line1, 5),
		"",
		centerText(line3, 5),
	)

	// Use the appropriate style to get the same dimensions
	style := cardStyle.Copy().Foreground(color)
	if isSelected {
		style = selectedCardStyle.Copy().Foreground(color)
	}

	// Render the full card
	fullCard := style.Render(content)

	// Extract just the top 2 lines
	lines := strings.Split(fullCard, "\n")
	if len(lines) >= 2 {
		return lines[0] + "\n" + lines[1]
	}
	return fullCard
}

// renderCardBackTopLine renders just the top 2 lines of a face-down card (for stacking)
func (m Model) renderCardBackTopLine() string {
	// Render a full card back using the same logic, then extract top 2 lines
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		centerText("░░░", 5),
		"",
	)

	// Use the card back style to get the same dimensions
	fullCard := cardBackStyle.Render(content)

	// Extract just the top 2 lines
	lines := strings.Split(fullCard, "\n")
	if len(lines) >= 2 {
		return lines[0] + "\n" + lines[1]
	}
	return fullCard
}

// viewWaterfall renders the waterfall animation
func (m Model) viewWaterfall() string {
	// Create a grid to place cards
	grid := make([][]rune, m.termHeight)
	for i := range grid {
		grid[i] = make([]rune, m.termWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw background text
	centerY := m.termHeight / 2
	centerX := m.termWidth / 2
	winText := "♠ ♥ YOU WON! ♦ ♣"
	startX := centerX - len(winText)/2
	if startX >= 0 && centerY >= 0 && centerY < m.termHeight {
		for i, ch := range winText {
			if startX+i < m.termWidth {
				grid[centerY][startX+i] = ch
			}
		}
	}

	// Draw falling cards
	for _, wc := range m.waterfallCards {
		x := int(wc.x)
		y := int(wc.y)

		if y >= 0 && y < m.termHeight && x >= 0 && x < m.termWidth-5 {
			// Simple card representation for animation
			rank := wc.card.Rank.RankSymbol()
			suit := wc.card.Suit.SuitSymbol()
			cardStr := rank + suit

			for i, ch := range cardStr {
				if x+i < m.termWidth {
					grid[y][x+i] = ch
				}
			}
		}
	}

	// Convert grid to string
	var b strings.Builder
	for _, row := range grid {
		b.WriteString(string(row))
		b.WriteRune('\n')
	}

	return b.String()
}

// viewWin renders the win screen
func (m Model) viewWin() string {
	elapsed := m.elapsedTime
	if m.state == StatePlaying {
		elapsed = time.Since(m.startTime)
	}

	content := fmt.Sprintf(
		"♠ ♥ YOU WON! ♦ ♣\n\n"+
			"Score: %d\n"+
			"Moves: %d\n"+
			"Time: %s\n\n"+
			"Press [N] for New Game or [Q] to Quit",
		m.score,
		m.moves,
		formatDuration(elapsed),
	)

	win := winStyle.Render(content)

	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		win,
	)
}

// Helper functions

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
