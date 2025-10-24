package balatro

// types.go - Type Definitions
// Purpose: All type definitions, structs, enums, and constants
// When to extend: Add new types here when introducing new data structures

// AppState represents the current application screen
type AppState int

const (
	StateLanding AppState = iota
	StateGame
	StateShop
	StateCollection
	StateOptions
)

// GamePhase represents the different phases within the game
type GamePhase int

const (
	PhaseSelectCards GamePhase = iota
	PhaseBlindComplete  // Victory screen after beating a blind
	PhaseShop           // Shopping between blinds
	PhaseGameOver       // Loss screen
)

// SortMode represents how cards in the hand are sorted
type SortMode int

const (
	SortNone SortMode = iota // No sorting (draw order)
	SortBySuit                // Sort by suit (Clubs, Diamonds, Hearts, Spades)
	SortByRank                // Sort by rank (A, 2, 3, ..., K)
)

// Model represents the application state
type Model struct {
	// Configuration
	config Config

	// UI State
	width  int
	height int
	state  AppState

	// Focus management
	focusedComponent string

	// Error handling
	err       error
	statusMsg string

	// Landing page
	landingPage *LandingPage

	// Game State
	deck              []Card           // Current deck
	hand              []Card           // Cards in player's hand
	playedCards       []Card           // Cards selected for play
	selectedCardIndex int              // Currently selected card for detail view (-1 = none)
	sortMode          SortMode         // Current hand sorting mode
	currentHandInfo   HandInfo         // Info about the currently selected hand
	lastScore         ScoreCalculation // Last score calculation

	// Phase 3: Jokers & Effects
	jokers        []Joker // Player's joker cards (max 5)
	shopJokers    []Joker // Jokers available for purchase in shop (2 per visit)
	selectedShopItem int  // Selected shop item index (-1 = none)

	// Phase 2: Round progression
	gamePhase    GamePhase   // Current game phase (playing, victory, shop, game over)
	roundState   RoundState  // Round state (ante, blind, scores, hands/discards)

	// Legacy fields (will be removed once fully integrated with RoundState)
	money             int              // Player's money
	currentRound      int              // Current round number
	handsRemaining    int              // Hands left to play this round
	discardsRemaining int              // Discards left this round
	targetScore       int              // Score needed to beat the round
	currentScore      int              // Score accumulated this round

	// Mouse state tracking
	mousePressX int // X position where mouse was pressed
	mousePressY int // Y position where mouse was pressed
}

// Config holds application configuration
type Config struct {
	// Theme
	Theme       string
	CustomTheme ThemeColors

	// Keybindings
	Keybindings       string
	CustomKeybindings map[string]string

	// Layout
	Layout LayoutConfig

	// UI Elements
	UI UIConfig

	// Performance
	Performance PerformanceConfig

	// Logging
	Logging LogConfig
}

// ThemeColors defines a color theme
type ThemeColors struct {
	Primary    string
	Secondary  string
	Background string
	Foreground string
	Accent     string
	Error      string
}

// LayoutConfig defines layout settings
type LayoutConfig struct {
	Type        string  // single, dual_pane, multi_panel, tabbed
	SplitRatio  float64 // For dual_pane
	ShowDivider bool
}

// UIConfig defines UI element settings
type UIConfig struct {
	ShowTitle       bool
	ShowStatus      bool
	ShowLineNumbers bool
	MouseEnabled    bool
	ShowIcons       bool
	IconSet         string
}

// PerformanceConfig defines performance settings
type PerformanceConfig struct {
	LazyLoading     bool
	CacheSize       int
	AsyncOperations bool
}

// LogConfig defines logging settings
type LogConfig struct {
	Enabled bool
	Level   string
	File    string
}

// Custom message types
// Add your application-specific messages here

type errMsg struct {
	err error
}

type statusMsg struct {
	message string
}

type resizeMsg struct {
	width  int
	height int
}

// Animation tick message for landing page
type tickMsg struct{}

// Game-specific messages
type startGameMsg struct{}

type cardSelectedMsg struct {
	index int
}

type playHandMsg struct{}

type discardMsg struct{}

type roundCompleteMsg struct {
	won bool
}
