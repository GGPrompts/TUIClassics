# 2048 - Implementation Plan

**Type**: Tile Merging Puzzle
**Complexity**: Simple (~450 LOC)
**Time Estimate**: 3-4 hours
**Status**: Ready to implement

---

## 🎮 Game Overview

Slide tiles on a 4×4 grid to combine matching numbers. Each move spawns a new tile (2 or 4). Combine tiles to reach 2048!

### Core Mechanics:
- 4×4 grid of tiles
- Arrow keys slide all tiles in that direction
- Matching tiles merge (2+2=4, 4+4=8, etc.)
- New tile (2 or 4) spawns after each move
- Win at 2048, can continue to higher scores
- Game over when no moves possible

### Visual Layout:
```
┌────────────────────────────────┐
│  2048                          │
│  Score: 1234    Best: 5678     │
├────────────────────────────────┤
│                                │
│   ┌────┬────┬────┬────┐        │
│   │  2 │  4 │  8 │ 16 │        │
│   ├────┼────┼────┼────┤        │
│   │ 32 │ 64 │128 │256 │        │
│   ├────┼────┼────┼────┤        │
│   │512 │1024│2048│  4 │        │
│   ├────┼────┼────┼────┤        │
│   │  2 │  8 │ 16 │  2 │        │
│   └────┴────┴────┴────┘        │
│                                │
└────────────────────────────────┘
Press ↑↓←→ to move | R to restart
```

---

## 📁 File Structure (9-file modular pattern)

```
games/2048/
├── types.go           - Game state, tile, direction
├── model.go           - Model initialization
├── update.go          - Main update loop
├── update_keyboard.go - Arrow key handling
├── view.go            - Grid rendering
├── styles.go          - Tile colors by value
├── grid.go            - Grid manipulation logic
├── merge.go           - Tile merging algorithm
└── spawn.go           - New tile spawning
```

---

## 🎯 Implementation Phases

### Phase 1: Core Data Structures (~100 LOC)

**Files**: `types.go`, `model.go`

**types.go** - Define structures:
```go
type Direction int
const (
    Up Direction = iota
    Down
    Left
    Right
)

type GameState int
const (
    StateMenu GameState = iota
    StatePlaying
    StateWon
    StateGameOver
)

type Tile struct {
    Value  int
    Merged bool // Tracks if tile merged this turn
}

type Model struct {
    state     GameState
    grid      [4][4]Tile
    score     int
    bestScore int
    wonOnce   bool // True if reached 2048 this game

    // Terminal
    termWidth  int
    termHeight int
}
```

**model.go** - Initialization:
```go
func New() Model {
    return Model{
        state:     StateMenu,
        bestScore: 0,
    }
}

func (m *Model) startGame() {
    // Clear grid
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            m.grid[i][j] = Tile{Value: 0}
        }
    }

    m.score = 0
    m.wonOnce = false
    m.state = StatePlaying

    // Spawn two initial tiles
    m.spawnTile()
    m.spawnTile()
}
```

---

### Phase 2: Tile Spawning (~60 LOC)

**File**: `spawn.go`

```go
func (m *Model) spawnTile() {
    // Find empty cells
    var empty [][2]int
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            if m.grid[i][j].Value == 0 {
                empty = append(empty, [2]int{i, j})
            }
        }
    }

    if len(empty) == 0 {
        return // No empty cells
    }

    // Pick random empty cell
    pos := empty[rand.Intn(len(empty))]

    // 90% chance of 2, 10% chance of 4
    value := 2
    if rand.Float64() < 0.1 {
        value = 4
    }

    m.grid[pos[0]][pos[1]] = Tile{Value: value}
}

func (m *Model) hasEmptyCell() bool {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            if m.grid[i][j].Value == 0 {
                return true
            }
        }
    }
    return false
}
```

---

### Phase 3: Grid Movement (~150 LOC)

**File**: `grid.go`

**Core Algorithm**: Slide tiles in direction, merging as we go

```go
func (m *Model) move(dir Direction) bool {
    // Copy grid to detect if anything changed
    oldGrid := m.grid
    moved := false

    // Clear merged flags
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            m.grid[i][j].Merged = false
        }
    }

    switch dir {
    case Left:
        moved = m.moveLeft()
    case Right:
        moved = m.moveRight()
    case Up:
        moved = m.moveUp()
    case Down:
        moved = m.moveDown()
    }

    // If grid changed, spawn new tile
    if moved {
        m.spawnTile()

        // Check win condition
        if !m.wonOnce && m.hasValue(2048) {
            m.state = StateWon
            m.wonOnce = true
        }

        // Check game over
        if !m.canMove() {
            m.state = StateGameOver
            if m.score > m.bestScore {
                m.bestScore = m.score
            }
        }
    }

    return moved
}

func (m *Model) moveLeft() bool {
    moved := false
    for i := 0; i < 4; i++ {
        if m.slideAndMergeRow(i) {
            moved = true
        }
    }
    return moved
}

func (m *Model) moveRight() bool {
    // Reverse each row, move left, reverse back
    m.reverseRows()
    moved := m.moveLeft()
    m.reverseRows()
    return moved
}

func (m *Model) moveUp() bool {
    // Transpose, move left, transpose back
    m.transpose()
    moved := m.moveLeft()
    m.transpose()
    return moved
}

func (m *Model) moveDown() bool {
    // Transpose, reverse, move left, reverse, transpose
    m.transpose()
    m.reverseRows()
    moved := m.moveLeft()
    m.reverseRows()
    m.transpose()
    return moved
}
```

---

### Phase 4: Merge Logic (~100 LOC)

**File**: `merge.go`

**Core merge algorithm for a single row**:

```go
func (m *Model) slideAndMergeRow(row int) bool {
    // Extract non-zero values
    values := []int{}
    for j := 0; j < 4; j++ {
        if m.grid[row][j].Value != 0 {
            values = append(values, m.grid[row][j].Value)
        }
    }

    // Merge adjacent equal values
    merged := []int{}
    i := 0
    for i < len(values) {
        if i+1 < len(values) && values[i] == values[i+1] {
            // Merge!
            newValue := values[i] * 2
            merged = append(merged, newValue)
            m.score += newValue
            i += 2
        } else {
            merged = append(merged, values[i])
            i++
        }
    }

    // Check if row changed
    changed := false
    for j := 0; j < 4; j++ {
        newValue := 0
        if j < len(merged) {
            newValue = merged[j]
        }

        if m.grid[row][j].Value != newValue {
            changed = true
        }

        m.grid[row][j] = Tile{Value: newValue}
    }

    return changed
}

// Helper functions
func (m *Model) transpose() {
    for i := 0; i < 4; i++ {
        for j := i + 1; j < 4; j++ {
            m.grid[i][j], m.grid[j][i] = m.grid[j][i], m.grid[i][j]
        }
    }
}

func (m *Model) reverseRows() {
    for i := 0; i < 4; i++ {
        for j := 0; j < 2; j++ {
            m.grid[i][j], m.grid[i][3-j] = m.grid[i][3-j], m.grid[i][j]
        }
    }
}
```

---

### Phase 5: Game Over Detection (~50 LOC)

**File**: `grid.go` (continued)

```go
func (m *Model) canMove() bool {
    // Check for empty cells
    if m.hasEmptyCell() {
        return true
    }

    // Check for possible merges (horizontal)
    for i := 0; i < 4; i++ {
        for j := 0; j < 3; j++ {
            if m.grid[i][j].Value == m.grid[i][j+1].Value {
                return true
            }
        }
    }

    // Check for possible merges (vertical)
    for i := 0; i < 3; i++ {
        for j := 0; j < 4; j++ {
            if m.grid[i][j].Value == m.grid[i+1][j].Value {
                return true
            }
        }
    }

    return false
}

func (m *Model) hasValue(target int) bool {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            if m.grid[i][j].Value == target {
                return true
            }
        }
    }
    return false
}
```

---

### Phase 6: Input Handling (~40 LOC)

**File**: `update_keyboard.go`

```go
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateMenu:
        if msg.String() == "enter" {
            m.startGame()
        }

    case StatePlaying:
        switch msg.String() {
        case "up", "k", "w":
            m.move(Up)
        case "down", "j", "s":
            m.move(Down)
        case "left", "h", "a":
            m.move(Left)
        case "right", "l", "d":
            m.move(Right)
        case "r":
            m.startGame()
        }

    case StateWon:
        switch msg.String() {
        case "c":
            m.state = StatePlaying // Continue playing
        case "r":
            m.startGame()
        }

    case StateGameOver:
        if msg.String() == "r" {
            m.startGame()
        }
    }

    return m, nil
}
```

---

### Phase 7: Rendering (~100 LOC)

**File**: `view.go`

```go
func (m Model) View() string {
    switch m.state {
    case StateMenu:
        return m.renderMenu()
    case StatePlaying:
        return m.renderGame()
    case StateWon:
        return m.renderWin()
    case StateGameOver:
        return m.renderGameOver()
    }
    return ""
}

func (m Model) renderGame() string {
    var b strings.Builder

    // Title and score
    b.WriteString(titleStyle.Render("2048") + "\n")
    b.WriteString(fmt.Sprintf("Score: %d    Best: %d\n\n", m.score, m.bestScore))

    // Grid
    b.WriteString("┌────┬────┬────┬────┐\n")

    for i := 0; i < 4; i++ {
        b.WriteString("│")
        for j := 0; j < 4; j++ {
            b.WriteString(m.renderTile(m.grid[i][j]))
            b.WriteString("│")
        }
        b.WriteString("\n")

        if i < 3 {
            b.WriteString("├────┼────┼────┼────┤\n")
        }
    }

    b.WriteString("└────┴────┴────┴────┘\n\n")
    b.WriteString("Press ↑↓←→ to move | R to restart | Q to quit\n")

    return lipgloss.Place(m.termWidth, m.termHeight,
        lipgloss.Center, lipgloss.Center, b.String())
}

func (m Model) renderTile(tile Tile) string {
    if tile.Value == 0 {
        return "    "
    }

    // Get style based on value
    style := getTileStyle(tile.Value)
    text := fmt.Sprintf("%4d", tile.Value)

    return style.Render(text)
}
```

---

## 🎨 Visual Design

**File**: `styles.go`

**Colors by tile value** (like original 2048):

```go
func getTileStyle(value int) lipgloss.Style {
    style := lipgloss.NewStyle().Bold(true)

    switch value {
    case 2:
        return style.Foreground(lipgloss.Color("238")) // Dark gray
    case 4:
        return style.Foreground(lipgloss.Color("239"))
    case 8:
        return style.Foreground(lipgloss.Color("208")) // Orange
    case 16:
        return style.Foreground(lipgloss.Color("202"))
    case 32:
        return style.Foreground(lipgloss.Color("196")) // Red
    case 64:
        return style.Foreground(lipgloss.Color("160"))
    case 128:
        return style.Foreground(lipgloss.Color("226")) // Yellow
    case 256:
        return style.Foreground(lipgloss.Color("220"))
    case 512:
        return style.Foreground(lipgloss.Color("214"))
    case 1024:
        return style.Foreground(lipgloss.Color("208"))
    case 2048:
        return style.Foreground(lipgloss.Color("46")) // Green (victory!)
    default:
        return style.Foreground(lipgloss.Color("201")) // Magenta (>2048)
    }
}
```

---

## ✅ Implementation Checklist

### Part 1: Setup
- [ ] Create all 9 files in `games/2048/`
- [ ] Define types (Direction, Tile, Model)
- [ ] Implement `New()` and `startGame()`

### Part 2: Tile Spawning
- [ ] Implement `spawnTile()` in `spawn.go`
- [ ] Spawn two tiles on game start
- [ ] Test: See tiles appear on grid

### Part 3: Core Movement
- [ ] Implement `slideAndMergeRow()` in `merge.go`
- [ ] Handle merging logic (2+2=4)
- [ ] Test: Single row slides and merges correctly

### Part 4: Directional Movement
- [ ] Implement `moveLeft()` in `grid.go`
- [ ] Implement `moveRight()` (reverse → left → reverse)
- [ ] Implement `moveUp()` (transpose → left → transpose)
- [ ] Implement `moveDown()` (transpose → reverse → left → reverse → transpose)
- [ ] Test: All four directions work

### Part 5: Game Over
- [ ] Implement `canMove()` in `grid.go`
- [ ] Check for empty cells
- [ ] Check for possible merges
- [ ] Test: Game over when stuck

### Part 6: Input
- [ ] Handle arrow keys in `update_keyboard.go`
- [ ] Call `move()` with appropriate direction
- [ ] Test: Arrow keys control game

### Part 7: Rendering
- [ ] Render grid in `view.go`
- [ ] Color tiles by value in `styles.go`
- [ ] Show score
- [ ] Test: Grid looks good

### Part 8: Polish
- [ ] Win screen at 2048
- [ ] Allow continuing after win
- [ ] Game over screen
- [ ] Save best score
- [ ] Restart functionality

---

## 🎯 Success Criteria

When complete, you should be able to:
1. ✅ Move tiles with arrow keys
2. ✅ Merge matching tiles
3. ✅ Spawn new tiles after each move
4. ✅ See score increase on merges
5. ✅ Win when reaching 2048
6. ✅ Game over when no moves left
7. ✅ Restart game

---

## 🚀 Quick Start

```bash
cd ~/projects/TUIClassics/games/2048

# Create files
touch types.go model.go update.go update_keyboard.go view.go styles.go grid.go merge.go spawn.go

# Implementation order:
# 1. types.go - Define all types
# 2. model.go - New() and startGame()
# 3. spawn.go - spawnTile()
# 4. merge.go - slideAndMergeRow()
# 5. grid.go - move(), transpose(), reverse()
# 6. update_keyboard.go - Arrow keys
# 7. view.go - Render grid
# 8. styles.go - Tile colors
# 9. update.go - Game loop

# Register in menu (games/menu/model.go)

# Build
cd ../..
make classics
./bin/classics
```

---

## 💡 Pro Tips

1. **Test merge first** - Get single row working before doing all directions
2. **Use transformations** - Transpose + reverse = all 4 directions with one algorithm
3. **Check for changes** - Only spawn tile if grid actually moved
4. **Visual feedback** - Different colors help see progress
5. **Score on merge** - Add merged value to score (not just 1 point)

---

## 🔜 Future Enhancements

- **Undo move** - Save grid history
- **Animations** - Slide tiles smoothly
- **Different grid sizes** - 3×3, 5×5, 6×6
- **Game modes** - Reach 4096, 8192, etc.
- **Statistics** - Track moves, time, merge counts
- **Save/load** - Continue later

---

**Ready to merge! 🎲 Let's build 2048!**
