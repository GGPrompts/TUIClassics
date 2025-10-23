# 🚀 Quick Start Prompt for 2048 Implementation

Copy and paste this into your next Claude Code session when working in the 2048 worktree:

---

## Session Goal: Implement 2048 Game (v1.0.0)

Working on TUI Classics 2048 game in `/home/matt/projects/TUIClassics-2048` (git worktree on branch `feature/2048`).

**Context:**
- Full implementation plan: `/docs/2048/PLAN.md`
- Following TUITemplate 9-file modular architecture pattern
- Reference CLAUDE.md for critical patterns (mouse coords, centering, grid alignment)
- Estimated time: 3-4 hours (~450 LOC)

**Implementation Phases:**

### Phase 1: Core Data Structures (~100 LOC)
- `games/2048/types.go` - Direction, GameState, Tile, Model structs
- `games/2048/model.go` - New() and startGame() initialization

### Phase 2: Tile Spawning (~60 LOC)
- `games/2048/spawn.go` - spawnTile(), hasEmptyCell()
- 90% chance of 2, 10% chance of 4

### Phase 3: Grid Movement (~150 LOC)
- `games/2048/grid.go` - move(), directional movements, canMove(), hasValue()
- Use transpose + reverse transformations for all 4 directions

### Phase 4: Merge Logic (~100 LOC)
- `games/2048/merge.go` - slideAndMergeRow(), transpose(), reverseRows()
- Core algorithm: slide non-zero values, merge adjacent equal values

### Phase 5: Input Handling (~40 LOC)
- `games/2048/update_keyboard.go` - Arrow keys (↑↓←→) + WASD + vim keys (hjkl)
- States: Menu, Playing, Won, GameOver

### Phase 6: Rendering (~100 LOC)
- `games/2048/view.go` - renderMenu(), renderGame(), renderWin(), renderGameOver()
- Use box-drawing characters for grid (IMPORTANT: ensure UTF-8 encoding!)
- Manual centering (avoid lipgloss.Place for precise layout)

### Phase 7: Styling (~50 LOC)
- `games/2048/styles.go` - getTileStyle() with color progression
- Colors: 2=gray, 8=orange, 32=red, 128=yellow, 2048=green, >2048=magenta

### Phase 8: Main Loop (~40 LOC)
- `games/2048/update.go` - Update() dispatcher, WindowSizeMsg handling
- `cmd/2048/main.go` - Launcher (21 lines max)

### Phase 9: Menu Integration
- Add 2048 to `games/menu/model.go` imports
- Register game with hotkey "2"
- Update Makefile with build targets

**Implementation Order (suggested):**
1. types.go → model.go → spawn.go (get tiles appearing)
2. merge.go → grid.go (get movement working)
3. update_keyboard.go → update.go (make it playable)
4. view.go → styles.go (make it pretty)
5. cmd/2048/main.go → menu integration (make it launchable)

**Critical Patterns to Follow:**

From CLAUDE.md:
- **Grid Alignment**: Use single-width characters only (no emojis!)
- **Box-Drawing**: Use `┌─┬─┐├─┼─┤└─┴─┘│` (ensure UTF-8 encoding)
- **Centering**: Manual padding calculation (avoid lipgloss.Place for grids)
- **Line Counting**: `"content\n\n"` = 2 lines total (content + 1 blank), not 3!

**Key Algorithm: Movement via Transformations**
```
Left:  moveLeft() directly
Right: reverse → moveLeft() → reverse
Up:    transpose → moveLeft() → transpose
Down:  transpose → reverse → moveLeft() → reverse → transpose
```

**Testing Checklist:**
- [ ] Grid displays correctly (4×4 with borders)
- [ ] Arrow keys move tiles in correct directions
- [ ] Matching tiles merge (2+2=4, 4+4=8, etc.)
- [ ] Score increases on merges
- [ ] New tile spawns after each move
- [ ] Win at 2048 with continue option
- [ ] Game over when no moves possible
- [ ] Restart works (R key)

**Build Commands:**
```bash
# Standalone game
make 2048
./bin/2048

# Via launcher menu
make classics
./bin/classics  # Press '2' for 2048
```

**Success Criteria:**
- ✅ All movement directions work correctly
- ✅ Tiles merge and score updates
- ✅ Win/lose states trigger properly
- ✅ Grid stays aligned at all times
- ✅ Game playable via launcher menu

**When Complete:**
1. Test thoroughly (all directions, win, lose, restart)
2. Build both standalone and classics launcher
3. Create PR from `feature/2048` → `main`
4. Update CLAUDE.md with any new patterns discovered

---

**Ready to merge! 🎲 Let's build 2048!**
