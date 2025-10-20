# TUI Classics Development Plan

Roadmap for building a collection of nostalgic terminal games with beautiful TUI interfaces.

## Vision

Bring classic Windows 95/XP games to the terminal with:
- **Beautiful visuals** - Styled with Lipgloss, centered layouts, nostalgic colors
- **Dual input modes** - Full mouse and keyboard support
- **Satisfying animations** - Explosion effects, waterfall wins, smooth transitions
- **Accessibility** - Keyboard-first design with mouse as convenience
- **TFE Integration** - Launch games from TFE file manager for seamless workflow

## Current Status

### ✅ Phase 1: Minesweeper (Complete)

**Status**: v0.1.0 - Feature complete with explosion animation

**Features Implemented**:
- [x] Classic Klondike gameplay
- [x] Three difficulty levels (Easy, Medium, Hard)
- [x] Mouse controls (left-click reveal, right-click flag)
- [x] Keyboard controls (arrows, Enter, F key)
- [x] Safe first click (never hit mine on first move)
- [x] Cascade reveal for empty cells
- [x] Classic Windows 95 color scheme
- [x] **Explosion animation** - Expanding shockwave on mine hit
- [x] Timer and score tracking
- [x] Centered, polished UI
- [x] Single-width character grid (perfect alignment)

**Code Stats**:
- 6 files, ~600 lines of code
- Clean modular architecture
- No known bugs

---

## 🚧 Phase 2: Solitaire (Klondike)

**Target**: v0.2.0

**Priority**: High - Most requested classic game after Minesweeper

### Core Gameplay

- [ ] Standard Klondike Solitaire rules
- [ ] 7 tableau columns (1-7 cards)
- [ ] 4 foundation piles (Ace to King, by suit)
- [ ] Stock pile (draw cards)
- [ ] Waste pile (drawn cards)
- [ ] Draw-1 and Draw-3 modes

### Input Methods

**Mouse** (Primary):
- [ ] Click card to select
- [ ] Drag-and-drop between piles
- [ ] Double-click to auto-move to foundation
- [ ] Click stock to draw card

**Keyboard** (Accessibility):
- [ ] Arrow keys to navigate cards
- [ ] Enter to pick up/drop card
- [ ] Number keys (1-4) for quick foundation moves
- [ ] Space to draw from stock
- [ ] U for undo

### Card Rendering

**Option A: Unicode Card Symbols** (Recommended)
```
Suits: ♠ ♥ ♦ ♣
Ranks: A 2 3 4 5 6 7 8 9 T J Q K

Example layout:
┌─────┐
│K ♠  │
│     │
│  ♠ K│
└─────┘
```

**Option B: Full Unicode Cards** (If terminal supports)
```
🂡 🂢 🂣 🂤 🂥 🂦 🂧 🂨 🂩 🂪 🂫 🂬 🂭 🂮  (Spades)
🂱 🂲 🂳 🂴 🂵 🂶 🂷 🂸 🂹 🂺 🂻 🂼 🂽 🂾  (Hearts)
```
⚠️ Issue: These are 2-char wide, would break alignment like minesweeper emojis

**Decision**: Use Option A (Unicode symbols) for consistent width

### Animations

#### 🎊 Waterfall Animation (On Win)

**The Signature Feature!**

When all cards are moved to foundation:
1. Cards cascade down from top of screen
2. Each card has random horizontal velocity
3. Gravity pulls them down
4. Bounce off bottom with decreasing elasticity
5. Spin while falling
6. After 3-5 seconds, show "You Won!" screen

**Implementation Pattern**:
```go
type WaterfallCard struct {
    card     Card
    x, y     float64  // Position (floats for smooth movement)
    vx, vy   float64  // Velocity
    rotation float64  // Spin angle
}

func (m *Model) animateWaterfall() {
    // Physics simulation
    for i := range m.waterfallCards {
        card := &m.waterfallCards[i]

        // Apply gravity
        card.vy += GRAVITY

        // Update position
        card.x += card.vx
        card.y += card.vy

        // Bounce at bottom
        if card.y >= m.termHeight-2 {
            card.y = m.termHeight - 2
            card.vy = -card.vy * ELASTICITY
        }

        // Spin
        card.rotation += 0.2
    }
}
```

Tick rate: 60fps (16ms) for smooth animation

#### Other Animations

- [ ] Card flip when revealing tableau cards
- [ ] Smooth slide when auto-completing to foundation
- [ ] Subtle highlight on valid drop targets during drag

### Validation & Rules

- [ ] Only King can be placed on empty tableau
- [ ] Tableau stacks alternate red/black, descending rank
- [ ] Foundation builds up from Ace, same suit
- [ ] Can only move fully revealed card sequences
- [ ] Stock cycles through deck (1 or 3 passes depending on mode)

### UI Layout

**Initial View**:
```
┌─ Solitaire ──────────────────────────────────────────────┐
│                                                           │
│  [Stock] [Waste]          [♠] [♥] [♦] [♣]  Score: 0      │
│                                                           │
│  [ K♠ ] [ 6♥ ] [ 4♣ ] [ 2♦ ] [ A♠ ] [ 9♥ ] [ 7♣ ]       │
│         [███] [███] [███] [███] [███] [███]              │
│               [███] [███] [███] [███] [███]              │
│                     [███] [███] [███] [███]              │
│                           [███] [███] [███]              │
│                                 [███] [███]              │
│                                       [███]              │
│                                                           │
│  Time: 1:23 | [U]ndo | [D]raw | [H]int | [N]ew | [Q]uit │
└───────────────────────────────────────────────────────────┘
```

**File Structure**:
```
games/solitaire/
├── types.go          (Card, Pile, Move structs)
├── model.go          (Deal cards, shuffle, game state)
├── update.go         (Main dispatcher)
├── update_keyboard.go (Keyboard controls)
├── update_mouse.go   (Drag-and-drop, click handlers)
├── view.go           (Render piles, cards)
├── cards.go          (Card logic, suit/rank helpers)
├── validation.go     (Move validation rules)
├── animations.go     (Waterfall, flips, slides)
└── styles.go         (Card styles, colors)
```

### Scoring (Optional)

Classic Windows Solitaire scoring:
- Move to foundation: +10
- Move from foundation back to tableau: -15
- Reveal tableau card: +5
- Draw from stock: -2 (Draw-3 mode)

### Statistics (Future)

- Games played
- Games won
- Win percentage
- Best time
- Best score
- Current streak

### Estimated Timeline

- **Core gameplay**: 8-12 hours
- **Drag-and-drop**: 3-4 hours
- **Waterfall animation**: 2-3 hours
- **Polish & testing**: 2-3 hours
- **Total**: 15-22 hours (~2-3 weekend sessions)

---

## 🔮 Phase 3: Game Launcher

**Target**: v0.3.0

**Purpose**: Unified menu to select which game to play

### Features

- [ ] Main menu showing all available games
- [ ] High scores for each game
- [ ] Last played indicator
- [ ] Keyboard shortcuts (M for Minesweeper, S for Solitaire)
- [ ] Beautiful title screen

### Binary: `classics`

```
┌─ Terminal Classics ───────────────────────────────────────┐
│                                                            │
│                   🎮 TERMINAL CLASSICS 🎮                  │
│                                                            │
│              Nostalgic games for your terminal             │
│                                                            │
│   ┌────────────────────────────────────────────────────┐  │
│   │  *  Minesweeper          [M]                       │  │
│   │      Clear the minefield without exploding         │  │
│   │      Best time: 1:23  |  Games played: 47          │  │
│   ├────────────────────────────────────────────────────┤  │
│   │  ♠  Solitaire            [S]                       │  │
│   │      Classic Klondike card game                    │  │
│   │      Best score: 3,245  |  Win rate: 23%           │  │
│   ├────────────────────────────────────────────────────┤  │
│   │  🐍  Snake               [N]  (Coming Soon!)       │  │
│   │      Eat apples, grow longer, don't crash          │  │
│   └────────────────────────────────────────────────────┘  │
│                                                            │
│            Press key to play, or Q to quit                 │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

**File Structure**:
```
cmd/classics/
└── main.go          (Launcher menu, detect available games)
```

---

## 🎯 Phase 4: TFE Integration

**Target**: v0.4.0

**Purpose**: Launch games from TFE file manager

### Integration Points

1. **Auto-detection**:
   ```go
   // TFE checks if binaries exist
   if exec.LookPath("minesweeper") != nil {
       // Add to games menu
   }
   ```

2. **Context Menu** (Right-click in TFE):
   ```
   > Games >
     - Minesweeper
     - Solitaire
     - Classics Launcher
   ```

3. **Keybinding**: `Ctrl+G` to open games menu

4. **Return to TFE**: Press any key after game exits

### Implementation

- [ ] Create helper in `games/shared/tfebridge/`
- [ ] Detect if launched from TFE
- [ ] "Press any key to return to TFE" message
- [ ] Clean terminal state on exit

---

## 🎨 Phase 5: Shared Components

**Target**: v0.5.0

**Purpose**: Reusable code across all games

### High Score System

**File**: `games/shared/highscores/storage.go`

```go
type HighScoreManager struct {
    dbPath string  // ~/.local/share/tui-classics/scores.json
}

func (h *HighScoreManager) SaveScore(game, player string, score int) error
func (h *HighScoreManager) GetTopScores(game string, limit int) []Score
func (h *HighScoreManager) GetPersonalBest(game, player string) (Score, error)
```

**Storage**: JSON file, SQLite, or BoltDB

### Theme System

**File**: `games/shared/themes/theme.go`

```go
type Theme struct {
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Success    lipgloss.Color
    Error      lipgloss.Color
    Border     lipgloss.Color
    Background lipgloss.Color
}

var Themes = map[string]Theme{
    "classic": ClassicTheme,  // Windows 95 colors
    "dark":    DarkTheme,     // Modern dark mode
    "retro":   RetroTheme,    // Green monochrome
    "nord":    NordTheme,     // Nord color scheme
}
```

### Animation Framework

**File**: `games/shared/animations/framework.go`

Reusable patterns for:
- Particle effects
- Easing functions
- Physics simulation (gravity, bounce)
- Fade in/out

---

## 💡 Phase 6: Additional Games (Future)

### Snake (Easy - 4-6 hours)

- [ ] Grid-based movement
- [ ] Apple spawning
- [ ] Collision detection
- [ ] Speed increase over time
- [ ] High score tracking

### Tetris (Medium - 12-16 hours)

- [ ] 7 tetromino shapes
- [ ] Rotation (SRS system)
- [ ] Line clearing
- [ ] Gravity/drop speed
- [ ] Next piece preview
- [ ] Score/level system

### Space Invaders (Medium - 10-14 hours)

- [ ] Enemy wave movement
- [ ] Shooting mechanics
- [ ] Shields that degrade
- [ ] UFO bonus ship
- [ ] Lives system

### Pac-Man (Hard - 20-30 hours)

- [ ] Maze rendering
- [ ] Ghost AI (chase, scatter, frightened)
- [ ] Power pellet effects
- [ ] Fruit spawning
- [ ] Complex collision detection

---

## 🚀 Phase 7: Polish & Distribution

**Target**: v1.0.0

### Features

- [ ] Comprehensive README with screenshots
- [ ] Installation via `go install`
- [ ] Prebuilt binaries for Linux, macOS, Windows
- [ ] GitHub releases with changelogs
- [ ] Demo GIFs/videos for each game
- [ ] Man pages for each binary

### Testing

- [ ] Test on multiple terminal emulators (Alacritty, iTerm2, Kitty, Windows Terminal)
- [ ] Test various terminal sizes (80x24 to 200x60)
- [ ] Test on macOS, Linux, Windows WSL
- [ ] Ensure consistent colors across terminals

### Distribution Channels

- [ ] GitHub Releases (binaries)
- [ ] Homebrew tap (macOS/Linux)
- [ ] AUR package (Arch Linux)
- [ ] Snap/Flatpak (Linux)
- [ ] Scoop (Windows)

---

## 📊 Success Metrics

**Goal**: Create the go-to collection of terminal games

**Metrics**:
- GitHub stars: 100+ (proves concept resonates)
- Games completed: 3+ (Minesweeper, Solitaire, Snake minimum)
- TFE integration: Working seamlessly
- User feedback: Positive nostalgia, smooth gameplay

---

## 🎯 Next Immediate Steps

1. **Create CHANGELOG.md** ✅
2. **Finalize documentation** ✅
3. **Begin Solitaire planning**:
   - Design card struct
   - Plan pile layouts
   - Prototype card rendering
4. **Start Solitaire implementation** (Weekend 1):
   - Scaffold files
   - Implement card logic
   - Basic rendering
5. **Drag-and-drop implementation** (Weekend 2):
   - Mouse tracking
   - Drop validation
   - Visual feedback
6. **Waterfall animation** (Weekend 2):
   - Physics simulation
   - Card rendering during fall
   - Win screen integration

---

**Created**: 2025-01-20
**Last Updated**: 2025-01-20
**Owner**: GGPrompts

**Status**:
- ✅ Minesweeper complete (v0.1.0)
- 🚧 Solitaire planned (v0.2.0)
- 📋 Launcher planned (v0.3.0)
