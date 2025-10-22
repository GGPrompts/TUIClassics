package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	game2048 "github.com/GGPrompts/TUIClassics/games/2048"
	"github.com/GGPrompts/TUIClassics/games/balatro"
	"github.com/GGPrompts/TUIClassics/games/hero"
	"github.com/GGPrompts/TUIClassics/games/minesweeper"
	"github.com/GGPrompts/TUIClassics/games/snake"
	"github.com/GGPrompts/TUIClassics/games/solitaire"
)

// New creates a new menu model with all available games
func New() Model {
	games := []GameInfo{
		{
			Name:        "2048",
			Description: "Slide tiles to reach 2048",
			Hotkey:      "2",
			NewFunc:     func() tea.Model { return game2048.New() },
		},
		{
			Name:        "Balatro 🚧",
			Description: "Poker roguelike with jokers and scoring combos",
			Hotkey:      "b",
			NewFunc:     func() tea.Model { return balatro.New() },
		},
		{
			Name:        "Keyboard Hero",
			Description: "Rhythm game - hit keys to the beat!",
			Hotkey:      "h",
			NewFunc:     func() tea.Model { return hero.New() },
		},
		{
			Name:        "Minesweeper",
			Description: "Classic mine-finding puzzle game",
			Hotkey:      "m",
			NewFunc:     func() tea.Model { return minesweeper.New() },
		},
		{
			Name:        "Snake",
			Description: "Classic snake game - eat, grow, survive!",
			Hotkey:      "n",
			NewFunc:     func() tea.Model { return snake.New() },
		},
		{
			Name:        "Solitaire",
			Description: "Klondike card game",
			Hotkey:      "s",
			NewFunc:     func() tea.Model { return solitaire.New() },
		},
	}

	m := Model{
		state:       StateMainMenu,
		games:       games,
		selectedIdx: 0,
	}
	m.landingPage = NewLandingPage(80, 24, games) // Default size, will be updated on WindowSizeMsg
	return m
}

// Init initializes the menu
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),  // Request terminal size immediately
		animationTick(),   // Start animation loop
	)
}

// launchGame starts the selected game
func (m *Model) launchGame(idx int) tea.Cmd {
	if idx < 0 || idx >= len(m.games) {
		return nil
	}

	game := m.games[idx]
	if game.NewFunc == nil {
		return nil // Game not implemented yet
	}

	m.currentGame = game.NewFunc()
	m.state = StateInGame

	return m.currentGame.Init()
}

// returnToMenu goes back to the main menu
func (m *Model) returnToMenu() {
	m.state = StateMainMenu
	m.currentGame = nil
}
