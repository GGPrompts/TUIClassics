package hero

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Constants for game layout
const (
	NoteAreaHeight = 25 // Number of rows for notes to scroll through
	GracePeriod    = 3  // Extra rows after hit zone before note is considered missed
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateFinished
)

// HitResult represents the quality of a note hit
type HitResult int

const (
	HitMiss HitResult = iota
	HitOK       // Within 150ms = 25 points
	HitGood     // Within 100ms = 50 points
	HitPerfect  // Within 50ms = 100 points
)

// Note represents a single note that falls down a lane
type Note struct {
	Lane    int       // 0-4 (A/S/D/F/J)
	Y       int       // Current Y position on screen
	HitTime time.Time // When note should be hit
	Hit     bool      // Has been successfully hit
}

// Model holds the game state
type Model struct {
	// Game state
	state      GameState
	notes      []Note
	score      int
	combo      int
	multiplier int
	lastHit    HitResult

	// Song data
	currentSong *Song
	songIndex   int

	// Visual dimensions
	termWidth  int
	termHeight int
	laneWidth  int

	// Timing
	startTime   time.Time
	currentTime time.Time
	scrollSpeed float64 // Rows per second

	// Feedback display
	showHitFeedback bool
	hitFeedbackTime time.Time
}

// TickMsg is sent on each animation frame (60 FPS)
type TickMsg time.Time

// Song represents a playable song with a chart
type Song struct {
	Title    string
	Artist   string
	BPM      int
	Duration float64 // seconds
	Chart    []ChartNote
}

// ChartNote represents a note in the song chart
type ChartNote struct {
	Time float64 // seconds from start
	Lane int     // 0-4
}

// tickCmd returns a command that sends TickMsg every 16ms (60 FPS)
func tickCmd() tea.Cmd {
	return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
