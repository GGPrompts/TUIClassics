# 🚀 Quick Copy-Paste Prompt for Next Session

Copy and paste this into your next Claude Code session:

---

## Session Goal: Fix Solitaire Bugs (v0.2.0)

Working on TUI Classics Solitaire game in `/home/matt/projects/TUIClassics`.

**Two critical bugs to fix:**

### 1. Card Border Visibility (5 min fix)
- **Problem**: Card borders use black color, invisible on black terminal backgrounds
- **Fix**: In `games/solitaire/styles.go`, change all `BorderForeground(blackColor)` to white/gray/cyan
- **Test**: Run game, verify rounded corners are visible on black background

### 2. Mouse Drag-and-Drop Broken (30-60 min fix)
- **Problem**: Y coordinate offset makes clicking wrong positions
  - Clicking waste (expect Y=8) → reports Y=86 (78 offset)
  - Clicking tableau (expect Y=15) → reports Y=26 (11 offset)
- **Reference**: Check `games/minesweeper/update_mouse.go:66` - this works perfectly
- **Key lesson from CLAUDE.md**: Mouse coordinates must exactly match rendering logic
  - When rendering `content + "\n\n"`, cursor moves 2 lines, not 4
  - Avoid `lipgloss.Place` for interactive elements - use manual centering
- **Files**:
  - `games/solitaire/update_mouse.go` - click detection
  - `games/solitaire/view.go` - rendering logic
- **Debug approach**:
  1. Add temp debug: `fmt.Printf("Mouse: X=%d Y=%d\n", msg.X, msg.Y)`
  2. Click piles, note actual vs expected coordinates
  3. Trace `view.go` to see where piles render
  4. Calculate correct offset
  5. Fix coordinate calculations in `update_mouse.go`
  6. Remove debug output

**Success criteria**:
- ✅ Visible card borders on black terminal
- ✅ Can click and drag cards between all piles
- ✅ Mouse and keyboard both work perfectly

**Commands**:
```bash
make solitaire      # Build
./bin/solitaire     # Run standalone
# OR
./bin/classics      # Launch menu, press 's'
```

See full details in `docs/solitaire-fix-prompt.md`.

---
