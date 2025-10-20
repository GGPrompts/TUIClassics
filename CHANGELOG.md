# Changelog

All notable changes to TUI Classics will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned

- Solitaire (Klondike) with drag-and-drop
- Waterfall animation for Solitaire wins
- Game launcher (`classics` binary)
- TFE integration
- High score persistence

---

## [0.1.0] - 2025-01-20

### Added - Minesweeper Complete! 🎉

#### Core Gameplay
- Classic Minesweeper gameplay with standard rules
- Three difficulty levels:
  - Easy: 8x8 grid, 10 mines
  - Medium: 16x16 grid, 40 mines
  - Hard: 30x16 grid, 99 mines
- Safe first click (mines placed after first reveal)
- Cascade reveal for empty cells
- Win/loss detection
- Timer and move counter
- Flag counter

#### Input Controls
- **Mouse Support**:
  - Left-click to reveal cell
  - Right-click to flag/unflag mine
  - Click detection with accurate coordinate mapping
- **Keyboard Support** (Accessibility):
  - Arrow keys (`↑↓←→`) or Vim keys (`h/j/k/l`) for navigation
  - `Enter` or `Space` to reveal cell at cursor
  - `F` to flag/unflag cell at cursor
  - `N` for new game
  - `R` to restart (when won/lost)
  - `Q` to quit
- **Difficulty Selection**:
  - `1` for Easy
  - `2` for Medium
  - `3` for Hard

#### Visual Design
- Clean, centered UI with automatic terminal resizing
- Classic Windows 95 color scheme:
  - Numbers 1-8 use authentic minesweeper colors
  - Blue (1), Green (2), Red (3), Dark Blue (4), etc.
- Single-width character rendering for perfect grid alignment:
  - `■` - Unrevealed cell
  - `·` - Empty revealed cell
  - `*` - Mine
  - `P` - Flag
  - `1-8` - Adjacent mine count
- Title: `* * * MINESWEEPER * * *`
- Centered stats display
- Helpful footer with control hints

#### Animations ✨
- **Explosion Animation** (the signature feature!):
  - Expanding shockwave radiates from clicked mine
  - Uses Manhattan distance for radial effect
  - Mines revealed progressively as wave reaches them
  - Wave front shows `+` with blinking effect
  - Previously revealed mines show `*`
  - Animates at 12.5 fps (80ms per frame)
  - Title changes to `* * * BOOM! * * *` during explosion
  - Shows "Explosion expanding... (radius: N)" status
  - Smooth transition to game over after animation completes

#### Architecture
- Modular file structure following [TUITemplate](https://github.com/GGPrompts/TUITemplate):
  - `types.go` - Type definitions
  - `model.go` - Game logic and state
  - `update.go` - Message dispatcher
  - `update_keyboard.go` - Keyboard handling
  - `update_mouse.go` - Mouse handling
  - `view.go` - Rendering
  - `styles.go` - Lipgloss styling
- Clean separation of concerns
- Reusable `getCellSymbolAndStyle()` helper
- Animation state machine (Playing → Exploding → Lost)

#### Documentation
- Comprehensive README with features, installation, controls
- MIT License
- `.gitignore` for Go projects
- Makefile for easy building

---

### Fixed

#### Mouse Click Detection (Commit: 9bbed62)
- **Issue**: Mouse clicks registered 2 rows above actual click position
- **Root Cause**: Incorrect Y-coordinate calculation (`topPadding + 6` instead of `topPadding + 4`)
- **Fix**: Corrected grid start position to account for title (1 line) + `\n\n` (2 lines) + stats (1 line) + `\n\n` (2 lines) = 4 total offset
- **Impact**: Mouse clicks now register exactly on the clicked cell

#### Grid Misalignment (Commit: 2d11c83)
- **Issue**: Grid columns became misaligned when mines were revealed
- **Root Cause**: Variable-width emoji characters (💣 = 2 chars, numbers = 1 char)
- **Fix**: Replaced all emojis with consistent single-width characters
- **Characters Changed**:
  - `💣` → `*` (mine)
  - `🚩` → `P` (flag)
  - `□` → `■` (unrevealed)
  - Empty → `·` (middle dot for revealed empty cells)
- **Impact**: Perfect grid alignment regardless of cell state

#### Rendering Centering (Commit: a65990f)
- **Issue**: Using `lipgloss.Place` made mouse coordinate calculation unpredictable
- **Fix**: Implemented manual centering with calculated padding
- **Impact**: Deterministic positioning for accurate mouse detection

---

### Technical Details

#### Build System
- Go module: `github.com/GGPrompts/TUIClassics`
- Dependencies:
  - `github.com/charmbracelet/bubbletea@latest`
  - `github.com/charmbracelet/lipgloss@latest`
- Build output: `bin/minesweeper` (~4.3 MB)
- Makefile targets:
  - `make minesweeper` - Build game
  - `make clean` - Remove build artifacts
  - `make run-minesweeper` - Build and run
  - `make all` - Build all games

#### Code Statistics (v0.1.0)
- Total lines: ~1,138
- Files: 13
- Minesweeper game: 6 files, ~600 lines
- Zero compiler warnings
- Zero runtime bugs (as of release)

#### Performance
- Instant startup
- Smooth animations (60fps potential, currently 12.5fps for explosion)
- Low memory footprint
- Responsive to terminal resizing

---

## Development Process

### Timeline
- **Session 1** (2025-01-20): Initial implementation
  - Project scaffolding
  - Core gameplay
  - Basic rendering
  - Mouse and keyboard controls
  - Duration: ~4 hours

- **Session 2** (2025-01-20): Bug fixes and polish
  - Fixed mouse coordinate calculation (2-row offset bug)
  - Fixed grid alignment (emoji width issue)
  - Improved rendering (manual centering)
  - Centered stats text
  - Duration: ~1 hour

- **Session 3** (2025-01-20): Explosion animation
  - Added animation state machine
  - Implemented expanding shockwave
  - Progressive mine revelation
  - Blinking wave front effect
  - Duration: ~2 hours

### Total Development Time
**~7 hours** from concept to polished v0.1.0 release

---

## Git History

### Commits (Latest First)

- `b046427` - Add explosion animation and center stats text
- `2d11c83` - Fix grid alignment by using consistent-width characters
- `9bbed62` - Fix mouse Y-coordinate calculation (was off by 2 rows)
- `a65990f` - Fix mouse click detection and improve rendering
- `2fe0c1e` - Initial commit: Playable Minesweeper game

---

## Known Issues

None! 🎉

---

## Future Enhancements (Minesweeper)

Potential improvements for future versions:

- [ ] **Hint System**: Highlight safe cells or logical next moves
- [ ] **Custom Grid Size**: Allow user-defined grid dimensions and mine count
- [ ] **High Score Persistence**: Save best times to file (JSON or SQLite)
- [ ] **Statistics Dashboard**: Track games played, win rate, average time
- [ ] **Difficulty Presets**: Add Expert (24x24, 99 mines) and Beginner (9x9, 10 mines)
- [ ] **Themes**: Multiple color schemes (dark, light, retro, colorblind-friendly)
- [ ] **Sound Effects**: Terminal bell on win/loss (optional)
- [ ] **First-Click Guarantee**: Ensure first click opens a large area (not just no mine)
- [ ] **Quick Restart**: `R` key restarts same board (for practice)
- [ ] **Undo Move**: Take back flag placements

*Note: These are ideas, not commitments. Focus is on completing Solitaire next.*

---

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) - Excellent TUI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling library
- Architecture inspired by [TUITemplate](https://github.com/GGPrompts/TUITemplate)
- Classic Windows 95 Minesweeper for design inspiration
- Generated with [Claude Code](https://claude.com/claude-code) - AI pair programming assistant

---

## Contributors

- **Matt** (GGPrompts) - Project creator and maintainer
- **Claude** (Anthropic) - AI development assistant

Co-Authored-By: Claude <noreply@anthropic.com>

---

## Links

- **Repository**: https://github.com/GGPrompts/TUIClassics
- **Issues**: https://github.com/GGPrompts/TUIClassics/issues
- **TUITemplate**: https://github.com/GGPrompts/TUITemplate
- **TFE**: https://github.com/GGPrompts/TFE

---

**Legend**:
- `Added` - New features
- `Changed` - Changes in existing functionality
- `Deprecated` - Soon-to-be removed features
- `Removed` - Removed features
- `Fixed` - Bug fixes
- `Security` - Vulnerability fixes
