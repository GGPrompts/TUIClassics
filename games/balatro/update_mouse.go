package balatro

import (
	tea "github.com/charmbracelet/bubbletea"
)

// update_mouse.go - Mouse Event Handling
// Purpose: All mouse input processing
// When to extend: Add new mouse interactions or clickable elements here

// handleMouseEvent handles mouse input
func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if !m.config.UI.MouseEnabled {
		return m, nil
	}

	switch msg.Type {
	case tea.MouseLeft:
		return m.handleLeftClick(msg)

	case tea.MouseRight:
		return m.handleRightClick(msg)

	case tea.MouseWheelUp:
		return m.handleWheelUp(msg)

	case tea.MouseWheelDown:
		return m.handleWheelDown(msg)

	case tea.MouseMotion:
		// Handle mouse motion if needed (for hover effects)
		return m.handleMouseMotion(msg)
	}

	return m, nil
}

// handleLeftClick handles left mouse button clicks
func (m Model) handleLeftClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	x, y := msg.X, msg.Y

	// Check if click is in a specific region
	// Example: check if clicked on an item in a list
	// if m.isInItemList(x, y) {
	//     itemIndex := m.getItemIndexAt(y)
	//     if itemIndex >= 0 && itemIndex < len(m.items) {
	//         m.cursor = itemIndex
	//         return m.selectItem()
	//     }
	// }

	// Check if clicked on UI elements
	if m.isInTitleBar(x, y) {
		return m.handleTitleBarClick(x, y)
	}

	if m.isInStatusBar(x, y) {
		return m.handleStatusBarClick(x, y)
	}

	// Add your application-specific click handlers here

	return m, nil
}

// handleRightClick handles right mouse button clicks
func (m Model) handleRightClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	x, y := msg.X, msg.Y

	// Example: show context menu
	// return m.showContextMenu(x, y)

	_ = x
	_ = y
	return m, nil
}

// handleWheelUp handles mouse wheel scroll up
func (m Model) handleWheelUp(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Scroll up in the focused component
	return m.moveUp()
}

// handleWheelDown handles mouse wheel scroll down
func (m Model) handleWheelDown(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Scroll down in the focused component
	return m.moveDown()
}

// handleMouseMotion handles mouse movement (for hover effects)
func (m Model) handleMouseMotion(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Example: highlight hovered item
	// x, y := msg.X, msg.Y
	// if m.isInItemList(x, y) {
	//     m.hoveredItem = m.getItemIndexAt(y)
	// }
	return m, nil
}

// Helper functions for click region detection

func (m Model) isInTitleBar(x, y int) bool {
	if !m.config.UI.ShowTitle {
		return false
	}
	return y < 2
}

func (m Model) isInStatusBar(x, y int) bool {
	if !m.config.UI.ShowStatus {
		return false
	}
	return y >= m.height-1
}

func (m Model) handleTitleBarClick(x, y int) (tea.Model, tea.Cmd) {
	// Example: click on breadcrumb navigation
	// or click on window control buttons
	_ = x
	_ = y
	return m, nil
}

func (m Model) handleStatusBarClick(x, y int) (tea.Model, tea.Cmd) {
	// Example: click on status bar items
	_ = x
	_ = y
	return m, nil
}

// Double-click detection (if needed)
type clickTracker struct {
	lastClickX    int
	lastClickY    int
	lastClickTime int64
}

var tracker clickTracker

func (m Model) isDoubleClick(msg tea.MouseMsg) bool {
	// Implement double-click detection
	// Compare with tracker.lastClickTime
	// Reset tracker.lastClickX, tracker.lastClickY, tracker.lastClickTime
	return false
}
