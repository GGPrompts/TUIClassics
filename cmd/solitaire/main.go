package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/GGPrompts/TUIClassics/games/solitaire"
)

func main() {
	p := tea.NewProgram(
		solitaire.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running solitaire: %v\n", err)
		os.Exit(1)
	}
}
