# Keyboard Hero - Implementation Plan

**Type**: Rhythm Game (Guitar Hero style)
**Complexity**: Medium (~900 LOC)
**Time Estimate**: 6-8 hours
**Status**: ✅ Complete

---

## ⌨️ Game Overview

A Guitar Hero style rhythm game for the terminal where notes fall down lanes and you press keys at the right time.

**Note**: Originally planned as "Terminal Hero" but renamed to "Keyboard Hero" to better reflect the keyboard-focused gameplay without audio.

### Core Mechanics:
- 5 lanes (A, S, D, F, J keys)
- Notes scroll down at constant speed
- Hit zone at bottom of screen
- Timing-based scoring (Perfect/Good/OK/Miss)
- Combo multiplier system
- Song chart system

### Visual Layout:
```
     ┌─────┬─────┬─────┬─────┬─────┐
     │  A  │  S  │  D  │  F  │  J  │  ← Lane headers
     ├─────┼─────┼─────┼─────┼─────┤
     │     │     │  ●  │     │     │
     │     │  ●  │     │     │  ●  │  ← Notes falling
     │  ●  │     │     │  ●  │     │
     │     │     │  ●  │     │     │
     ├─────┼─────┼─────┼─────┼─────┤
     │ [█] │ [█] │ [█] │ [█] │ [█] │  ← Hit zone
     └─────┴─────┴─────┴─────┴─────┘

     Score: 12,450  Combo: 37x  Multiplier: 4x
     ★ PERFECT! ★
```

---

## 📁 File Structure (9-file modular pattern)

```
games/hero/
├── types.go           - Game state, notes, timing enums
├── model.go           - Model initialization
├── update.go          - Main update loop & animation
├── update_keyboard.go - Key press handling
├── view.go            - Main rendering
├── styles.go          - Visual styles & colors
├── notes.go           - Note logic & scrolling
├── scoring.go         - Hit detection & scoring
└── songs.go           - Song data & chart loading
```

---

## 🎯 Implementation Phases

### Phase 1: Core Game Loop (~300 LOC)

**Files**: `types.go`, `model.go`, `update.go`, `view.go`, `styles.go`

**types.go** - Define core structures:
```go
type GameState int
const (
    StateMenu GameState = iota
    StatePlaying
    StateFinished
)

type Note struct {
    Lane      int       // 0-4 (A/S/D/F/J)
    Y         int       // Current Y position on screen
    HitTime   time.Time // When note should be hit
    Hit       bool      // Has been successfully hit
}

type Model struct {
    state      GameState
    notes      []Note
    score      int
    combo      int
    multiplier int

    // Visual
    termWidth  int
    termHeight int
    laneWidth  int
    hitZoneY   int

    // Timing
    startTime  time.Time
    currentTime time.Time
    scrollSpeed float64 // Notes per second
}
```

**model.go** - Initialization:
```go
func New() Model {
    return Model{
        state:       StateMenu,
        scrollSpeed: 10.0, // 10 rows per second
        laneWidth:   12,
        notes:       []Note{},
    }
}
```

**update.go** - Animation loop:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case TickMsg:
        m.currentTime = time.Now()
        m.updateNotePositions() // Scroll notes down
        m.checkMissedNotes()    // Remove passed notes
        return m, tickCmd()
}

func tickCmd() tea.Cmd {
    return tea.Tick(16*time.Millisecond, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}
```

**view.go** - Basic rendering:
```go
func (m Model) View() string {
    return m.renderLanes() +
           m.renderNotes() +
           m.renderHitZone() +
           m.renderScore()
}
```

---

### Phase 2: Note Logic & Scrolling (~200 LOC)

**File**: `notes.go`

**Key Functions**:
```go
func (m *Model) updateNotePositions() {
    elapsed := m.currentTime.Sub(m.startTime).Seconds()

    for i := range m.notes {
        if !m.notes[i].Hit {
            // Calculate Y position based on time
            timeToHit := m.notes[i].HitTime.Sub(m.startTime).Seconds()
            distanceFromHit := (timeToHit - elapsed) * m.scrollSpeed
            m.notes[i].Y = m.hitZoneY - int(distanceFromHit)
        }
    }
}

func (m *Model) checkMissedNotes() {
    newNotes := []Note{}
    for _, note := range m.notes {
        // Keep notes that haven't passed hit zone yet
        if note.Y <= m.hitZoneY + 3 { // Grace period
            newNotes = append(newNotes, note)
        } else if !note.Hit {
            // Missed! Reset combo
            m.combo = 0
            m.multiplier = 1
        }
    }
    m.notes = newNotes
}

func (m *Model) getNotesInLane(lane int) []Note {
    var laneNotes []Note
    for _, note := range m.notes {
        if note.Lane == lane && !note.Hit {
            laneNotes = append(laneNotes, note)
        }
    }
    return laneNotes
}
```

---

### Phase 3: Hit Detection & Scoring (~200 LOC)

**File**: `scoring.go`

**Timing Windows**:
```go
type HitResult int
const (
    Miss HitResult = iota
    OK       // ±150ms → 25 points
    Good     // ±100ms → 50 points
    Perfect  // ±50ms → 100 points
)

func (m *Model) checkHit(lane int) HitResult {
    laneNotes := m.getNotesInLane(lane)
    if len(laneNotes) == 0 {
        return Miss
    }

    // Find closest note to hit zone
    closest := laneNotes[0]
    for _, note := range laneNotes {
        if abs(note.Y - m.hitZoneY) < abs(closest.Y - m.hitZoneY) {
            closest = note
        }
    }

    // Check timing
    distance := abs(closest.Y - m.hitZoneY)

    if distance <= 1 {
        closest.Hit = true
        m.addScore(100)
        return Perfect
    } else if distance <= 2 {
        closest.Hit = true
        m.addScore(50)
        return Good
    } else if distance <= 3 {
        closest.Hit = true
        m.addScore(25)
        return OK
    }

    return Miss
}

func (m *Model) addScore(basePoints int) {
    m.combo++

    // Update multiplier based on combo
    if m.combo >= 30 {
        m.multiplier = 4
    } else if m.combo >= 20 {
        m.multiplier = 3
    } else if m.combo >= 10 {
        m.multiplier = 2
    } else {
        m.multiplier = 1
    }

    m.score += basePoints * m.multiplier
}
```

---

### Phase 4: Song System (~200 LOC)

**File**: `songs.go`

**Song Chart Format**:
```go
type Song struct {
    Title    string
    Artist   string
    BPM      int
    Duration float64 // seconds
    Chart    []ChartNote
}

type ChartNote struct {
    Time float64 // seconds from start
    Lane int     // 0-4
}

// Built-in demo songs
var demoSongs = []Song{
    {
        Title:  "Easy Street",
        Artist: "Tutorial",
        BPM:    120,
        Chart: []ChartNote{
            {Time: 1.0, Lane: 0},
            {Time: 2.0, Lane: 1},
            {Time: 3.0, Lane: 2},
            {Time: 4.0, Lane: 3},
            {Time: 5.0, Lane: 4},
            // ... more notes
        },
    },
    {
        Title:  "Speed Demon",
        Artist: "Challenge",
        BPM:    180,
        Chart: []ChartNote{
            {Time: 0.5, Lane: 0},
            {Time: 0.6, Lane: 2},
            {Time: 0.7, Lane: 4},
            // ... rapid notes
        },
    },
}

func (m *Model) loadSong(song Song) {
    m.notes = []Note{}
    m.startTime = time.Now()

    for _, chartNote := range song.Chart {
        m.notes = append(m.notes, Note{
            Lane:    chartNote.Lane,
            HitTime: m.startTime.Add(time.Duration(chartNote.Time * float64(time.Second))),
            Y:       -10, // Start off screen
            Hit:     false,
        })
    }
}
```

---

## 🎨 Visual Design

### Colors (Neon Style):
```go
// styles.go
var (
    laneColors = []lipgloss.Color{
        lipgloss.Color("46"),  // Green (A)
        lipgloss.Color("196"), // Red (S)
        lipgloss.Color("226"), // Yellow (D)
        lipgloss.Color("21"),  // Blue (F)
        lipgloss.Color("201"), // Magenta (J)
    }

    perfectColor = lipgloss.Color("46")  // Green
    goodColor    = lipgloss.Color("226") // Yellow
    okColor      = lipgloss.Color("214") // Orange
    missColor    = lipgloss.Color("196") // Red
)
```

### Note Rendering:
```go
func (m Model) renderNotes() string {
    var b strings.Builder

    // Render each row from top to bottom
    for y := 0; y < m.hitZoneY; y++ {
        b.WriteString("│")

        for lane := 0; lane < 5; lane++ {
            // Check if any note is at this position
            hasNote := false
            for _, note := range m.notes {
                if note.Lane == lane && note.Y == y && !note.Hit {
                    hasNote = true
                    break
                }
            }

            if hasNote {
                noteStyle := lipgloss.NewStyle().
                    Foreground(laneColors[lane]).
                    Bold(true)
                b.WriteString(centerString(noteStyle.Render("●"), m.laneWidth))
            } else {
                b.WriteString(strings.Repeat(" ", m.laneWidth))
            }

            b.WriteString("│")
        }
        b.WriteString("\n")
    }

    return b.String()
}
```

---

## ✅ Implementation Checklist

### Part 1: Setup ✅ COMPLETE
- [x] Create all 9 files in `games/hero/`
- [x] Define types in `types.go`
- [x] Create `New()` function in `model.go`
- [x] Add basic `Init()`, `Update()`, `View()` skeleton

### Part 2: Core Rendering ✅ COMPLETE
- [x] Implement lane rendering in `view.go`
- [x] Create lane styles in `styles.go`
- [x] Render hit zone at bottom
- [x] Test: See empty lanes on screen

### Part 3: Animation Loop ✅ COMPLETE
- [x] Add `TickMsg` and `tickCmd()` in `update.go`
- [x] Implement `updateNotePositions()` in `notes.go`
- [x] Create test notes that scroll down
- [x] Test: Notes fall from top to bottom

### Part 4: Input Handling ✅ COMPLETE
- [x] Implement key press detection in `update_keyboard.go`
- [x] Map A/S/D/F/SPACE to lanes 0-4 (improved from original J key)
- [x] Call `checkHit()` on key press
- [x] Test: Can press keys and see response

### Part 5: Hit Detection ✅ COMPLETE
- [x] Implement `checkHit()` in `scoring.go`
- [x] Calculate timing windows (Perfect/Good/OK/Miss)
- [x] Mark notes as hit when successful
- [x] Show visual feedback on hit
- [x] Test: Hit notes at right time → score increases

### Part 6: Scoring System ✅ COMPLETE
- [x] Implement combo counter
- [x] Calculate multiplier (2x/3x/4x)
- [x] Show score, combo, multiplier on screen
- [x] Reset combo on miss
- [x] Test: Build combo → see multiplier increase

### Part 7: Song System ✅ COMPLETE
- [x] Create demo songs in `songs.go` (3 difficulties: Easy/Medium/Hard)
- [x] Implement `loadSong()` function
- [x] Add song selection menu
- [x] Test: Load song → notes appear at right times

### Part 8: Polish ✅ COMPLETE
- [x] Add visual hit feedback (color-coded by result)
- [x] Show "PERFECT!" / "GOOD" / "OK" / "MISS" messages
- [x] Add end-of-song results screen
- [x] Show final score and max combo
- [x] Add restart option

---

## 🎯 Success Criteria ✅ ALL MET

When complete, you should be able to:
1. ✅ See 5 lanes with notes scrolling down at 60 FPS
2. ✅ Press A/S/D/F/SPACE to hit notes (ergonomic home row + thumb)
3. ✅ See Perfect/Good/OK/Miss feedback with color coding
4. ✅ Build combos and see multiplier increase (2x/3x/4x)
5. ✅ Play through a complete song with 3 difficulty levels
6. ✅ See final score and max combo on results screen

---

## 📝 Post-Implementation Notes

### What Was Built

**Core Game** (893 LOC across 9 files):
- `types.go` - Game state, note timing, constants (91 lines)
- `model.go` - Initialization and song end detection (66 lines)
- `update.go` - Main game loop with 60 FPS animation (62 lines)
- `update_keyboard.go` - Input handling for all game states (95 lines)
- `view.go` - Lane rendering and UI (228 lines)
- `styles.go` - Neon color palette and text styles (96 lines)
- `notes.go` - Note scrolling and miss detection (64 lines)
- `scoring.go` - Hit detection with timing windows (85 lines)
- `songs.go` - 3 demo songs with varied difficulty (161 lines)

### Critical Bug Fixes

**Issue 1: Hardcoded Coordinate Mismatch** (Similar to Minesweeper Issue 1 & Solitaire Issue 7)

**Problem**: Mixed use of hardcoded coordinates (Y=25) and dynamic calculations created hit detection failures.

**Files affected**:
- `scoring.go:16` - Hardcoded `hitZoneY := 25`
- `notes.go:18` - Hardcoded `m.notes[i].Y = 25 - int(distanceFromHit)`
- `notes.go:29` - Magic number `if note.Y <= 28`
- `model.go:42` - Dynamic `m.hitZoneY = m.termHeight - 10` (unused!)

**Solution**: Added constants to `types.go`:
```go
const (
    NoteAreaHeight = 25  // Fixed rendering height
    GracePeriod    = 3   // Extra rows for missed notes
)
```

**Key Lesson**: Following CLAUDE.md pattern - coordinate calculations MUST match rendering logic exactly. Use constants for fixed layout dimensions, not dynamic values based on terminal size.

### User Experience Improvements

**1. Ergonomic Controls** (update_keyboard.go:56)
- **Original**: A/S/D/F/J (requires 5 fingers spread across keyboard)
- **Improved**: A/S/D/F/SPACE (home row + thumb = much more comfortable!)
- **Impact**: Game is significantly easier to play at high speeds

**2. Clear Visual Feedback** (view.go:158)
- **Original**: Generic `[#]` symbols in hit zone
- **Improved**: Shows actual keys `[A] [S] [D] [F] [SPC]` in lane colors
- **Impact**: Immediately obvious which key to press for each lane

**3. Centered Results Screen** (styles.go:46-54, view.go:203-212)
- **Original**: Left-aligned text on completion screen
- **Improved**: All text centered and properly styled
- **Impact**: Professional, polished end-game experience

**4. Better Game Name**
- **Original**: "Terminal Hero" (implies audio that doesn't exist)
- **Improved**: "Keyboard Hero" (emphasizes keyboard skill)
- **Impact**: Honest branding that sets correct expectations

### Demo Songs

**Easy Street** (Tutorial)
- BPM: 80
- Pattern: Single-lane sequential notes with gaps
- Duration: ~20 seconds
- Perfect for learning controls

**Terminal Jam** (Medium)
- BPM: 140
- Pattern: Wave patterns, rapid center lane, multi-lane sequences
- Duration: ~25 seconds
- Moderate challenge with variety

**Speed Demon** (Hard)
- BPM: 200
- Pattern: Zigzag, simultaneous lanes, dense rapid-fire sections
- Duration: ~30 seconds
- High-speed challenge for experienced players

### Technical Achievements

✅ **Smooth 60 FPS Animation** - Uses 16ms tick intervals
✅ **Precise Timing Windows** - Perfect (1 row), Good (2 rows), OK (3 rows)
✅ **Dynamic Combo System** - Multiplier scales with streak (10/20/30 combo thresholds)
✅ **Color-Coded Feedback** - Green/Yellow/Orange/Red for hit quality
✅ **Clean State Management** - Menu → Playing → Finished flow with restart
✅ **Modular Architecture** - Follows TUITemplate 9-file pattern exactly

### Known Limitations

- **No Audio**: WSL2 environment doesn't support PC speaker beeps
- **Fixed Scroll Speed**: All songs use same speed (configurable but not per-song)
- **No Accuracy %**: Shows max combo but not hit/miss ratio
- **No Leaderboards**: Scores aren't persisted between sessions

### Future Enhancements (from original plan)

These remain good ideas for future iterations:
- **Hold notes** - Long notes you hold down
- **Chord detection** - Hit multiple lanes simultaneously
- **Star power** - Overdrive mode with bonus multiplier
- **Custom charts** - Load songs from JSON files
- **Accuracy stats** - Show perfect/good/ok/miss counts
- **High score persistence** - Save best scores per song

---

## 🚀 How to Play

```bash
cd ~/projects/TUIClassics

# Build the game
make classics

# Launch TUI Classics
./bin/classics

# Press 'h' to launch Keyboard Hero
```

### Controls

**Menu**:
- `↑/↓` or `j/k` - Select song
- `Enter` - Start selected song
- `q` - Quit

**Gameplay**:
- `A` - Hit lane 1 (Green)
- `S` - Hit lane 2 (Red)
- `D` - Hit lane 3 (Yellow)
- `F` - Hit lane 4 (Blue)
- `SPACE` - Hit lane 5 (Magenta)
- `ESC` - Return to menu
- `q` - Quit

**Results Screen**:
- `Enter` - Play again
- `ESC` - Return to menu
- `q` - Quit

### Tips for High Scores

1. **Home Row Position** - Keep fingers on ASDF and thumb on spacebar
2. **Watch the Colors** - Each lane has a unique color to track notes
3. **Timing is Everything** - Hit when notes reach the `[KEY]` boxes at bottom
4. **Build Combos** - Perfect hits build your multiplier (2x/3x/4x)
5. **Start Easy** - Master "Easy Street" before attempting "Speed Demon"

---

## 💡 Implementation Lessons Learned

1. **Coordinate Consistency** - Always use constants for fixed layout dimensions, not dynamic values
2. **Test Early** - Smooth scrolling and hit detection need testing from the start
3. **User Feedback** - Ergonomic controls (ASDF+SPACE) make a huge difference in playability
4. **Visual Clarity** - Showing hotkeys in hit zones eliminates confusion
5. **Pattern Variety** - Demo songs with different difficulties provide good progression

---

## 🔜 Potential Future Enhancements

Ideas that could be added in future iterations:

**Gameplay Features**:
- **Hold notes** - Long notes you hold down for duration
- **Chord detection** - Hit multiple lanes simultaneously
- **Star power** - Overdrive mode with bonus multiplier
- **Variable scroll speed** - Per-song difficulty adjustment
- **Accuracy percentage** - Show hit/miss ratio and grade (S/A/B/C)

**Content**:
- **Custom charts** - Load songs from JSON files
- **More songs** - Expand beyond 3 demo songs
- **Song editor** - Create custom note charts
- **Difficulty modes** - Same song, different note density

**Persistence**:
- **Leaderboards** - Save high scores per song
- **Progress tracking** - Unlock harder songs
- **Statistics** - Track total notes hit, best combo, etc.

**Technical**:
- **Audio integration** - Background music playback (mpv/beep library)
- **PC speaker beeps** - Sound effects for hits (Linux only)
- **Replay system** - Save and playback perfect runs

---

## 🎉 Project Complete!

**Keyboard Hero** is fully implemented and playable!

**Final Stats**:
- ✅ 893 lines of code across 9 modular files
- ✅ All success criteria met
- ✅ Critical coordinate bugs fixed
- ✅ UX improvements based on real testing
- ✅ 3 demo songs with varied difficulty
- ✅ Smooth 60 FPS gameplay
- ✅ Complete menu → play → results flow

**Play it now**: `./bin/classics` → Press `h` → Pick a song → Rock out! ⌨️🎸
