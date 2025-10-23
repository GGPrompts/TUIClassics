# Snake - Implementation Summary

**Type**: Classic Movement Game
**Complexity**: Medium (~550 LOC)
**Time Estimate**: 3-4 hours
**Status**: ✅ COMPLETED

---

## 🐍 Game Overview

Classic Snake game with personality! Control a growing emoji snake, eat apples, and avoid hitting walls or yourself.

### Core Mechanics:
- Snake moves continuously with tick-based animation
- Arrow keys (or WASD/HJKL) change direction
- Eating apples (🍎) makes snake grow
- Game over if you hit walls or yourself
- Score increases with each apple eaten
- Progressive speed increase based on difficulty

### Visual Layout (Actual Implementation):
```
╭──────────────────────────────────────────────╮
│         Score: 15    High: 42                │
│                                              │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ │  ← Visible walls
│ ▓▓                                        ▓▓ │
│ ▓▓      😊●●●●                            ▓▓ │  ← Emoji head!
│ ▓▓        ●                               ▓▓ │
│ ▓▓        ●●                              ▓▓ │
│ ▓▓                                        ▓▓ │
│ ▓▓              🍎                        ▓▓ │  ← Apple food
│ ▓▓                                        ▓▓ │
│ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ │
╰──────────────────────────────────────────────╯

Press arrow keys to move | P to pause | Q to quit
```

---

## 🎨 Key Features Implemented

### ✨ Emoji Snake with Expressions
- **😊 Happy Face** - Normal state (bright green)
- **😮 Surprise Face** - When eating (yellow flash)
- **😵 Dizzy Face** - Crash animation (red, 1-second pause)

### 🎯 Difficulty Selection
Three balanced difficulty levels with different speed curves:

| Difficulty | Start Speed | Max Speed | Progression | Scores to Max |
|-----------|-------------|-----------|-------------|---------------|
| 🟢 Easy    | 180ms (5.6/s) | 100ms (10/s) | -2ms/score | 40 |
| 🟡 Medium  | 200ms (5/s)   | 80ms (12.5/s) | -4ms/score | 30 |
| 🔴 Hard    | 150ms (6.7/s) | 50ms (20/s!) | -5ms/score | 20 |

### 🧱 Visible Walls
- Dark gray block walls (▓▓) around playable area
- Prevents "hidden boundary" confusion
- Crash faces always visible against walls
- Retro aesthetic

### 🎮 Modern UI/UX
- Lipgloss `RoundedBorder()` for polished look
- Difficulty selection screen with emoji indicators
- Centered layouts with proper aspect ratio (24×20 grid)
- Crash animation with 1-second pause before game over

---

## 📁 File Structure (9-file modular pattern)

```
games/snake/
├── types.go           - Game state, snake, direction, difficulty enums
├── model.go           - Model initialization, difficulty settings
├── update.go          - Main update loop, game tick, crash delay
├── update_keyboard.go - Direction control, difficulty selection
├── view.go            - Rendering (menu, difficulty, game, game over)
├── styles.go          - Visual styles (emojis, walls, UI)
├── snake.go           - Snake logic (movement, growth, pre-collision)
├── collision.go       - Collision detection (wouldCollide, checkCollision)
└── food.go            - Food spawning
```

**Entry Point:**
```
cmd/snake/main.go      - Standalone game launcher
```

---

## 🎯 Implementation Details

### Phase 1: Core Data Structures ✅

**types.go** - Extended from plan:
```go
type Difficulty int
const (
    Easy Difficulty = iota
    Medium
    Hard
)

type GameState int
const (
    StateMenu GameState = iota
    StateDifficultySelect  // NEW!
    StatePlaying
    StatePaused
    StateCrashed          // NEW! 1-second crash animation
    StateGameOver
)

type Model struct {
    state              GameState
    difficulty         Difficulty
    selectedDifficulty Difficulty // For menu cursor
    justAte            bool       // NEW! Eating animation flag
    // ... rest of fields
}
```

### Phase 2: Difficulty System ✅

**model.go** - getDifficultySettings():
```go
func (m *Model) getDifficultySettings() (initialSpeed, minSpeed time.Duration, speedDecrease int) {
    switch m.difficulty {
    case Easy:
        return 180 * time.Millisecond, 100 * time.Millisecond, 2
    case Medium:
        return 200 * time.Millisecond, 80 * time.Millisecond, 4
    case Hard:
        return 150 * time.Millisecond, 50 * time.Millisecond, 5
    }
}
```

### Phase 3: Pre-Collision Detection ✅

**Critical fix**: Check collision BEFORE moving to show crash face

**snake.go**:
```go
func (m *Model) moveSnake() bool {
    // Calculate new head position
    newHead := calculateNewPosition()

    // Check collision BEFORE moving (shows crash face!)
    if m.wouldCollide(newHead) {
        return false // Don't move
    }

    // Move snake...
    return true
}
```

**collision.go**:
```go
func (m *Model) wouldCollide(pos Point) bool {
    // Wall collision
    if pos.X < 0 || pos.X >= m.width ||
       pos.Y < 0 || pos.Y >= m.height {
        return true
    }

    // Self collision
    for i := 0; i < len(m.snake)-1; i++ {
        if pos == m.snake[i] {
            return true
        }
    }
    return false
}
```

### Phase 4: Crash Animation ✅

**update.go**:
```go
type CrashDelayMsg time.Time

func crashDelayCmd() tea.Cmd {
    return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
        return CrashDelayMsg(t)
    })
}

case TickMsg:
    if m.state == StatePlaying {
        if !m.moveSnake() {
            m.crash()
            return m, crashDelayCmd() // 1-second pause
        }
    }

case CrashDelayMsg:
    if m.state == StateCrashed {
        m.gameOver()
    }
```

### Phase 5: Visual Rendering ✅

**view.go** - Emoji rendering with state-based expressions:
```go
// Game board with visible walls
for y := -1; y <= m.height; y++ {
    for x := -1; x <= m.width; x++ {
        point := Point{x, y}
        isWall := x == -1 || x == m.width || y == -1 || y == m.height

        if isWall {
            b.WriteString(wallStyle.Render("▓▓"))
        } else if point == m.snake[0] {
            if m.state == StateCrashed {
                b.WriteString(headCrashedStyle.Render("😵"))
            } else if m.justAte {
                b.WriteString(headEatingStyle.Render("😮"))
            } else {
                b.WriteString(headStyle.Render("😊"))
            }
        } else if m.isSnakeBody(point) {
            b.WriteString(bodyStyle.Render("●●"))
        } else if point == m.food {
            b.WriteString(foodStyle.Render("🍎"))
        } else {
            b.WriteString("  ")
        }
    }
}

// Wrap in Lipgloss rounded border
bordered := gameBorderStyle.Render(b.String())
```

### Phase 6: Difficulty Selection Screen ✅

**view.go** - renderDifficultySelect():
```go
difficulties := []struct{
    diff        Difficulty
    name        string
    description string
}{
    {Easy, "🟢 EASY", "Gentle pace - 180ms start, 100ms max"},
    {Medium, "🟡 MEDIUM", "Balanced challenge - 200ms start, 80ms max"},
    {Hard, "🔴 HARD", "Fast & intense - 150ms start, 50ms max"},
}

for _, d := range difficulties {
    if d.diff == m.selectedDifficulty {
        b.WriteString(selectedDifficultyStyle.Render("► " + d.name))
    } else {
        b.WriteString(difficultyStyle.Render("  " + d.name))
    }
}
```

---

## 🎨 Visual Design (As Implemented)

**styles.go**:
```go
// Snake emoji styles
var (
    headStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("46")).  // Bright green
        Bold(true)

    headEatingStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("226")). // Yellow flash
        Bold(true)

    headCrashedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")). // Red
        Bold(true)

    bodyStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("46"))   // Bright green (matches head)

    foodStyle = lipgloss.NewStyle().
        Bold(true)

    wallStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")) // Dark gray
)

// Border with no padding (walls handle spacing)
gameBorderStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("46"))
```

---

## ✅ Implementation Checklist

### Part 1: Setup ✅
- [x] Create all 9 files in `games/snake/`
- [x] Define types (Direction, Point, GameState, Model, Difficulty)
- [x] Implement `New()` and `startGame()`

### Part 2: Core Movement ✅
- [x] Implement `moveSnake()` in `snake.go`
- [x] Handle direction changes with buffer
- [x] Prevent 180-degree turns
- [x] Pre-collision detection

### Part 3: Food System ✅
- [x] Implement `spawnFood()` in `food.go`
- [x] Check if snake ate food
- [x] Make snake grow when eating
- [x] Eating animation (😮 surprise face)

### Part 4: Collision ✅
- [x] Implement `wouldCollide()` for pre-check
- [x] Implement `checkCollision()` in `collision.go`
- [x] Detect wall hits
- [x] Detect self-collision
- [x] Crash animation (😵 dizzy face, 1-second pause)

### Part 5: Input ✅
- [x] Handle arrow keys in `update_keyboard.go`
- [x] Support WASD and HJKL alternatives
- [x] Difficulty selection navigation (↑↓)
- [x] Buffer direction changes

### Part 6: Rendering ✅
- [x] Render game board in `view.go`
- [x] Draw emoji snake with expressions
- [x] Draw apple food (🍎)
- [x] Draw visible walls (▓▓)
- [x] Lipgloss RoundedBorder
- [x] Show score and high score
- [x] Difficulty selection screen

### Part 7: Polish ✅
- [x] Add pause functionality
- [x] Show game over screen with restart option
- [x] Increase speed progressively based on difficulty
- [x] Save high score in memory
- [x] Difficulty selection system
- [x] Crash animation with delay
- [x] Eating animation with visual feedback

---

## 🎯 Success Criteria - ALL MET! ✅

1. ✅ Control snake with arrow keys (+ WASD/HJKL)
2. ✅ Eat food to grow and score
3. ✅ Game over on wall/self collision with crash animation
4. ✅ See score and high score
5. ✅ Pause and resume
6. ✅ Restart after game over (keeps difficulty)
7. ✅ **BONUS**: Emoji snake with expressions
8. ✅ **BONUS**: Three difficulty levels
9. ✅ **BONUS**: Visible walls for clear boundaries
10. ✅ **BONUS**: Crash face always visible

---

## 🚀 Build & Run

```bash
# Build from snake worktree
cd ~/projects/TUIClassics-snake
go build -o bin/snake ./cmd/snake
./bin/snake

# Or from main project
cd ~/projects/TUIClassics
make classics
./bin/classics  # Select Snake from menu
```

**Registered in menu**: Hotkey **N** (for s**N**ake)

---

## 💡 Key Implementation Lessons

1. **Aspect Ratio Matters**: Changed from 30×20 to 24×20 grid for better visual balance
2. **Pre-Collision Detection**: Check collision BEFORE moving to show crash face
3. **Emoji Width**: Emojis are 2-char wide, use `●●` for body to match
4. **Border Padding Issue**: Removed padding, added visible walls instead
5. **Easy Mode Paradox**: Too slow feels choppy (250ms → 180ms fixed this)
6. **Direction Buffer**: Essential for smooth controls at any speed
7. **Crash Animation**: 1-second pause gives player moment to process
8. **Lipgloss Borders**: `RoundedBorder()` looks professional with no manual drawing

---

## 🎮 Game Flow

```
Main Menu (ENTER)
    ↓
Difficulty Selection (↑↓ to select, ENTER to start, ESC to go back)
    ↓
Playing (arrow keys to move, P to pause)
    ↓
Crash Animation (😵 for 1 second)
    ↓
Game Over (R to restart with same difficulty, M for menu, Q to quit)
```

---

## 📊 Statistics

- **Total Lines of Code**: ~550 LOC
- **Files**: 9 game files + 1 main.go
- **Game States**: 6 (Menu, DifficultySelect, Playing, Paused, Crashed, GameOver)
- **Snake Expressions**: 3 (😊 😮 😵)
- **Difficulty Levels**: 3 (Easy, Medium, Hard)
- **Grid Size**: 24×20 cells (playable) + 26×22 with walls
- **Speed Range**: 50ms - 200ms depending on difficulty
- **Control Schemes**: 3 (Arrows, WASD, HJKL)

---

## 🔜 Possible Future Enhancements

- [ ] **Local high score persistence** - Save to file
- [ ] **Power-ups** - Slow down, speed up, invincibility
- [ ] **Obstacles** - Walls in middle of board
- [ ] **Multiple food types** - Different point values
- [ ] **Snake skins** - Different emoji sets
- [ ] **Sound effects** - Eating, crash, etc.
- [ ] **Leaderboard** - Top 10 scores
- [ ] **Custom difficulty** - Player-defined settings

---

**Status**: ✅ Complete and polished! Ready to merge to main.

**Highlights**: Emoji expressions, difficulty selection, visible walls, crash animation, perfect aspect ratio, responsive controls across all difficulty levels.

🐍 **The snake is ready to slither!**
