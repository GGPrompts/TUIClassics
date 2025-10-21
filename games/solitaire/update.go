package solitaire

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

	case TickMsg:
		if m.state == StatePlaying {
			m.elapsedTime = time.Since(m.startTime)
		}
		return m, tickCmd()

	case AnimationTickMsg:
		if m.animating {
			m.ProgressWaterfallAnimation()
			return m, animationTick()
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "n", "N":
		// Start new game
		m.NewGame()
		return m, nil
	}

	// State-specific keys
	switch m.state {
	case StateMenu:
		return m.handleMenuKeys(msg)
	case StatePlaying:
		return m.handleGameKeys(msg)
	case StateWon:
		return m.handleWinKeys(msg)
	}

	return m, nil
}

// handleMenuKeys handles keys in menu state
func (m Model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "N", "enter", " ":
		m.NewGame()
		return m, nil
	}
	return m, nil
}

// handleWinKeys handles keys in win state
func (m Model) handleWinKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "N", "enter", " ":
		m.NewGame()
		return m, nil
	}
	return m, nil
}

// handleGameKeys handles keys during gameplay
func (m Model) handleGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.moveCursorUp()
		return m, nil

	case "down", "j":
		m.moveCursorDown()
		return m, nil

	case "left", "h":
		m.moveCursorLeft()
		return m, nil

	case "right", "l":
		m.moveCursorRight()
		return m, nil

	case "enter", " ":
		return m.handleSelect()

	case "d", "D":
		// Draw from stock
		m.DrawFromStock()
		return m, nil

	case "1", "2", "3", "4":
		// Quick move to foundation
		return m.handleQuickFoundationMove(msg.String())

	case "w", "W":
		// Secret: Trigger waterfall animation (for testing/demo)
		m.state = StateWon
		m.elapsedTime = time.Since(m.startTime)
		m.StartWaterfallAnimation()
		return m, animationTick()
	}

	return m, nil
}

// moveCursorLeft moves the cursor left
func (m *Model) moveCursorLeft() {
	switch m.cursor.PileType {
	case TableauPile:
		if m.cursor.PileIndex > 0 {
			m.cursor.PileIndex--
		} else {
			// Wrap to foundation
			m.cursor.PileType = FoundationPile
			m.cursor.PileIndex = 3
		}
		// Reset card index when moving between piles
		m.resetCursorCardIndex()

	case FoundationPile:
		if m.cursor.PileIndex > 0 {
			m.cursor.PileIndex--
		} else {
			// Wrap to waste
			m.cursor.PileType = WastePile
			m.cursor.PileIndex = 0
		}

	case WastePile:
		m.cursor.PileType = StockPile
		m.cursor.PileIndex = 0

	case StockPile:
		// Wrap to tableau
		m.cursor.PileType = TableauPile
		m.cursor.PileIndex = 6
		m.resetCursorCardIndex()
	}
}

// moveCursorRight moves the cursor right
func (m *Model) moveCursorRight() {
	switch m.cursor.PileType {
	case StockPile:
		m.cursor.PileType = WastePile
		m.cursor.PileIndex = 0

	case WastePile:
		m.cursor.PileType = FoundationPile
		m.cursor.PileIndex = 0

	case FoundationPile:
		if m.cursor.PileIndex < 3 {
			m.cursor.PileIndex++
		} else {
			// Wrap to tableau
			m.cursor.PileType = TableauPile
			m.cursor.PileIndex = 0
			m.resetCursorCardIndex()
		}

	case TableauPile:
		if m.cursor.PileIndex < 6 {
			m.cursor.PileIndex++
		} else {
			// Wrap to stock
			m.cursor.PileType = StockPile
			m.cursor.PileIndex = 0
		}
		// Reset card index when moving between piles
		m.resetCursorCardIndex()
	}
}

// moveCursorUp moves the cursor up within a tableau stack
func (m *Model) moveCursorUp() {
	// Only works on tableau piles
	if m.cursor.PileType != TableauPile {
		return
	}

	pile := m.tableau[m.cursor.PileIndex]
	if len(pile.Cards) == 0 {
		return
	}

	// Move up to previous card (earlier in the stack)
	if m.cursorCardIndex > 0 {
		// Check if the previous card is face up
		if pile.Cards[m.cursorCardIndex-1].FaceUp {
			m.cursorCardIndex--
		}
	}
}

// moveCursorDown moves the cursor down within a tableau stack
func (m *Model) moveCursorDown() {
	// Only works on tableau piles
	if m.cursor.PileType != TableauPile {
		return
	}

	pile := m.tableau[m.cursor.PileIndex]
	if len(pile.Cards) == 0 {
		return
	}

	// Move down to next card (later in the stack)
	maxIndex := len(pile.Cards) - 1
	if m.cursorCardIndex < maxIndex {
		m.cursorCardIndex++
	}
}

// resetCursorCardIndex resets the card index to the top of the pile
func (m *Model) resetCursorCardIndex() {
	if m.cursor.PileType == TableauPile {
		pile := m.tableau[m.cursor.PileIndex]
		if len(pile.Cards) > 0 {
			m.cursorCardIndex = len(pile.Cards) - 1
		} else {
			m.cursorCardIndex = 0
		}
	}
}

// handleSelect handles the Enter/Space key for selecting and moving cards
func (m Model) handleSelect() (tea.Model, tea.Cmd) {
	// If we're on stock, draw a card
	if m.cursor.PileType == StockPile {
		m.DrawFromStock()
		return m, nil
	}

	// If no card is selected, select the current card
	if m.selectedPile == nil {
		return m.selectCard()
	}

	// Try to move the selected card to the current cursor position
	success := m.moveSelectedCard()
	if success {
		m.selectedPile = nil
		m.selectedCard = nil

		// Check win condition
		if m.CheckWin() {
			m.state = StateWon
			m.elapsedTime = time.Since(m.startTime)
			m.StartWaterfallAnimation()
			return m, animationTick()
		}
	} else {
		// If move failed, deselect
		m.selectedPile = nil
		m.selectedCard = nil
	}

	return m, nil
}

// selectCard selects the card at the current cursor position
func (m Model) selectCard() (tea.Model, tea.Cmd) {
	switch m.cursor.PileType {
	case TableauPile:
		pile := m.tableau[m.cursor.PileIndex]
		if len(pile.Cards) > 0 && m.cursorCardIndex < len(pile.Cards) {
			card := pile.Cards[m.cursorCardIndex]
			if card.FaceUp {
				m.selectedPile = &CursorLocation{
					PileType:  TableauPile,
					PileIndex: m.cursor.PileIndex,
				}
				m.selectedIndex = m.cursorCardIndex
				m.selectedCard = &card
			}
		}

	case WastePile:
		if len(m.waste.Cards) > 0 {
			topCard := m.waste.Cards[len(m.waste.Cards)-1]
			m.selectedPile = &CursorLocation{
				PileType:  WastePile,
				PileIndex: 0,
			}
			m.selectedCard = &topCard
		}

	case FoundationPile:
		pile := m.foundation[m.cursor.PileIndex]
		if len(pile.Cards) > 0 {
			topCard := pile.Cards[len(pile.Cards)-1]
			m.selectedPile = &CursorLocation{
				PileType:  FoundationPile,
				PileIndex: m.cursor.PileIndex,
			}
			m.selectedCard = &topCard
		}
	}

	return m, nil
}

// moveSelectedCard attempts to move the selected card to the cursor position
func (m *Model) moveSelectedCard() bool {
	if m.selectedPile == nil {
		return false
	}

	// Can't move to same location
	if m.cursor.PileType == m.selectedPile.PileType && m.cursor.PileIndex == m.selectedPile.PileIndex {
		return false
	}

	// Get the cards to move
	var cards []Card

	if m.selectedPile.PileType == TableauPile {
		// Get movable sequence from tableau
		sequence := m.GetMovableSequence(m.selectedPile.PileIndex, m.selectedIndex)
		if sequence == nil {
			return false
		}
		cards = sequence
	} else {
		// Single card from waste or foundation
		cards = []Card{*m.selectedCard}
	}

	// Attempt the move
	return m.MoveCards(*m.selectedPile, m.cursor, cards)
}

// handleQuickFoundationMove handles number keys for quick foundation moves
func (m Model) handleQuickFoundationMove(key string) (tea.Model, tea.Cmd) {
	foundationIndex := -1
	switch key {
	case "1":
		foundationIndex = 0
	case "2":
		foundationIndex = 1
	case "3":
		foundationIndex = 2
	case "4":
		foundationIndex = 3
	}

	if foundationIndex == -1 {
		return m, nil
	}

	// Try to move from current cursor to foundation
	var card *Card
	var from CursorLocation

	switch m.cursor.PileType {
	case TableauPile:
		pile := m.tableau[m.cursor.PileIndex]
		if len(pile.Cards) > 0 && pile.Cards[len(pile.Cards)-1].FaceUp {
			card = &pile.Cards[len(pile.Cards)-1]
			from = CursorLocation{PileType: TableauPile, PileIndex: m.cursor.PileIndex}
		}

	case WastePile:
		if len(m.waste.Cards) > 0 {
			card = &m.waste.Cards[len(m.waste.Cards)-1]
			from = CursorLocation{PileType: WastePile, PileIndex: 0}
		}
	}

	if card != nil {
		to := CursorLocation{PileType: FoundationPile, PileIndex: foundationIndex}
		m.MoveCards(from, to, []Card{*card})

		// Check win condition
		if m.CheckWin() {
			m.state = StateWon
			m.elapsedTime = time.Since(m.startTime)
			m.StartWaterfallAnimation()
			return m, animationTick()
		}
	}

	return m, nil
}
