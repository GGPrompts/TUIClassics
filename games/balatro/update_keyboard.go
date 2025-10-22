package balatro

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// update_keyboard.go - Keyboard Event Handling
// Purpose: All keyboard input processing
// When to extend: Add new keyboard shortcuts or key bindings here

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keybindings (work in all modes)
	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Help):
		return m.showHelp()

	case key.Matches(msg, keys.Refresh):
		return m.refresh()
	}

	// Mode-specific keybindings
	switch m.focusedComponent {
	case "landing":
		return m.handleLandingKeys(msg)

	case "game":
		return m.handleGameKeys(msg)

	case "main":
		return m.handleMainKeys(msg)

	// Add handlers for other components/modes
	// case "dialog":
	//     return m.handleDialogKeys(msg)
	//
	// case "menu":
	//     return m.handleMenuKeys(msg)
	}

	return m, nil
}

// handleLandingKeys handles keys in landing page
func (m Model) handleLandingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.landingPage == nil {
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		m.landingPage.SelectPrev()
		return m, nil

	case "down", "j":
		m.landingPage.SelectNext()
		return m, nil

	case "enter":
		// Handle menu selection
		selected := m.landingPage.GetSelectedItem()
		switch selected {
		case "NEW GAME":
			m.newGame()
			m.focusedComponent = "game"
		case "CONTINUE":
			m.statusMsg = "Continue not yet implemented"
		case "COLLECTION":
			m.state = StateCollection
			m.focusedComponent = "main"
			m.statusMsg = "Opening collection..."
		case "OPTIONS":
			m.state = StateOptions
			m.focusedComponent = "main"
			m.statusMsg = "Opening options..."
		case "QUIT":
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

// handleGameKeys handles keys in poker game
func (m Model) handleGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Phase 2: Handle different phases
	switch m.gamePhase {
	case PhaseBlindComplete:
		return m.handleBlindCompleteKeys(msg)
	case PhaseShop:
		return m.handleShopKeys(msg)
	case PhaseGameOver:
		return m.handleGameOverKeys(msg)
	case PhaseSelectCards:
		// Normal gameplay
		return m.handleSelectCardsKeys(msg)
	default:
		return m.handleSelectCardsKeys(msg)
	}
}

// handleSelectCardsKeys handles normal card selection and play
func (m Model) handleSelectCardsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Number keys 1-8 to select cards
	case "1":
		m.selectedCardIndex = 0
		return m, nil
	case "2":
		m.selectedCardIndex = 1
		return m, nil
	case "3":
		m.selectedCardIndex = 2
		return m, nil
	case "4":
		m.selectedCardIndex = 3
		return m, nil
	case "5":
		m.selectedCardIndex = 4
		return m, nil
	case "6":
		m.selectedCardIndex = 5
		return m, nil
	case "7":
		m.selectedCardIndex = 6
		return m, nil
	case "8":
		m.selectedCardIndex = 7
		return m, nil

	// Arrow keys for card selection
	case "left", "h":
		if m.selectedCardIndex > 0 {
			m.selectedCardIndex--
		}
		return m, nil

	case "right", "l":
		if m.selectedCardIndex < len(m.hand)-1 {
			m.selectedCardIndex++
		}
		return m, nil

	// Space to toggle card for play
	case " ":
		m.toggleCardForPlay(m.selectedCardIndex)
		return m, nil

	// Enter to play hand
	case "enter":
		m.playHand()
		return m, nil

	// D to discard
	case "d", "D":
		m.discardCards()
		return m, nil
	}

	return m, nil
}

// handleBlindCompleteKeys handles victory screen
func (m Model) handleBlindCompleteKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Continue to shop
		m.gamePhase = PhaseShop
		m.statusMsg = "Welcome to the shop!"
		return m, nil

	case "q":
		// Quit to menu
		m.state = StateLanding
		m.focusedComponent = "landing"
		return m, nil
	}

	return m, nil
}

// handleShopKeys handles shop screen
func (m Model) handleShopKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Advance to next blind
		m.roundState.AdvanceBlind()

		// Sync legacy fields
		m.money = m.roundState.Money
		m.currentScore = m.roundState.CurrentScore
		m.handsRemaining = m.roundState.HandsRemaining
		m.discardsRemaining = m.roundState.DiscardsRemaining
		m.targetScore = m.roundState.TargetScore

		// Reset for new blind
		m.gamePhase = PhaseSelectCards
		m.statusMsg = fmt.Sprintf("Ante %d - %s: Reach %d chips!",
			m.roundState.Ante,
			m.roundState.CurrentBlind.Name,
			m.roundState.TargetScore)

		// Create new shuffled deck
		m.deck = NewDeck()
		Shuffle(m.deck)

		// Deal new hand
		m.hand = make([]Card, 8)
		copy(m.hand, m.deck[:8])
		m.deck = m.deck[8:]

		// Reset played cards
		m.playedCards = make([]Card, 0)
		m.selectedCardIndex = 0
		m.lastScore = ScoreCalculation{}

		return m, nil

	case "q":
		// Quit to menu
		m.state = StateLanding
		m.focusedComponent = "landing"
		return m, nil
	}

	return m, nil
}

// handleGameOverKeys handles game over screen
func (m Model) handleGameOverKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		// Restart game
		m.newGame()
		return m, nil

	case "q":
		// Quit to menu
		m.state = StateLanding
		m.focusedComponent = "landing"
		return m, nil
	}

	return m, nil
}

// handleMainKeys handles keys in main view
func (m Model) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	// Navigation
	case "up", "k":
		return m.moveUp()

	case "down", "j":
		return m.moveDown()

	case "left", "h":
		return m.moveLeft()

	case "right", "l":
		return m.moveRight()

	case "pgup":
		return m.pageUp()

	case "pgdown":
		return m.pageDown()

	case "home", "g":
		return m.moveToTop()

	case "end", "G":
		return m.moveToBottom()

	// Actions
	case "enter":
		return m.selectItem()

	case " ": // space
		return m.toggleSelection()

	case "tab":
		return m.switchFocus()

	// Add your application-specific keys here
	// Example:
	// case "d":
	//     return m.deleteItem()
	//
	// case "e":
	//     return m.editItem()
	//
	// case "/":
	//     return m.startSearch()
	}

	return m, nil
}

// Navigation helper functions

func (m Model) moveUp() (tea.Model, tea.Cmd) {
	// Implement up navigation
	// Example: m.cursor = max(0, m.cursor-1)
	return m, nil
}

func (m Model) moveDown() (tea.Model, tea.Cmd) {
	// Implement down navigation
	// Example: m.cursor = min(len(m.items)-1, m.cursor+1)
	return m, nil
}

func (m Model) moveLeft() (tea.Model, tea.Cmd) {
	// Implement left navigation
	return m, nil
}

func (m Model) moveRight() (tea.Model, tea.Cmd) {
	// Implement right navigation
	return m, nil
}

func (m Model) pageUp() (tea.Model, tea.Cmd) {
	// Implement page up
	// Example: m.cursor = max(0, m.cursor-m.viewportHeight)
	return m, nil
}

func (m Model) pageDown() (tea.Model, tea.Cmd) {
	// Implement page down
	// Example: m.cursor = min(len(m.items)-1, m.cursor+m.viewportHeight)
	return m, nil
}

func (m Model) moveToTop() (tea.Model, tea.Cmd) {
	// Example: m.cursor = 0
	return m, nil
}

func (m Model) moveToBottom() (tea.Model, tea.Cmd) {
	// Example: m.cursor = len(m.items) - 1
	return m, nil
}

// Action helper functions

func (m Model) selectItem() (tea.Model, tea.Cmd) {
	// Implement item selection
	return m, nil
}

func (m Model) toggleSelection() (tea.Model, tea.Cmd) {
	// Implement toggle selection
	return m, nil
}

func (m Model) switchFocus() (tea.Model, tea.Cmd) {
	// Implement focus switching between components
	return m, nil
}

func (m Model) showHelp() (tea.Model, tea.Cmd) {
	// Show help dialog
	m.statusMsg = "Help: q=quit, ?=help, ↑↓=navigate, enter=select"
	return m, nil
}

func (m Model) refresh() (tea.Model, tea.Cmd) {
	// Refresh the current view
	m.statusMsg = "Refreshed"
	return m, nil
}

// Key bindings definition
type keyMap struct {
	Quit    key.Binding
	Help    key.Binding
	Refresh key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Select  key.Binding
	Toggle  key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "refresh"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	),
}
