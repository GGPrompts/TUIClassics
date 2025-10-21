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

## ✅ Phase 2: Solitaire (Klondike)

**Status**: v0.2.0 - Feature complete with waterfall animation

**Priority**: High - Most requested classic game after Minesweeper

**Current Status** (2025-01-21): 100% complete, all features working

### Core Gameplay

- [x] Standard Klondike Solitaire rules
- [x] 7 tableau columns (1-7 cards)
- [x] 4 foundation piles (Ace to King, by suit)
- [x] Stock pile (draw cards)
- [x] Waste pile (drawn cards)
- [x] Draw-1 mode (working)
- [ ] Draw-3 mode toggle (planned)

### Input Methods ✅ COMPLETE

**Mouse** (Primary):
- [x] Click-to-select (select with mouse, move with keyboard)
- [x] Drag-and-drop between piles (full visual feedback)
- [x] Right-click to auto-move to foundation
- [x] Click stock to draw card
- [x] Gold border shows selected/dragging cards
- [x] Accurate coordinate detection (centered layout)

**Keyboard** (Accessibility):
- [x] Arrow keys to navigate cards
- [x] Enter to pick up/drop card
- [x] Number keys (1-4) for quick foundation moves
- [x] Space to draw from stock
- [x] **W key** - Secret waterfall animation trigger
- [ ] U for undo (future enhancement)

### Card Rendering ✅ COMPLETE

**Implementation**: Unicode symbols with lipgloss rounded borders
```
Suits: ♠ ♥ ♦ ♣
Ranks: A 2 3 4 5 6 7 8 9 T J Q K

Actual rendering:
╭─────╮
│A ♠  │
│     │
│♠ A  │
╰─────╯
```

**Dimensions**:
- Width: 7 chars (5 content + 2 borders)
- Height: 5 lines (3 content + 2 borders)
- Style: `.Width(5).Height(3)` enforces consistent sizing

**Stacking**: Overlapping cards show top 2 lines (rounded border + rank/suit)

**Fixed Issues**:
- ✅ All cards now have consistent dimensions
- ✅ Proper centered alignment with `lipgloss.Center`
- ✅ Stacked cards match full card width
- ✅ No clipped symbols (hearts render correctly)

**Fixed Issues**:
- ✅ Gray borders now visible on all terminal backgrounds
- ✅ Mouse coordinates accurate with centered layout
- ✅ Cursor visible on empty foundation/tableau piles
- ✅ Hybrid mouse+keyboard workflow (click to select, arrows to move)

### Animations

#### 🎊 Waterfall Animation (On Win) ✅ COMPLETE

**The Signature Feature!**

When all cards are moved to foundation (or press **W** to test):
1. ✅ All visible cards cascade down from top of screen
2. ✅ Full card boxes with rounded borders
3. ✅ Random horizontal velocity and spread
4. ✅ Gravity pulls them down
5. ✅ Bounce off bottom with decreasing elasticity
6. ✅ Cards continue moving after bounce
7. ✅ After animation completes, show "You Won!" screen
8. ✅ Plain-text rendering (no ANSI artifacts)
9. ✅ Single-width suit letters (S/H/D/C) for perfect alignment

**Implementation**:
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

### Validation & Rules ✅ COMPLETE

- [x] Only King can be placed on empty tableau
- [x] Tableau stacks alternate red/black, descending rank
- [x] Foundation builds up from Ace, same suit
- [x] Can only move fully revealed card sequences
- [x] Stock cycles through deck (Draw-1 working)

### Layout & Centering ✅ COMPLETE

**Implementation**: Manual centering (like Minesweeper) for deterministic mouse coordinates
- ✅ Centered vertically and horizontally
- ✅ Title, stats, help text use lipgloss `.Width().Align()`
- ✅ Game board uses manual left padding (per-line)
- ✅ Mouse coordinates adjusted for padding offsets
- ✅ Responsive to terminal resize
- ✅ Requests `tea.WindowSize()` on init for immediate centering

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

## ✅ Phase 3: Game Launcher (Complete)

**Status**: v0.3.0 - Feature complete

**Purpose**: Unified menu to select which game to play

### Features

- [x] Main menu showing all available games
- [x] Keyboard shortcuts (M for Minesweeper, S for Solitaire, T for Tanks)
- [x] Beautiful title screen with gold/blue color scheme
- [x] Navigation with arrow keys and vim keys (j/k)
- [x] ESC to return from game to menu
- [ ] High scores for each game (planned for future)
- [ ] Last played indicator (planned for future)

### Binary: `classics`

**Launch**: `./bin/classics` or `make run-classics`

**Actual Implementation**:
```
TUI CLASSICS
Classic Terminal Games Collection

[m] Minesweeper
    Classic mine-finding puzzle game

[s] Solitaire
    Klondike card game (Work in Progress)

[t] Tanks
    Coming soon...

↑/↓: Navigate  •  Enter/Hotkey: Select  •  q: Quit
```

**Features**:
- Centered menu with gold/blue color scheme
- Selected item highlighting
- Disabled state for unimplemented games
- ESC returns from game to menu
- Seamless transitions without exiting app

**File Structure**:
```
cmd/classics/
└── main.go          (21 lines - launches menu)

games/menu/
├── types.go         (MenuState, GameInfo structs)
├── model.go         (Menu initialization, game launching)
├── view.go          (Menu rendering)
├── styles.go        (Color scheme and styling)
├── update.go        (Message routing, game delegation)
└── update_keyboard.go (Navigation, selection)
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

### Solitaire Completion (v0.2.0)

**Priority 1: Fix Card Border Visibility** 🔴
1. Change card border color from black to white/gray/cyan (visible on black terminals)
2. Test on black terminal background to confirm visibility
3. Ensure rounded corners are clearly visible

**Priority 2: Fix Mouse Drag-and-Drop** 🔴
1. Test with debug output to capture terminal dimensions
2. Identify source of Y coordinate offset (86 vs expected 8)
3. Compare coordinate calculation with Minesweeper (which works correctly)
4. Add offset compensation or fix coordinate calculation
5. Verify drag-and-drop works across all piles
6. Remove debug output once fixed

**Priority 3: Polish & Testing** 🟡
1. Add undo functionality (keyboard U key)
2. Implement Draw-3 mode toggle
3. Test on multiple terminals (Alacritty, iTerm2, Windows Terminal)
4. Verify all animations work smoothly
5. Clean up any remaining TODOs in code

**Priority 4: Documentation** 🟢
1. Add screenshots/GIFs to README
2. Document controls clearly
3. Add troubleshooting section for mouse issues
4. Update CHANGELOG with v0.2.0 release notes

---

**Created**: 2025-01-20
**Last Updated**: 2025-10-20 (Launcher complete!)
**Owner**: GGPrompts

**Status**:
- ✅ Minesweeper complete (v0.1.0) - No known bugs
- 🚧 Solitaire ~85% complete (v0.2.0) - Mouse drag blocked, keyboard works perfectly
- ✅ Launcher complete (v0.3.0) - Main menu working with game selection

**Session Notes**:
- Fixed card sizing/alignment issues (all cards now consistent dimensions)
- Fixed stacking rendering (cards overlap correctly with rounded corners)
- Waterfall animation complete and working
- Discovered mouse coordinate offset bug - needs investigation
- Game is fully playable via keyboard controls
- **NEW (2025-10-20)**: Built unified game launcher menu (Phase 3 complete!)
  - Main menu with game selection
  - ESC to return from games to menu
  - Gold/blue color scheme
  - Seamless transitions between menu and games
