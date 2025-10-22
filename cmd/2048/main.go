package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	game2048 "github.com/GGPrompts/TUIClassics/games/2048"
)

func main() {
	p := tea.NewProgram(game2048.New(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running 2048: %v\n", err)
		os.Exit(1)
	}
}
