package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/GGPrompts/TUIClassics/games/minesweeper"
	"github.com/GGPrompts/TUIClassics/games/solitaire"
)

// New creates a new menu model with all available games
func New() Model {
	games := []GameInfo{
		{
			Name:        "Minesweeper",
			Description: "Classic mine-finding puzzle game",
			Hotkey:      "m",
			NewFunc:     func() tea.Model { return minesweeper.New() },
		},
		{
			Name:        "Solitaire",
			Description: "Klondike card game (Work in Progress)",
			Hotkey:      "s",
			NewFunc:     func() tea.Model { return solitaire.New() },
		},
		{
			Name:        "Tanks",
			Description: "Coming soon...",
			Hotkey:      "t",
			NewFunc:     nil, // Not implemented yet
		},
	}

	return Model{
		state:       StateMainMenu,
		games:       games,
		selectedIdx: 0,
	}
}

// Init initializes the menu
func (m Model) Init() tea.Cmd {
	return nil
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
