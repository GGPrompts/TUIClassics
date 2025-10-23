package solitaire

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Suit represents a card suit
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// Rank represents a card rank
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card represents a playing card
type Card struct {
	Suit   Suit
	Rank   Rank
	FaceUp bool
}

// PileType represents the type of pile
type PileType int

const (
	TableauPile PileType = iota
	FoundationPile
	StockPile
	WastePile
)

// Pile represents a pile of cards
type Pile struct {
	Cards    []Card
	PileType PileType
}

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateWon
	StateLost
	StateStats
)

// CursorLocation tracks which pile the cursor is on
type CursorLocation struct {
	PileType  PileType
	PileIndex int // For tableau (0-6) and foundation (0-3)
}

// Model represents the Solitaire game state
type Model struct {
	// Terminal dimensions
	termWidth  int
	termHeight int

	// Game state
	state         GameState
	previousState GameState // For returning from stats view
	drawMode      int       // 1 or 3
	score         int
	moves         int
	startTime     time.Time
	elapsedTime   time.Duration

	// Stats/Achievements
	achievements []string // Latest achievements to display

	// Piles
	tableau    [7]Pile
	foundation [4]Pile
	stock      Pile
	waste      Pile

	// Cursor for keyboard navigation
	cursor         CursorLocation
	cursorCardIndex int // Which card in tableau stack (for up/down navigation)
	selectedCard   *Card
	selectedPile   *CursorLocation
	selectedIndex  int // Index of card in pile

	// Mouse drag state
	draggingCard  *Card
	draggingCards []Card // For moving sequences
	dragFromPile  *CursorLocation
	dragFromIndex int
	mousePressX   int // Track where mouse was pressed
	mousePressY   int

	// Double-click detection
	lastClickTime time.Time
	lastClickX    int
	lastClickY    int

	// Animation state
	animating         bool
	animationFrame    int
	waterfallCards    []WaterfallCard
	animationMaxSteps int

	// Undo stack (for future)
	undoStack []Move
}

// WaterfallCard represents a card in the waterfall animation
type WaterfallCard struct {
	card     Card
	x, y     float64 // Position
	vx, vy   float64 // Velocity
	rotation float64 // Spin angle
}

// Move represents a game move (for undo)
type Move struct {
	fromPile  CursorLocation
	toPile    CursorLocation
	cards     []Card
	revealed  bool // Whether a card was revealed
	fromWaste bool // Special case for waste pile
}

// AnimationTickMsg is sent on each animation frame
type AnimationTickMsg time.Time

// TickMsg is sent each second for timer updates
type TickMsg time.Time

// animationTick returns a command that waits for the next animation frame (60fps)
func animationTick() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return AnimationTickMsg(t)
	})
}

// tickCmd returns a command that ticks every second
func tickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
