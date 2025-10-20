# Solitaire Implementation Prompt

**Use this prompt to start implementing Solitaire in a new session:**

---

## Prompt

I'm working on TUIClassics, a collection of nostalgic terminal games. Minesweeper is complete (v0.1.0), and now I want to implement Solitaire (Klondike) for v0.2.0.

**Project Location**: `~/projects/TUIClassics`

**Architecture**: This project follows the modular pattern from [TUITemplate](https://github.com/GGPrompts/TUITemplate). Each game uses the 9-file structure:
- `types.go` - Type definitions, structs, enums
- `model.go` - Model initialization, game logic
- `update.go` - Main message dispatcher
- `update_keyboard.go` - Keyboard event handling
- `update_mouse.go` - Mouse event handling
- `view.go` - View rendering and layout
- `styles.go` - Lipgloss styling
- Plus game-specific files (e.g., `cards.go`, `validation.go`)

**Dependencies**:
- `github.com/charmbracelet/bubbletea` (TUI framework)
- `github.com/charmbracelet/lipgloss` (styling)

---

## Context Documents to Read

Before starting, please read these files for context:

1. **PLAN.md** - Contains detailed Solitaire specifications (Phase 2)
2. **CLAUDE.md** - Development notes with critical patterns from Minesweeper:
   - Mouse coordinate calculation patterns
   - Animation implementation (for waterfall effect)
   - Single-width vs multi-width character handling
   - Testing checklist
3. **CHANGELOG.md** - Shows what's been completed in Minesweeper

---

## Solitaire Requirements (from PLAN.md)

### Core Gameplay
- Standard Klondike Solitaire rules
- 7 tableau columns (1-7 cards)
- 4 foundation piles (Ace to King, by suit)
- Stock pile (draw cards)
- Waste pile (drawn cards)
- Draw-1 and Draw-3 modes

### Card Rendering (CRITICAL)
**Use Unicode symbols for cards** (single-width for alignment):

```
Suits: ♠ ♥ ♦ ♣
Ranks: A 2 3 4 5 6 7 8 9 T J Q K

Example card:
┌─────┐
│K ♠  │
│     │
│  ♠ K│
└─────┘
```

**DO NOT use full Unicode card emojis (🂡🂢🂣)** - they are 2 chars wide and will break alignment like Minesweeper's emoji issue.

### Input Controls

**Mouse** (Primary):
- Click to select card
- Drag-and-drop between piles
- Double-click to auto-move to foundation
- Click stock to draw

**Keyboard** (Accessibility):
- Arrow keys to navigate
- Enter to pick up/drop
- Number keys (1-4) for foundation quick-move
- Space to draw from stock
- U for undo

### Animations

#### Waterfall (On Win) - The Signature Feature!
When all cards moved to foundation:
1. Cards cascade from top of screen
2. Random horizontal velocity
3. Gravity pulls down
4. Bounce at bottom (decreasing elasticity)
5. Spin while falling
6. 3-5 seconds, then "You Won!"

**Implementation Pattern** (from CLAUDE.md):
```go
type WaterfallCard struct {
    card     Card
    x, y     float64  // Floats for smooth movement
    vx, vy   float64  // Velocity
    rotation float64  // Spin
}

func (m *Model) animateWaterfall() {
    // Physics: gravity, bounce, spin
    // Tick at 60fps (16ms)
}
```

---

## File Structure to Create

```
games/solitaire/
├── types.go          (Card, Pile, Move, SolitaireState)
├── model.go          (Deal, shuffle, initialization)
├── update.go         (Message dispatcher)
├── update_keyboard.go (Arrow keys, Enter, Undo)
├── update_mouse.go   (Drag-and-drop, click detection)
├── view.go           (Render piles, cards, layout)
├── cards.go          (Card logic, suit/rank helpers)
├── validation.go     (Move validation rules)
├── animations.go     (Waterfall animation!)
└── styles.go         (Card colors, borders)

cmd/solitaire/
└── main.go           (21 lines max - just launch)
```

---

## Key Lessons from Minesweeper (CLAUDE.md)

### 1. Mouse Coordinate Calculation
**MUST match rendering exactly**. Calculate layout positions consistently:
```go
topPadding := (termHeight - totalLines) / 2
gridStartY := topPadding + headerLines
```

When rendering `content + "\n\n"`, the content is on current line, then `\n\n` moves cursor forward 2 lines.

### 2. Single-Width Characters
Use consistent character widths for grid layouts. Emojis are OK for standalone elements (like smiley face), but NOT in grids.

### 3. Animation Pattern
```go
type AnimationTickMsg time.Time

func animationTick() tea.Cmd {
    return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
        return AnimationTickMsg(t)
    })
}

// In Update():
case AnimationTickMsg:
    if m.animating {
        m.progressAnimation()
        return m, animationTick()
    }
```

### 4. Drag-and-Drop Pattern
```go
type Model struct {
    draggingCard *Card
    dragStartX   int
    dragStartY   int
}

// Mouse down: Start drag
case tea.MouseActionPress:
    m.draggingCard = m.getCardAt(msg.X, msg.Y)

// Mouse up: Drop
case tea.MouseActionRelease:
    if m.draggingCard != nil {
        m.dropCard(msg.X, msg.Y)
        m.draggingCard = nil
    }
```

---

## Implementation Steps

### Phase 1: Core Structure (Start Here!)
1. Create file structure in `games/solitaire/`
2. Define types:
   - `Card` struct (Suit, Rank, FaceUp)
   - `Pile` struct (Cards, PileType)
   - `Model` struct (7 tableau, 4 foundation, stock, waste)
3. Implement `cards.go`:
   - Deck creation (52 cards)
   - Shuffle algorithm
   - Card comparison functions
4. Implement `model.go`:
   - `New()` - Create game
   - `InitGame()` - Deal cards (1-7 to tableau)
   - Basic layout calculation
5. Create minimal `view.go`:
   - Render placeholder piles
   - Show card backs (┌─────┐ with pattern)
   - Centered layout

**Goal**: Get a basic non-interactive view showing the dealt cards

### Phase 2: Validation & Basic Moves
1. Implement `validation.go`:
   - `canMoveToTableau()` - Alternating colors, descending rank
   - `canMoveToFoundation()` - Same suit, ascending from Ace
   - `canDrawFromStock()`
2. Add keyboard controls:
   - Arrow keys to navigate
   - Enter to select/move
   - Space to draw from stock
3. Update view to highlight selected card

**Goal**: Playable game with keyboard (no drag-and-drop yet)

### Phase 3: Mouse & Drag-Drop
1. Implement click detection for all piles
2. Add drag-and-drop:
   - Track mouse down/up
   - Visual feedback (card follows cursor)
   - Drop validation
3. Double-click auto-move to foundation

**Goal**: Full mouse support with drag-and-drop

### Phase 4: Waterfall Animation
1. Detect win condition (all cards in foundation)
2. Implement `WaterfallCard` physics
3. Render falling cards over existing view
4. Transition to win screen after animation

**Goal**: Satisfying win animation!

### Phase 5: Polish
- Undo functionality (stack of moves)
- Draw-1 vs Draw-3 mode selection
- Statistics tracking
- Hint system (optional)

---

## Testing Checklist

From CLAUDE.md, test these as you build:

### Grid Alignment
- [ ] All cards same width (use Unicode box chars)
- [ ] Piles stay aligned when cards added/removed
- [ ] Test on different terminal sizes

### Mouse Coordinates
- [ ] Click detection matches visual position
- [ ] Test all piles (tableau, foundation, stock, waste)
- [ ] Drag-and-drop drops at correct location

### Keyboard Navigation
- [ ] Arrow keys move between piles correctly
- [ ] Enter picks up and drops cards properly
- [ ] Space draws from stock

### Animations
- [ ] Waterfall runs at smooth 60fps
- [ ] Cards bounce and spin naturally
- [ ] Animation doesn't block other input

---

## Card Rendering Example

From PLAN.md, here's how to render cards:

```go
func (c Card) Render(faceUp bool) string {
    if !faceUp {
        return renderCardBack()
    }

    suit := c.SuitSymbol()  // ♠ ♥ ♦ ♣
    rank := c.RankSymbol()  // A 2-9 T J Q K
    color := c.Color()      // Red or Black

    top := fmt.Sprintf("│%s %-2s│", rank, suit)
    mid := "│     │"
    bot := fmt.Sprintf("│%2s %s│", suit, rank)

    return lipgloss.NewStyle().
        Foreground(color).
        Render(strings.Join([]string{
            "┌─────┐",
            top,
            mid,
            bot,
            "└─────┘",
        }, "\n"))
}

func renderCardBack() string {
    return "┌─────┐\n│░░░░░│\n│░░░░░│\n│░░░░░│\n└─────┘"
}
```

---

## Expected Timeline

From PLAN.md:
- Core gameplay: 8-12 hours
- Drag-and-drop: 3-4 hours
- Waterfall animation: 2-3 hours
- Polish: 2-3 hours
- **Total**: 15-22 hours (2-3 weekend sessions)

---

## Questions to Ask Me

Before you start implementing, ask me:

1. **Draw mode preference**: Start with Draw-1 or Draw-3? (I suggest Draw-1 for simplicity)
2. **Undo implementation**: Should undo be implemented in Phase 1 or saved for polish?
3. **Scoring**: Do I want Windows Solitaire scoring (+10 to foundation, -15 from foundation, +5 reveal, etc.)?
4. **Statistics**: Track games played, win rate, best time now or later?

---

## Success Criteria

When Solitaire is complete, it should:
- ✅ Play standard Klondike rules correctly
- ✅ Support both mouse (drag-drop) and keyboard
- ✅ Have beautiful card rendering with Unicode
- ✅ Show satisfying waterfall animation on win
- ✅ Be as polished as Minesweeper
- ✅ Follow the same modular architecture

---

## Build & Test Commands

```bash
# After creating files
cd ~/projects/TUIClassics

# Build
make solitaire

# Run
./bin/solitaire

# Or build and run
make run-solitaire
```

---

## Ready to Start!

Please:
1. Read PLAN.md (Phase 2 section)
2. Read CLAUDE.md (key patterns)
3. Ask me the clarifying questions above
4. Then start with Phase 1 (core structure and basic rendering)

Let's build an amazing Solitaire game! 🎴
