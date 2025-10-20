package solitaire

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// New creates a new Solitaire game model
func New() Model {
	m := Model{
		state:    StateMenu,
		drawMode: 1, // Start with Draw-1
		score:    0,
		moves:    0,
	}

	// Initialize piles
	for i := 0; i < 7; i++ {
		m.tableau[i] = Pile{Cards: []Card{}, PileType: TableauPile}
	}
	for i := 0; i < 4; i++ {
		m.foundation[i] = Pile{Cards: []Card{}, PileType: FoundationPile}
	}
	m.stock = Pile{Cards: []Card{}, PileType: StockPile}
	m.waste = Pile{Cards: []Card{}, PileType: WastePile}

	// Start cursor on first tableau pile
	m.cursor = CursorLocation{PileType: TableauPile, PileIndex: 0}

	return m
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

// NewGame starts a new game by dealing cards
func (m *Model) NewGame() {
	// Reset game state
	m.state = StatePlaying
	m.score = 0
	m.moves = 0
	m.startTime = time.Now()
	m.elapsedTime = 0

	// Clear piles
	for i := 0; i < 7; i++ {
		m.tableau[i].Cards = []Card{}
	}
	for i := 0; i < 4; i++ {
		m.foundation[i].Cards = []Card{}
	}
	m.stock.Cards = []Card{}
	m.waste.Cards = []Card{}

	// Create and shuffle deck
	deck := ShuffleDeck(NewDeck())

	// Deal to tableau (1-7 cards)
	deckIndex := 0
	for col := 0; col < 7; col++ {
		for row := 0; row <= col; row++ {
			card := deck[deckIndex]
			// Only the top card (last dealt) is face up
			if row == col {
				card.FaceUp = true
			}
			m.tableau[col].Cards = append(m.tableau[col].Cards, card)
			deckIndex++
		}
	}

	// Remaining cards go to stock
	for i := deckIndex; i < len(deck); i++ {
		m.stock.Cards = append(m.stock.Cards, deck[i])
	}

	// Reset cursor
	m.cursor = CursorLocation{PileType: TableauPile, PileIndex: 0}
	m.selectedCard = nil
	m.selectedPile = nil
}

// DrawFromStock draws a card from the stock pile
func (m *Model) DrawFromStock() {
	if len(m.stock.Cards) == 0 {
		// Recycle waste back to stock
		if len(m.waste.Cards) > 0 {
			m.stock.Cards = make([]Card, len(m.waste.Cards))
			// Reverse the waste pile and turn cards face down
			for i := 0; i < len(m.waste.Cards); i++ {
				card := m.waste.Cards[len(m.waste.Cards)-1-i]
				card.FaceUp = false
				m.stock.Cards[i] = card
			}
			m.waste.Cards = []Card{}
		}
		return
	}

	// Draw card(s) from stock to waste
	numToDraw := m.drawMode
	if numToDraw > len(m.stock.Cards) {
		numToDraw = len(m.stock.Cards)
	}

	for i := 0; i < numToDraw; i++ {
		card := m.stock.Cards[len(m.stock.Cards)-1]
		card.FaceUp = true
		m.stock.Cards = m.stock.Cards[:len(m.stock.Cards)-1]
		m.waste.Cards = append(m.waste.Cards, card)
	}

	m.moves++
}

// CheckWin checks if the player has won
func (m *Model) CheckWin() bool {
	// Win condition: all foundation piles have 13 cards
	for i := 0; i < 4; i++ {
		if len(m.foundation[i].Cards) != 13 {
			return false
		}
	}
	return true
}

// StartWaterfallAnimation initializes the waterfall animation
func (m *Model) StartWaterfallAnimation() {
	m.animating = true
	m.animationFrame = 0
	m.waterfallCards = []WaterfallCard{}

	// Create waterfall cards from all foundation piles
	startX := float64(m.termWidth) / 2
	for i := 0; i < 4; i++ {
		for j := 0; j < len(m.foundation[i].Cards); j++ {
			card := m.foundation[i].Cards[j]
			// Random horizontal velocity
			vx := (rand.Float64() - 0.5) * 4
			vy := rand.Float64() * 2

			m.waterfallCards = append(m.waterfallCards, WaterfallCard{
				card:     card,
				x:        startX + float64(i)*2,
				y:        5.0,
				vx:       vx,
				vy:       vy,
				rotation: 0,
			})
		}
	}

	// Animation runs for ~3 seconds at 60fps = 180 frames
	m.animationMaxSteps = 180
}

// ProgressWaterfallAnimation updates the waterfall animation
func (m *Model) ProgressWaterfallAnimation() {
	const (
		gravity    = 0.5
		elasticity = 0.7
	)

	m.animationFrame++

	for i := range m.waterfallCards {
		card := &m.waterfallCards[i]

		// Apply gravity
		card.vy += gravity

		// Update position
		card.x += card.vx
		card.y += card.vy

		// Bounce at bottom
		bottomY := float64(m.termHeight - 3)
		if card.y >= bottomY {
			card.y = bottomY
			card.vy = -card.vy * elasticity
			// Reduce horizontal velocity on bounce
			card.vx *= 0.9
		}

		// Bounce at sides
		if card.x <= 1 || card.x >= float64(m.termWidth-7) {
			card.vx = -card.vx
		}

		// Spin
		card.rotation += 0.2
	}

	// End animation
	if m.animationFrame >= m.animationMaxSteps {
		m.animating = false
		m.state = StateWon
	}
}

// AddScore adds points to the score
func (m *Model) AddScore(points int) {
	m.score += points
	if m.score < 0 {
		m.score = 0
	}
}
