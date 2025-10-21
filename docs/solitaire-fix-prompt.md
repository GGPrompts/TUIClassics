# Solitaire Bug Fix Session - Quick Start Prompt

**Project**: TUI Classics - Solitaire Game
**Location**: `/home/matt/projects/TUIClassics`
**Current Status**: ~85% complete, keyboard works perfectly, mouse has bugs

---

## 🎯 Session Goals

Fix two critical bugs in the Solitaire game to complete v0.2.0:

1. **Card border visibility** - Borders are black, invisible on black terminal backgrounds
2. **Mouse coordinate offset** - Drag-and-drop broken due to Y coordinate calculation bug

---

## 🔴 Priority 1: Fix Card Border Visibility (QUICK WIN!)

### Problem
Card borders are currently black (`blackColor`) which makes them invisible on black terminal backgrounds. This makes the rounded corners appear square and cards look broken.

### Files to Fix
- `games/solitaire/styles.go` - Change border colors

### Solution
Change card border colors from black to a visible color (white, gray, or cyan):

```go
// Current (INVISIBLE on black backgrounds):
Border(lipgloss.RoundedBorder()).
BorderForeground(blackColor)

// Fix options:
BorderForeground(lipgloss.Color("#FFFFFF"))  // White
BorderForeground(lipgloss.Color("#808080"))  // Gray
BorderForeground(lipgloss.Color("#00FFFF"))  // Cyan
```

### Testing
1. Run `make run-solitaire` or launch via `./bin/classics`
2. Verify rounded corners are clearly visible
3. Check on black terminal background

---

## 🔴 Priority 2: Fix Mouse Coordinate Offset

### Problem
Mouse Y coordinates have a massive offset, making drag-and-drop completely broken:
- Clicking waste pile (expected Y=8) reports Y=86 (offset of 78)
- Clicking tableau (expected Y=15) reports Y=26 (offset of 11)
- **Drag-and-drop is unusable**

### Current Symptoms
- Left-click on stock pile works (drawing cards)
- Right-click auto-move to foundation works
- **Cannot select or drag cards from tableau/waste**

### Files to Investigate
- `games/solitaire/update_mouse.go` - Mouse event handling
- `games/solitaire/view.go` - Card rendering and layout

### Debugging Strategy

#### Step 1: Compare with Working Minesweeper
Minesweeper's mouse handling works perfectly. Key file:
- `games/minesweeper/update_mouse.go:66` - Correct coordinate calculation

**Minesweeper's approach:**
```go
// Grid starts after: topPadding + title line + "\n\n" (2) + stats line + "\n\n" (2)
gridStartY := topPadding + 4
```

**Key lesson from CLAUDE.md**: Mouse coordinate calculations MUST exactly match rendering logic.

#### Step 2: Add Debug Output
Add temporary debug logging to `games/solitaire/update_mouse.go`:

```go
case tea.MouseMsg:
    // DEBUG: Print mouse coordinates and terminal dimensions
    fmt.Printf("DEBUG Mouse: X=%d Y=%d | Terminal: %dx%d\n",
        msg.X, msg.Y, m.termWidth, m.termHeight)
```

#### Step 3: Trace Rendering Logic
In `games/solitaire/view.go`, trace exactly where each element renders:
- How many lines does the header take?
- Are there any `\n\n` sequences adding extra spacing?
- Is `lipgloss.Place` being used for centering? (This breaks coordinate math!)

**From CLAUDE.md:**
> "For interactive elements (clickable grids), avoid `lipgloss.Place`. Use manual centering so mouse coordinate math is deterministic."

#### Step 4: Calculate Correct Offsets
Once you understand the rendering:
1. Calculate where each pile starts (X, Y coordinates)
2. Update the click detection logic in `update_mouse.go`
3. Test clicking each pile type: stock, waste, foundation, tableau

#### Step 5: Test Drag-and-Drop
After fixing click detection:
1. Test selecting a card (mouse down)
2. Test dragging (mouse motion with button held)
3. Test dropping (mouse up)
4. Verify across all pile types

---

## 📋 Implementation Checklist

### Card Border Fix
- [ ] Open `games/solitaire/styles.go`
- [ ] Find all `BorderForeground(blackColor)` occurrences
- [ ] Change to visible color (white/gray/cyan)
- [ ] Test on black terminal background
- [ ] Verify rounded corners are visible

### Mouse Coordinate Fix
- [ ] Add debug output to capture mouse coordinates
- [ ] Run game and click various piles, note coordinates
- [ ] Compare with expected coordinates from `view.go`
- [ ] Check if `lipgloss.Place` is being used (remove if so)
- [ ] Calculate correct Y offset based on rendering logic
- [ ] Update coordinate calculations in `update_mouse.go`
- [ ] Test click detection on all pile types
- [ ] Test drag-and-drop between piles
- [ ] Remove debug output

---

## 🧪 Testing Protocol

After fixes:

1. **Visual Test**
   - Cards have visible rounded borders
   - Colors look good on black background

2. **Mouse Test (Stock Pile)**
   - Click stock to draw cards ✓

3. **Mouse Test (Waste Pile)**
   - Click to select card
   - Drag to tableau or foundation

4. **Mouse Test (Tableau)**
   - Click to select card
   - Drag within tableau (build sequences)
   - Drag to foundation

5. **Mouse Test (Foundation)**
   - Right-click card to auto-move to foundation

6. **Edge Cases**
   - Try dragging invalid moves (should reject)
   - Test with different terminal sizes
   - Verify keyboard controls still work

---

## 📁 Key Files Reference

| File | Purpose | Lines of Interest |
|------|---------|-------------------|
| `games/solitaire/styles.go` | Card styling | Border colors |
| `games/solitaire/view.go` | Layout rendering | Pile positions, spacing |
| `games/solitaire/update_mouse.go` | Mouse events | Click detection, drag-and-drop |
| `games/minesweeper/update_mouse.go` | Working example | Line 66 - correct offset calc |

---

## 🎓 Key Principles from CLAUDE.md

1. **Mouse coordinates must match rendering exactly**
   - When rendering `content + "\n\n"`, content occupies current line, then `\n\n` moves cursor forward 2 lines
   - Next content is written where cursor lands, NOT 2 additional lines later

2. **Avoid `lipgloss.Place` for interactive elements**
   - Use manual centering with calculated padding
   - Makes coordinate math deterministic

3. **All grid characters must be single-width**
   - Already fixed for cards (7 chars wide × 5 lines tall)

---

## 🚀 Quick Start Commands

```bash
# Navigate to project
cd /home/matt/projects/TUIClassics

# Build Solitaire
make solitaire

# Run Solitaire standalone
./bin/solitaire

# OR run via launcher
make classics
./bin/classics
# Then press 's' to launch Solitaire

# Run with immediate testing
make run-solitaire
```

---

## 📝 Success Criteria

Session is complete when:
- ✅ Card borders are visible on black terminal backgrounds
- ✅ Can click and select cards from waste pile
- ✅ Can click and select cards from tableau
- ✅ Can drag cards between piles
- ✅ Drag-and-drop validates moves correctly
- ✅ Game is fully playable with mouse
- ✅ Keyboard controls still work
- ✅ No debug output in final code

---

## 💡 Expected Time

- **Border fix**: 5-10 minutes (quick find-and-replace)
- **Mouse coordinate debugging**: 30-60 minutes (investigation + fix)
- **Testing**: 15-30 minutes

**Total**: 1-2 hours

---

**Good luck! The game is so close to complete! 🎉**
