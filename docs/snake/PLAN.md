# Snake - Implementation Plan

**Type**: Classic Movement Game
**Complexity**: Simple (~400 LOC)
**Time Estimate**: 3-4 hours
**Status**: Ready to implement

---

## 🐍 Game Overview

Classic Snake game where you control a growing snake, eat food, and avoid hitting walls or yourself.

### Core Mechanics:
- Snake moves in one direction continuously
- Arrow keys change direction
- Eating food makes snake grow
- Game over if you hit wall or yourself
- Score increases with each food eaten

### Visual Layout:
```
┌────────────────────────────────┐
│  Score: 15    High: 42         │
├────────────────────────────────┤
│                                │
│      ●●●                       │
│        ●                       │
│        ●●                      │
│                                │
│              ◆                 │  ← Food
│                                │
│                                │
└────────────────────────────────┘
Press ↑↓←→ to move | Q to quit
```

---

## 📁 File Structure (9-file modular pattern)

```
games/snake/
├── types.go           - Game state, snake, direction enums
├── model.go           - Model initialization
├── update.go          - Main update loop & game tick
├── update_keyboard.go - Direction control
├── view.go            - Rendering
├── styles.go          - Visual styles
├── snake.go           - Snake logic (movement, growth)
├── collision.go       - Collision detection
└── food.go            - Food spawning
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

type Point struct {
    X int
    Y int
}

type GameState int
const (
    StateMenu GameState = iota
    StatePlaying
    StatePaused
    StateGameOver
)

type Model struct {
    state      GameState
    snake      []Point   // Head is at index 0
    direction  Direction
    nextDir    Direction // Buffer for next move
    food       Point
    score      int
    highScore  int

    // Game board
    width      int
    height     int
    speed      time.Duration // Time between moves

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
        width:     30,
        height:    20,
        speed:     100 * time.Millisecond,
        highScore: 0,
    }
}

func (m *Model) startGame() {
    // Place snake in middle
    centerX := m.width / 2
    centerY := m.height / 2

    m.snake = []Point{
        {centerX, centerY},
        {centerX - 1, centerY},
        {centerX - 2, centerY},
    }

    m.direction = Right
    m.nextDir = Right
    m.score = 0
    m.spawnFood()
    m.state = StatePlaying
}
```

---

### Phase 2: Snake Movement (~100 LOC)

**File**: `snake.go`

**Key Functions**:
```go
func (m *Model) moveSnake() {
    // Update direction (prevent 180-degree turns)
    if !m.isOpposite(m.direction, m.nextDir) {
        m.direction = m.nextDir
    }

    // Calculate new head position
    head := m.snake[0]
    var newHead Point

    switch m.direction {
    case Up:
        newHead = Point{head.X, head.Y - 1}
    case Down:
        newHead = Point{head.X, head.Y + 1}
    case Left:
        newHead = Point{head.X - 1, head.Y}
    case Right:
        newHead = Point{head.X + 1, head.Y}
    }

    // Add new head
    m.snake = append([]Point{newHead}, m.snake...)

    // Check if ate food
    if newHead == m.food {
        m.score++
        m.spawnFood()
        // Snake grows (don't remove tail)
    } else {
        // Remove tail (snake moves forward)
        m.snake = m.snake[:len(m.snake)-1]
    }
}

func (m *Model) isOpposite(d1, d2 Direction) bool {
    return (d1 == Up && d2 == Down) ||
           (d1 == Down && d2 == Up) ||
           (d1 == Left && d2 == Right) ||
           (d1 == Right && d2 == Left)
}
```

---

### Phase 3: Collision Detection (~80 LOC)

**File**: `collision.go`

**Functions**:
```go
func (m *Model) checkCollision() bool {
    head := m.snake[0]

    // Wall collision
    if head.X < 0 || head.X >= m.width ||
       head.Y < 0 || head.Y >= m.height {
        return true
    }

    // Self collision (check if head hits body)
    for i := 1; i < len(m.snake); i++ {
        if head == m.snake[i] {
            return true
        }
    }

    return false
}

func (m *Model) gameOver() {
    m.state = StateGameOver
    if m.score > m.highScore {
        m.highScore = m.score
        // Could save to file
    }
}
```

---

### Phase 4: Food System (~60 LOC)

**File**: `food.go`

```go
func (m *Model) spawnFood() {
    // Keep trying random positions until we find an empty spot
    for {
        food := Point{
            X: rand.Intn(m.width),
            Y: rand.Intn(m.height),
        }

        // Check if food spawned on snake
        onSnake := false
        for _, segment := range m.snake {
            if food == segment {
                onSnake = true
                break
            }
        }

        if !onSnake {
            m.food = food
            return
        }
    }
}
```

---

### Phase 5: Update Loop (~60 LOC)

**File**: `update.go`

```go
type TickMsg time.Time

func tickCmd(d time.Duration) tea.Cmd {
    return tea.Tick(d, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        tea.WindowSize(),
        tickCmd(m.speed),
    )
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.termWidth = msg.Width
        m.termHeight = msg.Height
        return m, nil

    case tea.KeyMsg:
        return m.handleKeyPress(msg)

    case TickMsg:
        if m.state == StatePlaying {
            m.moveSnake()

            if m.checkCollision() {
                m.gameOver()
            }

            // Speed up as score increases
            newSpeed := 100*time.Millisecond - time.Duration(m.score)*2*time.Millisecond
            if newSpeed < 30*time.Millisecond {
                newSpeed = 30 * time.Millisecond
            }
            m.speed = newSpeed
        }

        return m, tickCmd(m.speed)
    }

    return m, nil
}
```

---

### Phase 6: Keyboard Input (~50 LOC)

**File**: `update_keyboard.go`

```go
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateMenu:
        switch msg.String() {
        case "enter":
            m.startGame()
            return m, tickCmd(m.speed)
        }

    case StatePlaying:
        switch msg.String() {
        case "up", "k", "w":
            m.nextDir = Up
        case "down", "j", "s":
            m.nextDir = Down
        case "left", "h", "a":
            m.nextDir = Left
        case "right", "l", "d":
            m.nextDir = Right
        case "p", " ":
            m.state = StatePaused
        }

    case StatePaused:
        switch msg.String() {
        case "p", " ":
            m.state = StatePlaying
        }

    case StateGameOver:
        switch msg.String() {
        case "r":
            m.startGame()
            return m, tickCmd(m.speed)
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
    case StatePlaying, StatePaused:
        return m.renderGame()
    case StateGameOver:
        return m.renderGameOver()
    }
    return ""
}

func (m Model) renderGame() string {
    var b strings.Builder

    // Top border with score
    b.WriteString("┌" + strings.Repeat("─", m.width*2) + "┐\n")
    b.WriteString(fmt.Sprintf("│  Score: %d    High: %d%s│\n",
        m.score, m.highScore,
        strings.Repeat(" ", m.width*2-20)))

    b.WriteString("├" + strings.Repeat("─", m.width*2) + "┤\n")

    // Game board
    for y := 0; y < m.height; y++ {
        b.WriteString("│")
        for x := 0; x < m.width; x++ {
            point := Point{x, y}

            if point == m.snake[0] {
                // Snake head
                b.WriteString(headStyle.Render("●●"))
            } else if m.isSnakeBody(point) {
                // Snake body
                b.WriteString(bodyStyle.Render("●●"))
            } else if point == m.food {
                // Food
                b.WriteString(foodStyle.Render("◆◆"))
            } else {
                // Empty space
                b.WriteString("  ")
            }
        }
        b.WriteString("│\n")
    }

    // Bottom border
    b.WriteString("└" + strings.Repeat("─", m.width*2) + "┘\n")

    if m.state == StatePaused {
        b.WriteString("\nPAUSED - Press P to continue\n")
    } else {
        b.WriteString("\nPress ↑↓←→ to move | P to pause | Q to quit\n")
    }

    return lipgloss.Place(m.termWidth, m.termHeight,
        lipgloss.Center, lipgloss.Center, b.String())
}

func (m Model) isSnakeBody(point Point) bool {
    for _, segment := range m.snake {
        if point == segment {
            return true
        }
    }
    return false
}
```

---

## 🎨 Visual Design

**File**: `styles.go`

```go
var (
    headStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("46")).  // Green
        Bold(true)

    bodyStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("82"))   // Light green

    foodStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")). // Red
        Bold(true)

    scoreStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("226")). // Yellow
        Bold(true)
)
```

---

## ✅ Implementation Checklist

### Part 1: Setup
- [ ] Create all 9 files in `games/snake/`
- [ ] Define types (Direction, Point, GameState, Model)
- [ ] Implement `New()` and `startGame()`

### Part 2: Core Movement
- [ ] Implement `moveSnake()` in `snake.go`
- [ ] Handle direction changes
- [ ] Test: Snake moves continuously

### Part 3: Food System
- [ ] Implement `spawnFood()` in `food.go`
- [ ] Check if snake ate food
- [ ] Make snake grow when eating
- [ ] Test: Eating food increases score

### Part 4: Collision
- [ ] Implement `checkCollision()` in `collision.go`
- [ ] Detect wall hits
- [ ] Detect self-collision
- [ ] Test: Game over on collision

### Part 5: Input
- [ ] Handle arrow keys in `update_keyboard.go`
- [ ] Prevent 180-degree turns
- [ ] Buffer direction changes
- [ ] Test: Smooth direction control

### Part 6: Rendering
- [ ] Render game board in `view.go`
- [ ] Draw snake with distinct head/body
- [ ] Draw food
- [ ] Show score
- [ ] Test: Game looks good on screen

### Part 7: Polish
- [ ] Add pause functionality
- [ ] Show game over screen
- [ ] Add restart option
- [ ] Increase speed as score increases
- [ ] Save high score

---

## 🎯 Success Criteria

When complete, you should be able to:
1. ✅ Control snake with arrow keys
2. ✅ Eat food to grow and score
3. ✅ Game over on wall/self collision
4. ✅ See score and high score
5. ✅ Pause and resume
6. ✅ Restart after game over

---

## 🚀 Quick Start

```bash
cd ~/projects/TUIClassics/games/snake

# Create files
touch types.go model.go update.go update_keyboard.go view.go styles.go snake.go collision.go food.go

# Implementation order:
# 1. types.go - Define all types
# 2. model.go - New() and startGame()
# 3. snake.go - moveSnake()
# 4. food.go - spawnFood()
# 5. collision.go - checkCollision()
# 6. update.go - Game loop
# 7. update_keyboard.go - Controls
# 8. view.go - Rendering
# 9. styles.go - Colors

# Register in menu (games/menu/model.go)

# Build
cd ../..
make classics
./bin/classics
```

---

## 💡 Pro Tips

1. **Start with static rendering** - Draw board first
2. **Test movement early** - Get snake moving smoothly
3. **Direction buffer** - Store next direction to prevent missed inputs
4. **Speed curve** - Don't make it too fast too quickly
5. **Visual feedback** - Distinct head color helps orientation

---

## 🔜 Future Enhancements

- **Difficulty levels** - Different board sizes/speeds
- **Power-ups** - Slow down, speed up, invincibility
- **Obstacles** - Walls in the middle of board
- **Multiplayer** - Two snakes competing
- **Different game modes** - Walls that wrap around, etc.

---

**Ready to slither! 🐍 Let's build Snake!**
