# Claude Code Development Notes

Documentation of key patterns, fixes, and architectural decisions for AI-assisted development of TUI Classics.

## Project Architecture

This project follows the modular architecture pattern from [TUITemplate](https://github.com/GGPrompts/TUITemplate).

### File Structure (Per Game)

Each game follows the 9-file modular pattern:

```
games/minesweeper/
├── main.go           (cmd/minesweeper/main.go - 21 lines max, just launch)
├── types.go          (All type definitions, structs, enums, constants)
├── model.go          (Model initialization, game logic)
├── update.go         (Main message dispatcher)
├── update_keyboard.go (Keyboard event handling)
├── update_mouse.go   (Mouse event handling)
├── view.go           (View rendering and layout)
├── styles.go         (Lipgloss style definitions)
└── [game-specific]   (e.g., grid.go, cards.go)
```

**Key Principle**: Clean separation of concerns - each file has a single responsibility.

## Critical Fixes & Patterns

### Issue 1: Mouse Coordinate Calculation

**Symptom**: Mouse clicks registered 2 rows above where user clicked.

**Root Cause**: Incorrect calculation of grid starting Y position.

**The Math Error**:
```go
// WRONG (what we had initially):
gridStartY := topPadding + 1 + 2 + 1 + 2  // = topPadding + 6

// CORRECT:
gridStartY := topPadding + 4
```

**Why**: When rendering `content + "\n\n"`, the content occupies the current line, then `\n\n` moves cursor forward 2 lines. The next content is written where cursor lands, not 2 additional lines later.

**Trace Through Rendering**:
```
Line topPadding:     title written here
                     \n\n moves to topPadding+2
Line topPadding+2:   stats written here
                     \n\n moves to topPadding+4
Line topPadding+4:   grid starts here
```

**Fix** (minesweeper/update_mouse.go:66):
```go
// Grid starts after: topPadding + title line + "\n\n" (2) + stats line + "\n\n" (2)
gridStartY := topPadding + 4
```

**Key Lesson**: Mouse coordinate calculations MUST exactly match the rendering logic. Keep both in sync.

---

### Issue 2: Grid Misalignment from Variable-Width Emojis

**Symptom**: When mines were revealed, the grid columns became misaligned.

**Root Cause**: Emojis have inconsistent terminal widths:
- 💣 (bomb) = 2 characters wide
- 🚩 (flag) = 2 characters wide
- Numbers = 1 character wide

**Fix**: Replace all emojis with consistent single-width characters:

```go
// BEFORE (minesweeper/view.go):
if cell.IsFlagged {
    symbol = "🚩"  // 2 chars wide!
} else if cell.IsMine {
    symbol = "💣"  // 2 chars wide!
}

// AFTER:
if cell.IsFlagged {
    symbol = "P"   // 1 char wide ✓
} else if cell.IsMine {
    symbol = "*"   // 1 char wide ✓
}
```

**Character Choices**:
- Unrevealed cell: `■` (solid block)
- Revealed empty: `·` (middle dot)
- Mine: `*` (asterisk)
- Flag: `P` (for "possible mine")
- Numbers: `1-8` (classic colors preserved)

**Key Lesson**: In terminal UIs, always use single-width characters for grid-based layouts. Unicode box-drawing characters (─│┌┐└┘) are safe; emojis are not.

---

### Issue 3: Centering with `lipgloss.Place` Breaks Coordinate Math

**Symptom**: Mouse clicks didn't work when using `lipgloss.Place` for centering.

**Root Cause**: `lipgloss.Place` centers content visually but makes coordinate calculation unpredictable.

**Fix**: Manual centering with calculated padding:

```go
// DON'T use Place for interactive grids:
return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, content)

// DO use manual padding:
topPadding := (m.termHeight - totalLines) / 2
leftPadding := (m.termWidth - gridWidth) / 2

// Then use same calculations in mouse handler
gridStartY := topPadding + 4
gridStartX := leftPadding
```

**Key Lesson**: For interactive elements (clickable grids), avoid `lipgloss.Place`. Use manual centering so mouse coordinate math is deterministic.

---

## Animation Patterns

### Explosion Animation (Minesweeper)

**Implementation**: Expanding shockwave using Manhattan distance.

**Key Components**:

1. **Animation State** (types.go:68-72):
```go
type Model struct {
    // ...
    explosionCenterX  int // Where mine was clicked
    explosionCenterY  int
    explosionRadius   int // Current wave radius
    explosionMaxSteps int // Total frames
}
```

2. **Animation Messages** (types.go:94-108):
```go
type AnimationTickMsg time.Time

func animationTick() tea.Cmd {
    return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
        return AnimationTickMsg(t)
    })
}
```

3. **Frame Updates** (model.go:204-228):
```go
func (m *Model) progressExplosion() {
    m.explosionRadius++

    // Reveal mines within current radius
    for y := 0; y < m.height; y++ {
        for x := 0; x < m.width; x++ {
            if m.grid[y][x].IsMine {
                dist := abs(x-m.explosionCenterX) + abs(y-m.explosionCenterY)
                if dist <= m.explosionRadius {
                    m.grid[y][x].IsRevealed = true
                }
            }
        }
    }

    // End when complete
    if m.explosionRadius >= m.explosionMaxSteps {
        m.state = StateLost
    }
}
```

4. **Visual Effects** (view.go:186-189):
```go
// Mines at current radius blink
if dist == m.explosionRadius && cell.IsMine {
    symbol = "+"
    style = mineCellStyle.Copy().Blink(true)
}
```

**Key Lessons**:
- Use separate message type (`AnimationTickMsg`) for animation frames
- Calculate max steps upfront based on grid size
- Use Manhattan distance for radial effects in grid layouts
- Progressive reveal creates satisfying visual feedback

---

## Shared Code Patterns

### Cell Rendering Helper

Refactored cell rendering into reusable helper to avoid duplication:

```go
// view.go:208-235
func (m Model) getCellSymbolAndStyle(cell Cell, isCursor bool) (string, lipgloss.Style) {
    var symbol string
    var style lipgloss.Style

    if cell.IsFlagged && !cell.IsRevealed {
        symbol = "P"
        style = flagCellStyle
    } else if !cell.IsRevealed {
        symbol = "■"
        style = unknownCellStyle
    } else if cell.IsMine {
        symbol = "*"
        style = mineCellStyle
    } else if cell.Adjacent == 0 {
        symbol = "·"
        style = revealedCellStyle
    } else {
        symbol = fmt.Sprintf("%d", cell.Adjacent)
        style = getNumberStyle(cell.Adjacent)
    }

    if isCursor && !cell.IsRevealed {
        style = cursorCellStyle
    }

    return symbol, style
}
```

**Usage**: Both `renderCell()` and `renderExplosionGrid()` use this helper, ensuring consistent rendering.

---

## Color Palette (Classic Minesweeper)

Following Windows 95 Minesweeper exactly:

```go
// styles.go:12-19
var (
    color1 = lipgloss.Color("#0000FF") // Blue
    color2 = lipgloss.Color("#008000") // Green
    color3 = lipgloss.Color("#FF0000") // Red
    color4 = lipgloss.Color("#000080") // Dark Blue
    color5 = lipgloss.Color("#800000") // Maroon
    color6 = lipgloss.Color("#008080") // Teal
    color7 = lipgloss.Color("#000000") // Black
    color8 = lipgloss.Color("#808080") // Gray
)
```

**Why These Colors**: Nostalgia and readability. These exact colors are instantly recognizable to anyone who played Windows 95 Minesweeper.

---

## Testing Checklist

When implementing new games or features:

1. **Grid Alignment**:
   - [ ] All characters are single-width
   - [ ] Grid stays aligned when all cell states revealed
   - [ ] Test on different terminal sizes

2. **Mouse Coordinates**:
   - [ ] Click detection matches visual position
   - [ ] Test all four corners of grid
   - [ ] Test after terminal resize

3. **Keyboard Navigation**:
   - [ ] Arrow keys move cursor correctly
   - [ ] Enter/Space triggers correct cell
   - [ ] Vim keys (h/j/k/l) work if implemented

4. **State Transitions**:
   - [ ] Menu → Game → Win/Loss flow works
   - [ ] New game resets all state
   - [ ] Restart preserves difficulty

5. **Edge Cases**:
   - [ ] Very small terminal (< 80 cols)
   - [ ] Very large grids (30x16)
   - [ ] Rapid clicking doesn't break state

---

## Dependencies

From [TUITemplate](https://github.com/GGPrompts/TUITemplate):
- `github.com/charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `github.com/charmbracelet/lipgloss` - Terminal styling

**Version Compatibility**: Using `@latest` for both. Bubbletea v1.x is stable.

---

## Future Games: Lessons to Apply

### For Solitaire

**Mouse Handling**: Will need drag-and-drop:
```go
type Model struct {
    draggingCard *Card
    dragStartX   int
    dragStartY   int
}

// In update_mouse.go:
case tea.MouseButtonLeft:
    if msg.Action == tea.MouseActionPress {
        m.draggingCard = m.getCardAt(msg.X, msg.Y)
    } else if msg.Action == tea.MouseActionRelease {
        m.dropCard(msg.X, msg.Y)
        m.draggingCard = nil
    }
```

**Card Rendering**: Use single-width Unicode card symbols:
```go
// ♠ ♥ ♦ ♣ for suits
// A 2 3 4 5 6 7 8 9 T J Q K for ranks
// Or use full card Unicode: 🂡🂢🂣 etc. IF terminal width is consistent
```

**Waterfall Animation**: Similar to explosion, but:
- Cards fall from top (foundation piles)
- Use Y velocity + gravity simulation
- Bounce at bottom with decreasing elasticity

### For Snake

**Movement**: Simple tick-based animation:
```go
type TickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}
```

**Collision Detection**: Check if head position == body segment position.

---

## Solitaire Implementation Notes

### Issue 4: Centering with Manual Padding

**Symptom**: Mouse coordinates don't work with `lipgloss.Place` centering.

**Solution**: Use manual padding with exact coordinate calculations:

```go
// Calculate padding
topPadding := (m.termHeight - totalLines) / 2
leftPadding := (m.termWidth - contentWidth) / 2

// Apply padding to EACH line
for _, line := range contentLines {
    b.WriteString(strings.Repeat(" ", leftPadding))
    b.WriteString(line)
    b.WriteString("\n")
}
```

**Mouse Coordinate Adjustment**:
```go
// Subtract padding to get content-relative coordinates
relX := mouseX - leftPadding
relY := mouseY - topPadding
```

**Key Lesson**: Multi-line content (like cards) needs per-line padding, not just first line.

---

### Issue 5: Double Spacing from Margins

**Symptom**: Cards misaligned after terminal resize; top row broken.

**Root Cause**: Style margins + manual `\n\n` created double spacing:
```go
// WRONG:
titleStyle = lipgloss.NewStyle().MarginBottom(1)
b.WriteString(titleStyle.Render(title))
b.WriteString("\n\n")  // 1 + 2 = 3 lines!

// CORRECT:
titleStyle = lipgloss.NewStyle()  // No margin
b.WriteString(titleStyle.Render(title))
b.WriteString("\n\n")  // Exactly 2 lines
```

**Key Lesson**: Don't mix lipgloss margins with manual spacing - pick one approach.

---

### Issue 6: Waterfall Animation with Styled Text

**Symptom**: ANSI escape codes visible in grid; white backgrounds don't connect.

**Root Cause**: Lipgloss styles add ANSI codes that break plain character grids.

**Solution**: Render plain-text cards without styling:
```go
// Plain box-drawing characters only
╭─────╮
│ A S │  // Use S/H/D/C instead of ♠♥♦♣ (double-width)
│     │
│ S A │
╰─────╯
```

**Key Lesson**: Character grids can't handle ANSI styling - use plain text for animations.

---

### Issue 7: Mouse Clicks Off by Multiple Lines (totalLines Calculation)

**Symptom**: Mouse clicks register 1-2 lines higher than actual click position. Clicks don't register on top 1/4 of cards.

**Root Cause**: Misunderstanding how `"\n\n"` spacing works in line count calculations.

**The Math Error**:
```go
// WRONG (what we had):
totalLines := 1 + 2 + 1 + 2 + 5 + 2 + maxTableauHeight + 1 + 1
//           title  "\n\n" stats "\n\n" topRow "\n\n" ...
//           = 15 + maxTableauHeight (4 lines too many!)

// CORRECT (per CLAUDE.md Minesweeper fix):
totalLines := 11 + maxTableauHeight
```

**Why**: When you write `"title\n\n"`, it doesn't occupy 1 + 2 = 3 lines. It occupies **2 lines total**:
- Line 0: title (content written)
- Line 1: blank (from first `\n`)
- Cursor moves to line 2 (from second `\n`)

**Trace Through Rendering**:
```
Line 0: title
Line 1: blank (from "\n\n")
Line 2: stats
Line 3: blank (from "\n\n")
Lines 4-8: topRow (5 lines)
Line 9: blank (from "\n\n")
Lines 10+: tableau starts
```

**Fix** (view.go:75, update_mouse.go:209, update_mouse.go:293):
```go
// Each "content\n\n" = 2 lines (content + 1 blank)
// title(1) + blank(1) + stats(1) + blank(1) + topRow(5) + blank(1) + tableau + help(1)
totalLines := 11 + maxTableauHeight
```

**Impact**: Old calculation made `topPadding` ~2 lines too small, causing content to render higher than expected. When users clicked on visually-rendered cards, the coordinate math thought clicks were 2 lines above the calculated card positions, preventing clicks on the top portion of cards.

**Key Lesson**: This is the **same pattern** as Minesweeper Issue 1. The `"\n\n"` spacing adds 1 blank line to the count, not 2. **Always trace through actual rendering line-by-line** to verify coordinate calculations.

---

### Pattern: Click-to-Select for Hybrid Input

**Implementation**: Detect click vs drag by measuring distance:

```go
// On press: Store position
m.mousePressX = msg.X
m.mousePressY = msg.Y

// On release: Check if moved
dx := msg.X - m.mousePressX
dy := msg.Y - m.mousePressY
distanceMoved := dx*dx + dy*dy

if distanceMoved < 4 {  // Threshold: ~2 pixels
    // Click: Select for keyboard movement
    m.selectedPile = m.dragFromPile
    m.cursor = *m.dragFromPile
} else {
    // Drag: Move cards immediately
    m.MoveCards(...)
}
```

**Benefits**: Users can click to select, then use arrow keys to move - best of both worlds.

---

### Pattern: Initial Terminal Size Request

**Problem**: Games appear uncentered on first launch from launcher.

**Solution**: Request terminal dimensions in `Init()`:
```go
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        tea.EnterAltScreen,
        tea.WindowSize(),  // Request size immediately
        tickCmd(),
    )
}
```

**Key Lesson**: Don't wait for manual resize - request dimensions on init.

---

## Key Files Reference

| Feature | File | Lines |
|---------|------|-------|
| Mouse coordinate fix | `minesweeper/update_mouse.go` | 66 |
| Explosion animation | `minesweeper/model.go` | 204-236 |
| Animation rendering | `minesweeper/view.go` | 116-206 |
| Cell rendering helper | `minesweeper/view.go` | 208-235 |
| Classic color palette | `minesweeper/styles.go` | 12-43 |
| Grid alignment fix | `minesweeper/view.go` | 256-260 |

---

**Created**: 2025-01-20
**Last Updated**: 2025-01-21
**Purpose**: Preserve critical patterns and bug fixes for future AI-assisted development

**Related Documentation**:
- [TUITemplate CLAUDE.md](https://github.com/GGPrompts/TUITemplate/blob/main/CLAUDE.md) - Original layout patterns
- [TUITemplate ARCHITECTURE.md](https://github.com/GGPrompts/TUITemplate/blob/main/ARCHITECTURE.md) - Modular structure guide
