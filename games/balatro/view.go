package balatro

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// view.go - View Rendering
// Purpose: Top-level view rendering and layout
// When to extend: Add new view modes or modify layout logic

// View renders the entire application
func (m Model) View() string {
	// Check if terminal size is sufficient
	if !m.isValidSize() {
		return m.renderMinimalView()
	}

	// Handle errors
	if m.err != nil {
		return m.renderErrorView()
	}

	// Render based on application state
	switch m.state {
	case StateLanding:
		return m.renderLandingPage()

	case StateGame:
		// Phase 2: Dispatch based on game phase
		switch m.gamePhase {
		case PhaseSelectCards:
			return m.renderGameView()
		case PhaseBlindComplete:
			return m.renderBlindCompleteScreen()
		case PhaseShop:
			return m.renderShopScreen()
		case PhaseGameOver:
			return m.renderGameOverScreen()
		default:
			return m.renderGameView()
		}

	case StateShop, StateCollection, StateOptions:
		// Render based on layout type for other states
		switch m.config.Layout.Type {
		case "single":
			return m.renderSinglePane()

		case "dual_pane":
			return m.renderDualPane()

		case "multi_panel":
			return m.renderMultiPanel()

		case "tabbed":
			return m.renderTabbed()

		default:
			return m.renderSinglePane()
		}

	default:
		return m.renderSinglePane()
	}
}

// renderLandingPage renders the trippy lava lamp landing page
func (m Model) renderLandingPage() string {
	if m.landingPage == nil {
		return "Loading..."
	}

	return m.landingPage.Render()
}

// renderSinglePane renders a single-pane layout
func (m Model) renderSinglePane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Main content
	sections = append(sections, m.renderMainContent())

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderDualPane renders a dual-pane layout (side-by-side)
func (m Model) renderDualPane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Calculate pane dimensions
	leftWidth, rightWidth := m.calculateDualPaneLayout()

	// Left pane
	leftPane := m.renderLeftPane(leftWidth)

	// Divider
	divider := ""
	if m.config.Layout.ShowDivider {
		divider = m.renderDivider()
	}

	// Right pane
	rightPane := m.renderRightPane(rightWidth)

	// Join panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, divider, rightPane)
	sections = append(sections, panes)

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMultiPanel renders a multi-panel layout
func (m Model) renderMultiPanel() string {
	// Implement multi-panel layout
	// This is a placeholder - customize based on your needs
	return m.renderSinglePane()
}

// renderTabbed renders a tabbed interface
func (m Model) renderTabbed() string {
	// Implement tabbed interface
	// This is a placeholder - customize based on your needs
	return m.renderSinglePane()
}

// Component rendering functions

// renderTitleBar renders the title bar
func (m Model) renderTitleBar() string {
	return titleStyle.Render("♠ ♥ BALATRO TUI ♦ ♣")
}

// renderStatusBar renders the status bar
func (m Model) renderStatusBar() string {
	return statusStyle.Render(m.statusMsg)
}

// renderMainContent renders the main content area
func (m Model) renderMainContent() string {
	contentWidth, contentHeight := m.calculateLayout()

	// Implement your main content rendering here
	// Example:
	// return m.renderItemList(contentWidth, contentHeight)

	placeholder := "Main content area\n\n"
	placeholder += "Implement your content rendering in renderMainContent()\n\n"
	placeholder += "Press ? for help\n"
	placeholder += "Press q to quit"

	return contentStyle.Width(contentWidth).Height(contentHeight).Render(placeholder)
}

// renderLeftPane renders the left pane in dual-pane mode
func (m Model) renderLeftPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement left pane content
	content := "Left Pane\n\n"
	content += "Width: " + string(rune(width))

	return leftPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderRightPane renders the right pane in dual-pane mode
func (m Model) renderRightPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement right pane content
	content := "Right Pane (Preview)\n\n"
	content += "Width: " + string(rune(width))

	return rightPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderDivider renders the vertical divider between panes
func (m Model) renderDivider() string {
	_, contentHeight := m.calculateLayout()
	divider := strings.Repeat("│\n", contentHeight)
	return dividerStyle.Render(divider)
}

// Error and minimal views

// renderErrorView renders an error message
func (m Model) renderErrorView() string {
	content := "Error: " + m.err.Error() + "\n\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// renderMinimalView renders a minimal view for small terminals
func (m Model) renderMinimalView() string {
	content := "Terminal too small\n"
	content += "Minimum: 40x10\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// Helper functions

// truncateString truncates a string to fit within maxWidth
func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}

// padRight pads a string with spaces to reach the desired width
func padRight(s string, width int) string {
	currentWidth := lipgloss.Width(s)
	if currentWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-currentWidth)
}

// centerString centers a string within the given width
func centerString(s string, width int) string {
	strWidth := lipgloss.Width(s)
	if strWidth >= width {
		return s
	}
	leftPad := (width - strWidth) / 2
	rightPad := width - strWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// Game-specific rendering functions

// renderGameView renders the main poker game view (centered)
func (m Model) renderGameView() string {
	var sections []string

	// Game info (round, score, money) - centered
	gameInfo := m.renderGameInfo()
	centeredGameInfo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(gameInfo)
	sections = append(sections, centeredGameInfo)
	sections = append(sections, "")

	// Jokers display (if any)
	if len(m.jokers) > 0 {
		sections = append(sections, m.renderJokers())
		sections = append(sections, "")
	}

	// Hand display
	sections = append(sections, m.renderHand())
	sections = append(sections, "")

	// Selected card detail panel
	if m.selectedCardIndex >= 0 && m.selectedCardIndex < len(m.hand) {
		sections = append(sections, m.renderCardDetail(m.hand[m.selectedCardIndex]))
		sections = append(sections, "")
	}

	// Current hand info (always reserve space to prevent UI jumping)
	sections = append(sections, m.renderCurrentHandInfo())
	sections = append(sections, "")

	// Score breakdown (if we just scored)
	if m.lastScore.FinalScore > 0 {
		sections = append(sections, m.renderScoreBreakdown())
		sections = append(sections, "")
	}

	// Controls - centered
	controls := m.renderGameControls()
	centeredControls := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(controls)
	sections = append(sections, centeredControls)

	// Join all content
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Calculate vertical centering
	contentHeight := lipgloss.Height(content)
	availableHeight := m.height - 2 // -2 for title and status bars

	// Only center if content fits; otherwise start from top
	topPadding := 0
	if contentHeight < availableHeight {
		topPadding = (availableHeight - contentHeight) / 2
	}

	// Add vertical padding
	var result strings.Builder

	// Title bar (centered)
	if m.config.UI.ShowTitle {
		title := m.renderTitleBar()
		centeredTitle := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(title)
		result.WriteString(centeredTitle)
		result.WriteString("\n")
	}

	// Top padding (only if content fits)
	for i := 0; i < topPadding; i++ {
		result.WriteString("\n")
	}

	// Center content horizontally using lipgloss.Place
	centeredContent := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Top,
		content,
	)
	result.WriteString(centeredContent)

	// Status bar (centered)
	if m.config.UI.ShowStatus {
		result.WriteString("\n")
		status := m.renderStatusBar()
		centeredStatus := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(status)
		result.WriteString(centeredStatus)
	}

	return result.String()
}

// renderGameInfo renders round/score/money info
func (m Model) renderGameInfo() string {
	// Phase 2: Use roundState if initialized
	var parts []string

	// Default style for most text
	defaultStyle := lipgloss.NewStyle().
		Foreground(neonColors.Gold).
		Bold(true)

	// Blue style for Hands
	handsStyle := lipgloss.NewStyle().
		Foreground(neonColors.Blue).
		Bold(true)

	// Red style for Discards
	discardsStyle := lipgloss.NewStyle().
		Foreground(neonColors.Red).
		Bold(true)

	if m.roundState.Ante > 0 {
		// Show blind name and ante
		blindName := m.roundState.CurrentBlind.Name
		baseInfo := fmt.Sprintf("%s (Ante %d) | Score: %d/%d | Money: $%d | ",
			blindName,
			m.roundState.Ante,
			m.roundState.CurrentScore,
			m.roundState.TargetScore,
			m.roundState.Money)
		parts = append(parts, defaultStyle.Render(baseInfo))

		// Add colored Hands
		parts = append(parts, handsStyle.Render(fmt.Sprintf("Hands: %d", m.roundState.HandsRemaining)))
		parts = append(parts, defaultStyle.Render(" | "))

		// Add colored Discards
		parts = append(parts, discardsStyle.Render(fmt.Sprintf("Discards: %d", m.roundState.DiscardsRemaining)))
	} else {
		// Legacy format (fallback)
		baseInfo := fmt.Sprintf("Round %d | Score: %d/%d | Money: $%d | ",
			m.currentRound,
			m.currentScore,
			m.targetScore,
			m.money)
		parts = append(parts, defaultStyle.Render(baseInfo))

		// Add colored Hands
		parts = append(parts, handsStyle.Render(fmt.Sprintf("Hands: %d", m.handsRemaining)))
		parts = append(parts, defaultStyle.Render(" | "))

		// Add colored Discards
		parts = append(parts, discardsStyle.Render(fmt.Sprintf("Discards: %d", m.discardsRemaining)))
	}

	return strings.Join(parts, "")
}

// renderJokers renders owned jokers in a compact horizontal display
func (m Model) renderJokers() string {
	if len(m.jokers) == 0 {
		return ""
	}

	var jokerCards []string
	for _, joker := range m.jokers {
		jokerCards = append(jokerCards, m.renderJokerCard(joker))
	}

	// Join jokers horizontally with spacing
	jokersRow := lipgloss.JoinHorizontal(lipgloss.Top, jokerCards...)

	// Label
	label := lipgloss.NewStyle().
		Foreground(neonColors.Purple).
		Bold(true).
		Render("JOKERS")

	return lipgloss.JoinVertical(lipgloss.Center, label, jokersRow)
}

// renderJokerCard renders a single joker card (compact)
func (m Model) renderJokerCard(joker Joker) string {
	// Card border style based on rarity
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(joker.Rarity.Color()).
		Width(18).
		Height(4).
		Padding(0, 1).
		MarginRight(1)

	// Content
	var content strings.Builder

	// Name (colored by rarity)
	nameStyle := lipgloss.NewStyle().
		Foreground(joker.Rarity.Color()).
		Bold(true)
	content.WriteString(nameStyle.Render(joker.Name))
	content.WriteString("\n")

	// Description (wrap if needed)
	desc := joker.Description
	if len(desc) > 18 {
		// Truncate with ellipsis
		desc = desc[:15] + "..."
	}
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))
	content.WriteString(descStyle.Render(desc))

	return cardStyle.Render(content.String())
}

// renderShopJoker renders a joker for sale in the shop
func (m Model) renderShopJoker(joker Joker, selected bool, index int) string {
	// Card border style
	borderColor := joker.Rarity.Color()
	if selected {
		borderColor = neonColors.Gold // Highlight selected
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(22).
		Height(6).
		Padding(0, 1).
		MarginRight(2)

	// Content
	var content strings.Builder

	// Number for selection
	numStyle := lipgloss.NewStyle().
		Foreground(neonColors.Cyan).
		Bold(true)
	content.WriteString(numStyle.Render(fmt.Sprintf("[%d]", index+1)))
	content.WriteString("\n")

	// Name (colored by rarity)
	nameStyle := lipgloss.NewStyle().
		Foreground(joker.Rarity.Color()).
		Bold(true)
	content.WriteString(nameStyle.Render(joker.Name))
	content.WriteString("\n")

	// Description
	desc := joker.Description
	if len(desc) > 22 {
		desc = desc[:19] + "..."
	}
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250"))
	content.WriteString(descStyle.Render(desc))
	content.WriteString("\n")

	// Cost
	costStyle := lipgloss.NewStyle().
		Foreground(neonColors.Gold).
		Bold(true)
	content.WriteString(costStyle.Render(fmt.Sprintf("$%d", joker.GetCost())))

	return cardStyle.Render(content.String())
}

// renderHand renders all cards in the hand (centered)
func (m Model) renderHand() string {
	if len(m.hand) == 0 {
		return "No cards in hand"
	}

	var cards []string
	for i, card := range m.hand {
		isSelected := i == m.selectedCardIndex
		isPlayed := false
		for _, pc := range m.playedCards {
			if card.Suit == pc.Suit && card.Rank == pc.Rank {
				isPlayed = true
				break
			}
		}
		cards = append(cards, m.renderCard(card, isSelected, isPlayed))
	}

	// Join cards horizontally with space
	handDisplay := lipgloss.JoinHorizontal(lipgloss.Top, cards...)

	// Add number indicators below cards
	var numbers []string
	for i := range m.hand {
		num := fmt.Sprintf(" [%d] ", i+1)
		numStyle := lipgloss.NewStyle().Width(7).Align(lipgloss.Center)
		if i == m.selectedCardIndex {
			numStyle = numStyle.Foreground(neonColors.Gold).Bold(true)
		}
		numbers = append(numbers, numStyle.Render(num))
	}
	numbersDisplay := lipgloss.JoinHorizontal(lipgloss.Top, numbers...)

	// Combine hand and numbers
	handWithNumbers := lipgloss.JoinVertical(lipgloss.Left, handDisplay, numbersDisplay)

	// Center the entire hand display
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(handWithNumbers)
}

// renderCard renders a single compact card (5 chars wide)
func (m Model) renderCard(card Card, selected bool, played bool) string {
	// Top line: Rank + Suit symbol
	rankStr := card.Rank.ShortString()
	if len(rankStr) == 1 {
		rankStr = " " + rankStr
	}
	topLine := rankStr + " " + card.Suit.Symbol()

	// Middle line: Enhancement indicator
	midLine := card.Enhancement.ShortString()
	if midLine == "" {
		midLine = " "
	}

	// Bottom line: Edition/Seal indicator
	botLine := card.Edition.ShortString()
	if card.Seal != NoSeal {
		botLine += card.Seal.Symbol()
	}
	if botLine == "" {
		botLine = " "
	}

	content := topLine + "\n" + centerString(midLine, 5) + "\n" + centerString(botLine, 5)

	// Style with colored border
	borderColor := card.Enhancement.BorderColor()
	if selected {
		borderColor = neonColors.Gold
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(5).
		Height(3).
		Padding(0, 0)

	if played {
		style = style.Background(lipgloss.Color("235")) // Dim background
	}

	// Color the suit symbol
	styledContent := strings.Replace(content, card.Suit.Symbol(),
		lipgloss.NewStyle().Foreground(card.Suit.Color()).Render(card.Suit.Symbol()), 1)

	// Color the enhancement text
	if card.Enhancement != NoEnhancement {
		enhText := card.Enhancement.ShortString()
		styledEnhText := lipgloss.NewStyle().
			Foreground(card.Enhancement.TextColor()).
			Render(enhText)
		styledContent = strings.Replace(styledContent, enhText, styledEnhText, 1)
	}

	return style.Render(styledContent)
}

// renderCardDetail renders detailed info about a card (compact & centered)
func (m Model) renderCardDetail(card Card) string {
	var lines []string

	// Title: ACE OF SPADES
	title := card.Rank.String() + " OF " + card.Suit.String()
	styledTitle := lipgloss.NewStyle().
		Foreground(card.Suit.Color()).
		Bold(true).
		Render(title)
	lines = append(lines, styledTitle)

	// Chip value
	lines = append(lines, fmt.Sprintf("Value: %d", card.GetChipValue()))

	// Enhancement (compact - no description)
	if card.Enhancement != NoEnhancement {
		enhStyle := lipgloss.NewStyle().
			Foreground(card.Enhancement.TextColor())
		lines = append(lines, fmt.Sprintf("Enh: %s",
			enhStyle.Render(card.Enhancement.String())))
	}

	// Edition
	if card.Edition != NoEdition {
		lines = append(lines, fmt.Sprintf("Ed: %s", card.Edition.String()))
	}

	// Seal
	if card.Seal != NoSeal {
		lines = append(lines, fmt.Sprintf("Seal: %s", card.Seal.String()))
	}

	content := strings.Join(lines, "\n")

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neonColors.Cyan).
		Padding(0, 1).
		Width(28).
		Render(content)

	// Center the panel
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(panel)
}

// renderCurrentHandInfo shows info about the currently selected cards for play (centered)
func (m Model) renderCurrentHandInfo() string {
	var info string

	if len(m.playedCards) == 0 {
		// Return placeholder to reserve space and prevent UI jumping
		info = " " // Single space to maintain height
	} else {
		handInfo := m.currentHandInfo
		handName := handInfo.Type.String()

		style := lipgloss.NewStyle().
			Foreground(neonColors.Magenta).
			Bold(true)

		info = style.Render(fmt.Sprintf("Selected Hand: %s (%d chips × %d mult)",
			handName, handInfo.BaseChips, handInfo.BaseMult))
	}

	// Center it
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(info)
}

// renderScoreBreakdown shows the detailed score calculation (compact & centered)
func (m Model) renderScoreBreakdown() string {
	score := m.lastScore

	// Compact format: single line with breakdown
	scoreText := lipgloss.NewStyle().
		Foreground(neonColors.Green).
		Bold(true).
		Render(fmt.Sprintf("SCORE: %d", score.FinalScore))

	calculation := fmt.Sprintf("%d chips × %d mult = %d",
		score.TotalChips, score.TotalMult, score.FinalScore)

	// Show modifiers only if they exist (keep it minimal)
	var modSummary string
	chipMods := len(score.ChipModifiers)
	multMods := len(score.MultModifiers)
	if chipMods > 0 || multMods > 0 {
		modSummary = fmt.Sprintf("(%d chip mods, %d mult mods)", chipMods, multMods)
	}

	// Compact 2-3 line format
	content := scoreText + "\n" + calculation
	if modSummary != "" {
		content += "\n" + modSummary
	}

	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neonColors.Green).
		Padding(0, 2).
		Render(content)

	// Center the panel
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(panel)
}

// renderGameControls shows keyboard controls
func (m Model) renderGameControls() string {
	// Base controls
	controls := "1-8: Select | Space: Toggle for play | Enter: Play hand | D: Discard | S: Sort | Q: Quit"

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)

	// Show current sort mode
	var sortModeText string
	switch m.sortMode {
	case SortBySuit:
		sortModeText = "Sort: SUIT (♣→♦→♥→♠)"
	case SortByRank:
		sortModeText = "Sort: RANK (A→K)"
	case SortNone:
		sortModeText = "Sort: NONE"
	}

	sortStyle := lipgloss.NewStyle().
		Foreground(neonColors.Cyan).
		Bold(true)

	// Combine controls and sort mode indicator
	return style.Render(controls) + "\n" + sortStyle.Render(sortModeText)
}

// Phase 2: New game state screens

// renderBlindCompleteScreen shows victory screen after beating a blind
func (m Model) renderBlindCompleteScreen() string {
	var sections []string

	// Victory title
	title := lipgloss.NewStyle().
		Foreground(neonColors.Green).
		Bold(true).
		Render("★ BLIND COMPLETE! ★")
	sections = append(sections, title)
	sections = append(sections, "")

	// Blind info
	blindInfo := fmt.Sprintf("Defeated: %s (Ante %d)",
		m.roundState.CurrentBlind.Name,
		m.roundState.Ante)
	sections = append(sections, blindInfo)

	// Score achieved
	scoreText := fmt.Sprintf("Score: %d / %d",
		m.roundState.CurrentScore,
		m.roundState.TargetScore)
	scoreStyle := lipgloss.NewStyle().
		Foreground(neonColors.Magenta).
		Bold(true)
	sections = append(sections, scoreStyle.Render(scoreText))
	sections = append(sections, "")

	// Money earned
	moneyEarned := m.roundState.CurrentBlind.Reward
	moneyText := fmt.Sprintf("💰 Earned: $%d (Total: $%d)",
		moneyEarned,
		m.roundState.Money)
	moneyStyle := lipgloss.NewStyle().
		Foreground(neonColors.Gold).
		Bold(true)
	sections = append(sections, moneyStyle.Render(moneyText))
	sections = append(sections, "")

	// Next blind info
	nextBlindText := m.getNextBlindText()
	sections = append(sections, nextBlindText)
	sections = append(sections, "")

	// Controls
	controls := "Press Enter to continue to shop | Q: Quit"
	controlsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)
	sections = append(sections, controlsStyle.Render(controls))

	// Join all sections
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Create border panel
	panel := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(neonColors.Green).
		Padding(2, 4).
		Render(content)

	// Center the panel on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		panel,
	)
}

// renderGameOverScreen shows loss screen when running out of hands
func (m Model) renderGameOverScreen() string {
	var sections []string

	// Game over title
	title := lipgloss.NewStyle().
		Foreground(neonColors.Red).
		Bold(true).
		Render("☠ GAME OVER ☠")
	sections = append(sections, title)
	sections = append(sections, "")

	// Failed blind
	blindInfo := fmt.Sprintf("Failed: %s (Ante %d)",
		m.roundState.CurrentBlind.Name,
		m.roundState.Ante)
	sections = append(sections, blindInfo)

	// Final score
	scoreText := fmt.Sprintf("Score: %d / %d",
		m.roundState.CurrentScore,
		m.roundState.TargetScore)
	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")) // Bright red
	sections = append(sections, scoreStyle.Render(scoreText))
	sections = append(sections, "")

	// Final money
	moneyText := fmt.Sprintf("Final Money: $%d", m.roundState.Money)
	sections = append(sections, moneyText)
	sections = append(sections, "")

	// Stats
	statsText := fmt.Sprintf("Reached: Ante %d - %s",
		m.roundState.Ante,
		m.roundState.GetProgress())
	sections = append(sections, statsText)
	sections = append(sections, "")

	// Controls
	controls := "Press R to restart | Q: Quit to menu"
	controlsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)
	sections = append(sections, controlsStyle.Render(controls))

	// Join all sections
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Create border panel
	panel := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(neonColors.Red).
		Padding(2, 4).
		Render(content)

	// Center the panel on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		panel,
	)
}

// renderShopScreen shows basic shop interface (Phase 2: simple version)
func (m Model) renderShopScreen() string {
	var sections []string

	// Shop title
	title := lipgloss.NewStyle().
		Foreground(neonColors.Cyan).
		Bold(true).
		Render("🛍 SHOP 🛍")
	sections = append(sections, title)
	sections = append(sections, "")

	// Money display
	moneyText := fmt.Sprintf("Money: $%d", m.roundState.Money)
	moneyStyle := lipgloss.NewStyle().
		Foreground(neonColors.Gold).
		Bold(true)
	sections = append(sections, moneyStyle.Render(moneyText))
	sections = append(sections, "")

	// Next blind preview
	nextBlindText := m.getNextBlindText()
	sections = append(sections, nextBlindText)
	sections = append(sections, "")

	// Shop items - Jokers for sale
	if len(m.shopJokers) > 0 {
		jokersLabel := lipgloss.NewStyle().
			Foreground(neonColors.Purple).
			Bold(true).
			Render("JOKERS FOR SALE")
		sections = append(sections, jokersLabel)
		sections = append(sections, "")

		// Render each shop joker
		var shopJokerCards []string
		for i, joker := range m.shopJokers {
			selected := (i == m.selectedShopItem)
			shopJokerCards = append(shopJokerCards, m.renderShopJoker(joker, selected, i))
		}
		jokersRow := lipgloss.JoinHorizontal(lipgloss.Top, shopJokerCards...)
		sections = append(sections, jokersRow)
		sections = append(sections, "")
	}

	// Current jokers owned
	if len(m.jokers) > 0 {
		ownedLabel := lipgloss.NewStyle().
			Foreground(neonColors.Gold).
			Bold(true).
			Render(fmt.Sprintf("YOUR JOKERS (%d/5)", len(m.jokers)))
		sections = append(sections, ownedLabel)
		sections = append(sections, "")

		var ownedJokerCards []string
		for _, joker := range m.jokers {
			ownedJokerCards = append(ownedJokerCards, m.renderJokerCard(joker))
		}
		ownedRow := lipgloss.JoinHorizontal(lipgloss.Top, ownedJokerCards...)
		sections = append(sections, ownedRow)
		sections = append(sections, "")
	}

	// Controls
	controls := "1-2: Select | Enter: Buy | Space: Continue | Q: Quit"
	controlsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		Italic(true)
	sections = append(sections, controlsStyle.Render(controls))

	// Join all sections
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Create border panel
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neonColors.Cyan).
		Padding(2, 4).
		Render(content)

	// Center the panel on screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		panel,
	)
}

// getNextBlindText returns text describing the next blind
func (m Model) getNextBlindText() string {
	// Calculate next blind
	nextBlindProgress := m.roundState.BlindProgress + 1
	nextAnte := m.roundState.Ante

	if nextBlindProgress > 2 {
		nextBlindProgress = 0
		nextAnte++
	}

	nextBlind := GetBlind(nextAnte, BlindType(nextBlindProgress))

	return fmt.Sprintf("Next: %s (Ante %d) - Target: %d chips",
		nextBlind.Name,
		nextAnte,
		nextBlind.TargetScore)
}
