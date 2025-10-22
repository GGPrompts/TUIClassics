package balatro

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// landing_page.go - Trippy Lava Lamp Landing Page
// Purpose: Recreate Balatro's mesmerizing opening scene with wavy grid and floating metaballs

// Neon color palette for dark mode
var neonColors = struct {
	Gold    lipgloss.Color
	Cyan    lipgloss.Color
	Magenta lipgloss.Color
	Green   lipgloss.Color
	Orange  lipgloss.Color
	Red     lipgloss.Color
	Blue    lipgloss.Color
	Purple  lipgloss.Color
	Yellow  lipgloss.Color
	Pink    lipgloss.Color
}{
	Gold:    lipgloss.Color("226"),
	Cyan:    lipgloss.Color("51"),
	Magenta: lipgloss.Color("201"),
	Green:   lipgloss.Color("46"),
	Orange:  lipgloss.Color("214"),
	Red:     lipgloss.Color("196"),
	Blue:    lipgloss.Color("39"),
	Purple:  lipgloss.Color("129"),
	Yellow:  lipgloss.Color("226"),
	Pink:    lipgloss.Color("213"),
}

// WavyGrid represents an animated grid background with sine wave distortion
type WavyGrid struct {
	width    int
	height   int
	frame    int
	gridSize int
}

// NewWavyGrid creates a new wavy grid background
func NewWavyGrid(width, height int) *WavyGrid {
	return &WavyGrid{
		width:    width,
		height:   height,
		gridSize: 10,
	}
}

// Update advances the animation frame
func (g *WavyGrid) Update() {
	g.frame++
}

// Render generates the wavy grid as a string
func (g *WavyGrid) Render() string {
	var b strings.Builder

	// Create grid with wave distortion
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			// Calculate wave offset using sine
			waveX := math.Sin(float64(y)/5.0+float64(g.frame)/20.0) * 2
			waveY := math.Sin(float64(x)/5.0+float64(g.frame)/20.0) * 2

			// Determine if this position should be a grid line
			gridX := int(float64(x) + waveX)
			gridY := int(float64(y) + waveY)

			isGridLine := (gridX%g.gridSize == 0) || (gridY%g.gridSize == 0)

			var char string
			var color lipgloss.Color

			if isGridLine {
				// Grid line
				if gridX%g.gridSize == 0 && gridY%g.gridSize == 0 {
					char = "+"              // Intersection
					color = neonColors.Purple // Purple
				} else if gridX%g.gridSize == 0 {
					char = "│"                          // Vertical line
					color = lipgloss.Color("61") // Dark purple
				} else {
					char = "─"                          // Horizontal line
					color = lipgloss.Color("61")
				}
			} else {
				char = " "
				color = lipgloss.Color("0")
			}

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}
		if y < g.height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Blob represents a single floating blob for the lava lamp effect
type Blob struct {
	x, y       float64        // Position
	vx, vy     float64        // Velocity
	radius     float64        // Size
	color      lipgloss.Color // Color
	colorIndex int            // For tracking
}

// LavaLamp represents the metaball/lava lamp effect with floating blobs
type LavaLamp struct {
	blobs  []Blob
	width  int
	height int
	frame  int
}

// NewLavaLamp creates a new lava lamp effect with colorful blobs
func NewLavaLamp(width, height int) *LavaLamp {
	blobs := []Blob{
		{
			x:          float64(width) / 4,
			y:          float64(height) / 3,
			vx:         0.3,
			vy:         0.2,
			radius:     8,
			color:      neonColors.Magenta,
			colorIndex: 0,
		},
		{
			x:          float64(width) * 3 / 4,
			y:          float64(height) / 2,
			vx:         -0.25,
			vy:         0.15,
			radius:     10,
			color:      neonColors.Cyan,
			colorIndex: 1,
		},
		{
			x:          float64(width) / 2,
			y:          float64(height) * 2 / 3,
			vx:         0.2,
			vy:         -0.3,
			radius:     7,
			color:      neonColors.Yellow,
			colorIndex: 2,
		},
		{
			x:          float64(width) / 3,
			y:          float64(height) / 4,
			vx:         -0.15,
			vy:         0.25,
			radius:     9,
			color:      neonColors.Purple,
			colorIndex: 3,
		},
	}

	return &LavaLamp{
		blobs:  blobs,
		width:  width,
		height: height,
	}
}

// Update moves the blobs and adds organic wobble
func (l *LavaLamp) Update() {
	l.frame++

	for i := range l.blobs {
		blob := &l.blobs[i]

		// Update position
		blob.x += blob.vx
		blob.y += blob.vy

		// Bounce off edges
		if blob.x < 0 || blob.x > float64(l.width) {
			blob.vx = -blob.vx
		}
		if blob.y < 0 || blob.y > float64(l.height) {
			blob.vy = -blob.vy
		}

		// Add some organic wobble
		wobbleX := math.Sin(float64(l.frame+i*37)/30.0) * 0.05
		wobbleY := math.Cos(float64(l.frame+i*41)/25.0) * 0.05
		blob.vx += wobbleX
		blob.vy += wobbleY

		// Damping to prevent too much speed
		blob.vx *= 0.99
		blob.vy *= 0.99
	}
}

// Render generates the metaball effect with gradient characters
func (l *LavaLamp) Render() string {
	// Create field map
	field := make([][]float64, l.height)
	colorMap := make([][]int, l.height)

	for y := 0; y < l.height; y++ {
		field[y] = make([]float64, l.width)
		colorMap[y] = make([]int, l.width)

		for x := 0; x < l.width; x++ {
			// Calculate field strength from all blobs
			totalField := 0.0
			closestBlob := 0
			maxField := 0.0

			for i, blob := range l.blobs {
				dx := float64(x) - blob.x
				dy := float64(y) - blob.y
				distance := math.Sqrt(dx*dx + dy*dy)

				// Metaball field: 1/distance^2 * radius^2
				if distance > 0 {
					blobField := (blob.radius * blob.radius) / (distance * distance)
					totalField += blobField

					// Track which blob contributes most (for color)
					if blobField > maxField {
						maxField = blobField
						closestBlob = i
					}
				}
			}

			field[y][x] = totalField
			colorMap[y][x] = closestBlob
		}
	}

	// Render with gradient characters
	var b strings.Builder

	gradientChars := []string{" ", "░", "▒", "▓", "█"}

	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			strength := field[y][x]
			blobIndex := colorMap[y][x]

			// Choose character based on field strength
			var char string
			if strength < 0.3 {
				char = " " // Empty space
			} else if strength < 0.8 {
				char = gradientChars[1] // ░
			} else if strength < 1.5 {
				char = gradientChars[2] // ▒
			} else if strength < 2.5 {
				char = gradientChars[3] // ▓
			} else {
				char = gradientChars[4] // █
			}

			// Color from closest blob
			color := l.blobs[blobIndex].color

			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(char)

			b.WriteString(styled)
		}
		if y < l.height-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// Big ASCII art for "BALATRO" title
var balatroBigText = []string{
	"██████╗  █████╗ ██╗      █████╗ ████████╗██████╗  ██████╗ ",
	"██╔══██╗██╔══██╗██║     ██╔══██╗╚══██╔══╝██╔══██╗██╔═══██╗",
	"██████╔╝███████║██║     ███████║   ██║   ██████╔╝██║   ██║",
	"██╔══██╗██╔══██║██║     ██╔══██║   ██║   ██╔══██╗██║   ██║",
	"██████╔╝██║  ██║███████╗██║  ██║   ██║   ██║  ██║╚██████╔╝",
	"╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ",
}

// renderTitle renders the big BALATRO title with rainbow color cycling
func renderTitle(frame int) string {
	var b strings.Builder

	// Cycle through rainbow colors
	colors := []lipgloss.Color{
		neonColors.Red,
		neonColors.Orange,
		neonColors.Yellow,
		neonColors.Green,
		neonColors.Cyan,
		neonColors.Blue,
		neonColors.Magenta,
	}

	for lineIdx, line := range balatroBigText {
		for charIdx, char := range line {
			if char == ' ' || char == '╚' || char == '═' || char == '╝' {
				b.WriteRune(char)
			} else {
				// Rainbow effect: different color per character, cycling
				colorIdx := (charIdx + lineIdx + frame/5) % len(colors)
				colored := lipgloss.NewStyle().
					Foreground(colors[colorIdx]).
					Bold(true).
					Render(string(char))
				b.WriteString(colored)
			}
		}
		if lineIdx < len(balatroBigText)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// MenuItem represents a menu option
type MenuItem struct {
	Label    string
	Selected bool
}

// renderMenu renders the menu with gold selection indicator
func renderMenu(items []MenuItem, frame int) string {
	var b strings.Builder

	for i, item := range items {
		if item.Selected {
			// Selected: gold with pulse effect
			color := neonColors.Gold

			style := lipgloss.NewStyle().
				Foreground(color).
				Bold(true)

			// Add selection indicator
			indicator := "► "
			styledIndicator := lipgloss.NewStyle().
				Foreground(color).
				Render(indicator)

			b.WriteString(styledIndicator)
			b.WriteString(style.Render(item.Label))
		} else {
			// Not selected: dim white
			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color("250"))

			b.WriteString("  ")
			b.WriteString(style.Render(item.Label))
		}

		if i < len(items)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// LandingPage composites all layers into the final landing page
type LandingPage struct {
	grid         *WavyGrid
	lavaLamp     *LavaLamp
	width        int
	height       int
	frame        int
	selectedItem int
	menuItems    []string
}

// NewLandingPage creates a new landing page with all effects
func NewLandingPage(width, height int) *LandingPage {
	return &LandingPage{
		grid:     NewWavyGrid(width, height),
		lavaLamp: NewLavaLamp(width, height),
		width:    width,
		height:   height,
		menuItems: []string{
			"NEW GAME",
			"CONTINUE",
			"COLLECTION",
			"OPTIONS",
			"QUIT",
		},
		selectedItem: 0,
	}
}

// Update updates all animation components
func (lp *LandingPage) Update() {
	lp.frame++
	lp.grid.Update()
	lp.lavaLamp.Update()
}

// Resize updates dimensions for the landing page
func (lp *LandingPage) Resize(width, height int) {
	lp.width = width
	lp.height = height
	lp.grid.width = width
	lp.grid.height = height
	lp.lavaLamp.width = width
	lp.lavaLamp.height = height
}

// SelectNext moves selection to next menu item
func (lp *LandingPage) SelectNext() {
	lp.selectedItem = (lp.selectedItem + 1) % len(lp.menuItems)
}

// SelectPrev moves selection to previous menu item
func (lp *LandingPage) SelectPrev() {
	lp.selectedItem--
	if lp.selectedItem < 0 {
		lp.selectedItem = len(lp.menuItems) - 1
	}
}

// GetSelectedItem returns the currently selected menu item
func (lp *LandingPage) GetSelectedItem() string {
	return lp.menuItems[lp.selectedItem]
}

// RenderBackground composites the grid and lava lamp layers
func (lp *LandingPage) RenderBackground() []string {
	// Render both layers
	gridLines := strings.Split(lp.grid.Render(), "\n")
	blobLines := strings.Split(lp.lavaLamp.Render(), "\n")

	result := make([]string, lp.height)

	// Composite: blobs on top of grid
	for y := 0; y < lp.height && y < len(gridLines) && y < len(blobLines); y++ {
		var lineBuilder strings.Builder

		// This is a simplified compositing - in reality, we need to handle
		// styled strings more carefully. For now, prefer blob over grid.
		// If blob line has content (not just spaces), use it; otherwise use grid
		blobLine := blobLines[y]
		gridLine := gridLines[y]

		// Simple heuristic: if blob line has more than 50% spaces, show grid
		// Otherwise show blob
		if strings.Count(blobLine, " ") > len(blobLine)/2 {
			lineBuilder.WriteString(gridLine)
		} else {
			lineBuilder.WriteString(blobLine)
		}

		result[y] = lineBuilder.String()
	}

	return result
}

// Render generates the complete landing page with proper compositing
func (lp *LandingPage) Render() string {
	// Layer 1: Animated background (grid + blobs)
	background := lp.RenderBackground()

	// Layer 2: Title in a box (centered)
	title := renderTitle(lp.frame)
	titleBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neonColors.Magenta).
		Padding(1, 2).
		Render(title)

	// Layer 3: Menu in a box (centered below title)
	menuItemsList := make([]MenuItem, len(lp.menuItems))
	for i, label := range lp.menuItems {
		menuItemsList[i] = MenuItem{
			Label:    label,
			Selected: i == lp.selectedItem,
		}
	}
	menu := renderMenu(menuItemsList, lp.frame)
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neonColors.Gold).
		Padding(1, 2).
		Render(menu)

	// Use lipgloss.Place to overlay title and menu on background
	// First, join the background lines into a single string
	backgroundStr := strings.Join(background, "\n")

	// Calculate positions for centering
	titleHeight := lipgloss.Height(titleBox)
	menuHeight := lipgloss.Height(menuBox)
	totalContentHeight := titleHeight + 2 + menuHeight // 2 lines gap

	startY := (lp.height - totalContentHeight) / 2
	if startY < 0 {
		startY = 0
	}

	titleY := startY
	menuY := titleY + titleHeight + 2

	// Split both into lines for manual compositing
	bgLines := strings.Split(backgroundStr, "\n")
	titleLines := strings.Split(titleBox, "\n")
	menuLines := strings.Split(menuBox, "\n")

	// Ensure we have enough background lines
	for len(bgLines) < lp.height {
		bgLines = append(bgLines, strings.Repeat(" ", lp.width))
	}

	// Composite title onto background
	titleWidth := lipgloss.Width(titleLines[0])
	titleStartX := (lp.width - titleWidth) / 2
	if titleStartX < 0 {
		titleStartX = 0
	}

	result := make([]string, lp.height)
	copy(result, bgLines)

	// Overlay title
	for i, titleLine := range titleLines {
		y := titleY + i
		if y >= 0 && y < lp.height {
			result[y] = overlayString(result[y], titleLine, titleStartX, lp.width)
		}
	}

	// Overlay menu
	menuWidth := lipgloss.Width(menuLines[0])
	menuStartX := (lp.width - menuWidth) / 2
	if menuStartX < 0 {
		menuStartX = 0
	}

	for i, menuLine := range menuLines {
		y := menuY + i
		if y >= 0 && y < lp.height {
			result[y] = overlayString(result[y], menuLine, menuStartX, lp.width)
		}
	}

	return strings.Join(result, "\n")
}

// overlayString overlays src onto dst at position x, preserving background on left and right
// Uses visual width to handle ANSI-styled strings properly
func overlayString(dst, src string, x, maxWidth int) string {
	if x < 0 {
		return dst
	}

	// Get visual widths (ignoring ANSI codes)
	dstWidth := lipgloss.Width(dst)
	srcWidth := lipgloss.Width(src)

	// Build result by extracting visible portions
	var result strings.Builder

	// Left side: extract x visible characters from dst
	if x > 0 {
		leftPart := extractVisibleChars(dst, 0, x)
		result.WriteString(leftPart)
	}

	// Middle: the overlay
	result.WriteString(src)

	// Right side: extract remaining characters from dst after the overlay
	rightStart := x + srcWidth
	if rightStart < dstWidth {
		rightPart := extractVisibleChars(dst, rightStart, dstWidth-rightStart)
		result.WriteString(rightPart)
	}

	return result.String()
}

// extractVisibleChars extracts count visible characters starting from position start
// Handles ANSI escape codes properly
func extractVisibleChars(s string, start, count int) string {
	if count <= 0 {
		return ""
	}

	runes := []rune(s)
	visibleCount := 0
	inEscape := false
	startIdx := -1
	endIdx := -1

	for i, r := range runes {
		// Track ANSI escape sequences
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
			continue
		}

		// Only count visible characters
		if !inEscape {
			if visibleCount == start && startIdx == -1 {
				startIdx = i
			}
			if visibleCount >= start {
				visibleCount++
				if visibleCount-start >= count {
					endIdx = i + 1
					break
				}
			} else {
				visibleCount++
			}
		}
	}

	if startIdx == -1 {
		return ""
	}
	if endIdx == -1 {
		endIdx = len(runes)
	}

	// Need to include any ANSI codes before startIdx for proper coloring
	// Find the last complete ANSI sequence before startIdx
	result := string(runes[startIdx:endIdx])

	// Add ANSI reset at the end to prevent color bleed
	return result
}
