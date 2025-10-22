# Terminal Hero - Implementation Plan

**Type**: Rhythm Game (Guitar Hero style)
**Complexity**: Medium (~900 LOC)
**Time Estimate**: 6-8 hours
**Status**: Ready to implement

---

## 🎸 Game Overview

A Guitar Hero style rhythm game for the terminal where notes fall down lanes and you press keys at the right time.

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

### Part 1: Setup
- [ ] Create all 9 files in `games/hero/`
- [ ] Define types in `types.go`
- [ ] Create `New()` function in `model.go`
- [ ] Add basic `Init()`, `Update()`, `View()` skeleton

### Part 2: Core Rendering
- [ ] Implement lane rendering in `view.go`
- [ ] Create lane styles in `styles.go`
- [ ] Render hit zone at bottom
- [ ] Test: See empty lanes on screen

### Part 3: Animation Loop
- [ ] Add `TickMsg` and `tickCmd()` in `update.go`
- [ ] Implement `updateNotePositions()` in `notes.go`
- [ ] Create test notes that scroll down
- [ ] Test: Notes fall from top to bottom

### Part 4: Input Handling
- [ ] Implement key press detection in `update_keyboard.go`
- [ ] Map A/S/D/F/J to lanes 0-4
- [ ] Call `checkHit()` on key press
- [ ] Test: Can press keys and see response

### Part 5: Hit Detection
- [ ] Implement `checkHit()` in `scoring.go`
- [ ] Calculate timing windows (Perfect/Good/OK/Miss)
- [ ] Mark notes as hit when successful
- [ ] Show visual feedback on hit
- [ ] Test: Hit notes at right time → score increases

### Part 6: Scoring System
- [ ] Implement combo counter
- [ ] Calculate multiplier (2x/3x/4x)
- [ ] Show score, combo, multiplier on screen
- [ ] Reset combo on miss
- [ ] Test: Build combo → see multiplier increase

### Part 7: Song System
- [ ] Create demo songs in `songs.go`
- [ ] Implement `loadSong()` function
- [ ] Add song selection menu
- [ ] Test: Load song → notes appear at right times

### Part 8: Polish
- [ ] Add visual hit feedback (flash on Perfect)
- [ ] Show "PERFECT!" / "GOOD" / "MISS" messages
- [ ] Add end-of-song results screen
- [ ] Show accuracy percentage
- [ ] Add restart option

---

## 🎯 Success Criteria

When complete, you should be able to:
1. ✅ See 5 lanes with notes scrolling down
2. ✅ Press A/S/D/F/J to hit notes
3. ✅ See Perfect/Good/OK/Miss feedback
4. ✅ Build combos and see multiplier increase
5. ✅ Play through a complete song
6. ✅ See final score and accuracy

---

## 🚀 Quick Start

```bash
cd ~/projects/TUIClassics

# Create the game files
cd games/hero
touch types.go model.go update.go update_keyboard.go view.go styles.go notes.go scoring.go songs.go

# Start implementing!
# 1. types.go - Define GameState, Note, Model
# 2. model.go - Implement New() function
# 3. update.go - Add Init(), Update(), tickCmd()
# 4. view.go - Render lanes
# 5. notes.go - Scroll notes down
# 6. scoring.go - Hit detection
# 7. songs.go - Load demo songs

# Register in menu
# Add to games/menu/model.go

# Build and test
cd ../..
make classics
./bin/classics
# Press 'h' for Hero
```

---

## 💡 Pro Tips

1. **Start simple** - Get lanes rendering first, then add notes
2. **Test scrolling** - Make sure note positions update smoothly
3. **Timing is critical** - Use consistent time calculations
4. **Visual feedback** - Flash colors make hits feel satisfying
5. **Demo songs** - Start with slow, simple patterns

---

## 🎵 Song Ideas

**Easy (Tutorial)**:
- Single lane patterns
- Slow BPM (60-80)
- Long gaps between notes

**Medium**:
- 2-3 lane patterns
- Medium BPM (100-140)
- Some rapid sections

**Hard**:
- All 5 lanes
- Fast BPM (160-200)
- Dense note patterns
- Chord hits (multiple lanes at once)

---

## 🔜 Future Enhancements

- **Hold notes** - Long notes you hold down
- **Chord detection** - Hit multiple lanes simultaneously
- **Star power** - Overdrive mode with bonus multiplier
- **Custom charts** - Load songs from JSON files
- **Audio** - PC speaker beeps (optional)
- **Leaderboards** - Save high scores

---

**Ready to rock! 🎸 Let's build Terminal Hero!**
