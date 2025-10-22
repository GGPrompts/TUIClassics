package hero

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the game
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.renderMenu()
	case StatePlaying:
		return m.renderGame()
	case StateFinished:
		return m.renderFinished()
	default:
		return ""
	}
}

// renderMenu shows the song selection menu
func (m Model) renderMenu() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("*** KEYBOARD HERO ***")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Song list
	songs := getDemoSongs()
	b.WriteString("Select a song:\n\n")
	for i, song := range songs {
		prefix := "  "
		if i == m.songIndex {
			prefix = "> "
		}
		b.WriteString(fmt.Sprintf("%s%d. %s - %s (BPM: %d)\n",
			prefix, i+1, song.Title, song.Artist, song.BPM))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Up/Down: Select  Enter: Play  q: Quit"))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderGame shows the active gameplay
func (m Model) renderGame() string {
	var b strings.Builder

	// Title
	songTitle := ""
	if m.currentSong != nil {
		songTitle = fmt.Sprintf("%s - %s", m.currentSong.Title, m.currentSong.Artist)
	}
	b.WriteString(titleStyle.Render(songTitle))
	b.WriteString("\n\n")

	// Render the lanes with notes
	b.WriteString(m.renderLanes())

	// Score display
	b.WriteString("\n")
	b.WriteString(m.renderScore())

	// Hit feedback
	if m.showHitFeedback {
		b.WriteString("\n")
		hitText := getHitResultText(m.lastHit)
		hitColor := getHitResultColor(m.lastHit)
		feedbackStyle := lipgloss.NewStyle().Foreground(hitColor).Bold(true)
		b.WriteString(feedbackStyle.Render(hitText))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("A S D F SPACE to hit notes  |  ESC: Menu  q: Quit"))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// renderLanes renders the 5 lanes with notes and hit zone
func (m Model) renderLanes() string {
	var b strings.Builder
	laneKeys := []string{"A", "S", "D", "F", "SPC"}

	// Top border and lane headers
	b.WriteString("     +")
	for i := 0; i < 5; i++ {
		b.WriteString(strings.Repeat("-", m.laneWidth))
		if i < 4 {
			b.WriteString("+")
		}
	}
	b.WriteString("+\n")

	// Lane key labels
	b.WriteString("     |")
	for i := 0; i < 5; i++ {
		keyStyle := laneHeaderStyle.Copy().Foreground(getLaneColor(i))
		label := keyStyle.Render(laneKeys[i])
		b.WriteString(centerString(label, m.laneWidth))
		b.WriteString("|")
	}
	b.WriteString("\n")

	// Separator
	b.WriteString("     +")
	for i := 0; i < 5; i++ {
		b.WriteString(strings.Repeat("-", m.laneWidth))
		if i < 4 {
			b.WriteString("+")
		}
	}
	b.WriteString("+\n")

	// Note area
	for y := 0; y < NoteAreaHeight; y++ {
		b.WriteString("     |")
		for lane := 0; lane < 5; lane++ {
			// Check if there's a note at this position
			hasNote := m.hasNoteAt(lane, y)

			if hasNote {
				noteSymbol := noteStyle.Copy().
					Foreground(getLaneColor(lane)).
					Render("O")
				b.WriteString(centerString(noteSymbol, m.laneWidth))
			} else {
				b.WriteString(strings.Repeat(" ", m.laneWidth))
			}

			b.WriteString("|")
		}
		b.WriteString("\n")
	}

	// Hit zone separator
	b.WriteString("     +")
	for i := 0; i < 5; i++ {
		b.WriteString(strings.Repeat("-", m.laneWidth))
		if i < 4 {
			b.WriteString("+")
		}
	}
	b.WriteString("+\n")

	// Hit zone
	b.WriteString("     |")
	for lane := 0; lane < 5; lane++ {
		hitBox := hitZoneStyle.Copy().
			Foreground(getLaneColor(lane)).
			Render("[" + laneKeys[lane] + "]")
		b.WriteString(centerString(hitBox, m.laneWidth))
		b.WriteString("|")
	}
	b.WriteString("\n")

	// Bottom border
	b.WriteString("     +")
	for i := 0; i < 5; i++ {
		b.WriteString(strings.Repeat("-", m.laneWidth))
		if i < 4 {
			b.WriteString("+")
		}
	}
	b.WriteString("+")

	return b.String()
}

// hasNoteAt checks if there's a note at a specific lane and Y position
func (m Model) hasNoteAt(lane, y int) bool {
	for _, note := range m.notes {
		if note.Lane == lane && note.Y == y && !note.Hit {
			return true
		}
	}
	return false
}

// renderScore displays score, combo, and multiplier
func (m Model) renderScore() string {
	return fmt.Sprintf("%s  %s  Multiplier: %dx",
		scoreStyle.Render(fmt.Sprintf("Score: %d", m.score)),
		comboStyle.Render(fmt.Sprintf("Combo: %dx", m.combo)),
		m.multiplier)
}

// renderFinished shows the end-of-song results
func (m Model) renderFinished() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("*** SONG COMPLETE! ***"))
	b.WriteString("\n\n")

	if m.currentSong != nil {
		songInfo := lipgloss.NewStyle().Align(lipgloss.Center).Render(
			fmt.Sprintf("Song: %s - %s", m.currentSong.Title, m.currentSong.Artist))
		b.WriteString(songInfo)
		b.WriteString("\n\n")
	}

	b.WriteString(scoreStyle.Render(fmt.Sprintf("Final Score: %d", m.score)))
	b.WriteString("\n")
	b.WriteString(comboStyle.Render(fmt.Sprintf("Max Combo: %dx", m.combo)))
	b.WriteString("\n")

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Enter: Play Again  ESC: Menu  q: Quit"))

	return lipgloss.Place(m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center, b.String())
}

// centerString centers a string within a given width
func centerString(s string, width int) string {
	// Remove ANSI codes to get actual length
	strLen := lipgloss.Width(s)
	if strLen >= width {
		return s
	}
	leftPad := (width - strLen) / 2
	rightPad := width - strLen - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}
