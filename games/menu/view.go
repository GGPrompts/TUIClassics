package menu

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the menu
func (m Model) View() string {
	// If in game, delegate to the game's view
	if m.state == StateInGame && m.currentGame != nil {
		return m.currentGame.View()
	}

	// Render main menu
	return m.renderMainMenu()
}

func (m Model) renderMainMenu() string {
	// Use the Windows 95-style landing page
	if m.landingPage != nil {
		return m.landingPage.Render()
	}

	// Fallback to simple menu if landing page not initialized
	var b strings.Builder

	// Title
	title := titleStyle.Render("TUI CLASSICS")
	subtitle := subtitleStyle.Render("Classic Terminal Games Collection")

	b.WriteString(title + "\n")
	b.WriteString(subtitle + "\n\n")

	// Game list
	for i, game := range m.games {
		var line string

		// Format: [hotkey] Game Name
		prefix := fmt.Sprintf("[%s] ", game.Hotkey)

		if game.NewFunc == nil {
			// Game not implemented yet
			line = disabledMenuItemStyle.Render(prefix + game.Name)
		} else if i == m.selectedIdx {
			// Selected game
			line = selectedMenuItemStyle.Render(prefix + game.Name)
		} else {
			// Regular game
			line = menuItemStyle.Render(prefix + game.Name)
		}

		b.WriteString(line + "\n")

		// Show description for selected item
		if i == m.selectedIdx {
			desc := descriptionStyle.Render(game.Description)
			b.WriteString(desc + "\n")
		}

		b.WriteString("\n")
	}

	// Footer hints
	b.WriteString("\n")
	hints := hintStyle.Render("↑/↓: Navigate  •  Enter/Hotkey: Select  •  q: Quit")
	b.WriteString(hints)

	// Center the content
	content := b.String()
	return lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
