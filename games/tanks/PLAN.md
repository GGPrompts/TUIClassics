# Tanks (Artillery Game) - Implementation Plan

## Game Overview

A turn-based artillery game inspired by **Gorillas.bas** and **Scorched Earth**. Two players control tanks on opposite sides of a destructible cityscape, taking turns adjusting angle and velocity to launch projectiles at each other.

## Visual Style

**Muted Palette with Strategic Color:**
- Buildings: Grayscale/dark gray blocks with lighter gray borders
- Sky/Background: Dark background (consistent with neon theme)
- Tanks: Dark brown/gray bodies with **bright brown/orange borders**
- Projectiles: Gray body with **bright yellow borders** (high visibility)
- Projectile trail: Faint yellow dotted line (optional)
- Explosion: **Bright yellow/orange burst effect** (visual highlight)
- Craters: Darker gray/black (shows damage)
- UI elements: Subtle cyan or white text
- Active player indicator: Bright cyan/white pulsing border on tank

## Game Mechanics

### Core Gameplay Loop
1. Player 1's turn: Adjust angle (0-90°) and velocity
2. Launch projectile with physics-based arc
3. Projectile travels through air, potentially destroying buildings
4. Hit detection: Check for tank hit, building collision, or ground impact
5. Explosion animation with terrain destruction
6. Switch to Player 2
7. Repeat until one tank is destroyed

### Physics
- **Projectile motion**: Standard kinematic equations
- Horizontal velocity: `vx = v * cos(angle)`
- Vertical velocity: `vy = v * sin(angle)`
- Gravity: Constant downward acceleration
- Collision detection: Check each frame for terrain/tank intersection

### Win Conditions
- Direct hit on enemy tank (health-based or instant kill)
- Best of 3/5/7 rounds
- Score tracking across multiple rounds

## Technical Architecture

Following the TUIClassics modular pattern:

### File Structure
```
tanks/
├── main.go              - Entry point (minimal)
├── types.go             - Type definitions
├── styles.go            - Lipgloss styling
├── model.go             - Model initialization
├── update.go            - Update dispatcher
├── update_keyboard.go   - Keyboard handling
├── view.go              - View dispatcher
├── render_game.go       - Main game rendering
├── render_ui.go         - UI elements (angle, velocity, scores)
├── physics.go           - Projectile physics & collision
├── terrain.go           - Terrain generation & destruction
├── sprites.go           - Tank ASCII art at different angles
└── animation.go         - Explosion & projectile animations
```

### Core Types

```go
type model struct {
    // Game state
    gameState    gameState    // menu, playing, gameOver
    currentTurn  int          // 0 = player1, 1 = player2
    winner       int          // -1 = none, 0 = player1, 1 = player2

    // Players
    player1      player
    player2      player

    // Terrain
    buildings    []building
    groundLevel  int

    // Projectile
    projectile   *projectile  // nil when not in flight
    projectileFrames []projectileFrame  // For animation replay

    // Input
    angleInput   string
    velocityInput string
    inputMode    inputMode    // angle, velocity, none

    // Animation
    explosionFrame int
    explosionPos   position
    animating      bool

    // Display
    width        int
    height       int
}

type player struct {
    position     position
    angle        float64      // 0-90 degrees
    velocity     float64      // Launch power
    health       int          // Optional: multi-hit gameplay
    score        int          // Rounds won
    color        lipgloss.Color
}

type building struct {
    x            int
    height       int
    width        int
    cells        [][]bool     // For destruction: true = intact, false = destroyed
}

type projectile struct {
    x            float64
    y            float64
    vx           float64      // Horizontal velocity
    vy           float64      // Vertical velocity
    trail        []position   // For visual trail
}

type position struct {
    x            int
    y            int
}

type gameState int
const (
    stateMenu gameState = iota
    statePlaying
    stateAnimating
    stateGameOver
)

type inputMode int
const (
    inputNone inputMode = iota
    inputAngle
    inputVelocity
)
```

## Implementation Phases

### Phase 1: Basic Structure (Foundation)
- [ ] Create file structure with all modules
- [ ] Define types in `types.go`
- [ ] Create basic lipgloss styles in `styles.go`
- [ ] Implement model initialization in `model.go`
- [ ] Set up update dispatcher in `update.go`
- [ ] Create basic view dispatcher in `view.go`

### Phase 2: Terrain Generation
- [ ] Implement random building generation in `terrain.go`
- [ ] Create building rendering function
- [ ] Generate varied skyline (different heights/widths)
- [ ] Position tanks on tallest buildings on each side
- [ ] Add ground rendering

### Phase 3: Tank Sprites
- [ ] Design tank ASCII art for 9-11 angles (0°, 15°, 30°, 45°, 60°, 75°, 90°)
- [ ] Create sprite map in `sprites.go`
- [ ] Implement sprite selection based on angle
- [ ] Add tank rendering with borders
- [ ] Add turret/cannon visualization
- [ ] Mirror sprites for right-facing tank

### Phase 4: Input System
- [ ] Implement angle input in `update_keyboard.go`
- [ ] Implement velocity input
- [ ] Add input validation (angle: 0-90, velocity: reasonable range)
- [ ] Create UI for displaying current angle/velocity
- [ ] Add arrow keys for fine-tuning angle/velocity
- [ ] Add Enter to confirm and launch

### Phase 5: Physics Engine
- [ ] Implement projectile motion equations in `physics.go`
- [ ] Create projectile update function (position per frame)
- [ ] Add gravity constant
- [ ] Implement velocity decomposition (vx, vy from angle/power)
- [ ] Test projectile arcs with different inputs

### Phase 6: Collision Detection
- [ ] Implement ground collision detection
- [ ] Implement building collision detection
- [ ] Implement tank collision detection
- [ ] Add terrain destruction on building hit
- [ ] Create crater/hole in buildings when hit

### Phase 7: Projectile Rendering
- [ ] Render projectile in flight with yellow border
- [ ] Add projectile trail effect (fading yellow dots)
- [ ] Smooth animation (update position each tick)
- [ ] Clear trail after impact

### Phase 8: Explosion Animation
- [ ] Design explosion frames in `animation.go`
- [ ] Create expanding burst effect (yellow/orange)
- [ ] Add particle-like characters radiating outward
- [ ] Implement frame-by-frame animation
- [ ] Add slight screen shake effect (optional)

### Phase 9: Game Loop
- [ ] Implement turn switching
- [ ] Add turn indicator (highlight active tank)
- [ ] Reset input fields after launch
- [ ] Add game over detection
- [ ] Implement win screen

### Phase 10: UI Polish
- [ ] Create header with scores in `render_ui.go`
- [ ] Add "Player 1's Turn" / "Player 2's Turn" indicator
- [ ] Display current angle and velocity during input
- [ ] Add help text footer (controls)
- [ ] Add launch power visualization (bar/meter)
- [ ] Create game over screen with winner announcement

### Phase 11: Advanced Features (Optional)
- [ ] Add wind (affects horizontal velocity)
- [ ] Multiple weapon types (different explosion radii)
- [ ] Tank movement (limited distance per turn)
- [ ] Power-ups scattered on buildings
- [ ] AI opponent for single-player mode
- [ ] Best of X rounds with score tracking
- [ ] Replay system (show last shot)
- [ ] Settings menu (gravity, wind, round count)

### Phase 12: Menu System
- [ ] Create main menu (New Game, Settings, Quit)
- [ ] Add pause menu (Resume, Restart, Main Menu)
- [ ] Settings screen (rounds to win, starting health)

## Controls

### During Input Phase
- **Arrow Up/Down**: Adjust angle (+/- 5°)
- **Arrow Left/Right**: Adjust velocity (+/- 5)
- **A**: Enter angle input mode (type exact value)
- **V**: Enter velocity input mode (type exact value)
- **Space/Enter**: Launch projectile
- **ESC**: Cancel input, return to menu
- **Q**: Quit game

### During Animation
- **No input** - wait for animation to complete

### Menu Navigation
- **Arrow Up/Down**: Select menu option
- **Enter**: Confirm selection
- **ESC**: Go back / Quit

## Technical Considerations

### Performance
- Target 30-60 FPS for smooth projectile animation
- Optimize collision detection (only check visible terrain)
- Cache rendered building sprites to avoid recalculation

### Terminal Compatibility
- Use standard Unicode box-drawing characters
- Ensure colors work on dark backgrounds
- Test on different terminal sizes (min: 80x24)

### Animation Timing
- Use Bubbletea's `tick` command for frame updates
- Projectile speed: Balance between realistic and playable
- Explosion duration: ~500ms (15-20 frames at 30 FPS)

## Testing Checklist

- [ ] Projectile arcs correctly at all angles (0°, 45°, 90°)
- [ ] Collisions detect properly (buildings, ground, tanks)
- [ ] Buildings show destruction/craters after hits
- [ ] Turn switching works correctly
- [ ] Input validation prevents invalid values
- [ ] Game over triggers on tank hit
- [ ] UI displays correct information
- [ ] Colors render properly on dark background
- [ ] No crashes on edge cases (projectile off-screen, etc.)

## Future Expansion Ideas

- **Gorillas Mode**: Swap tank sprites for gorilla ASCII art, bananas for projectiles
- **Theme Toggle**: Switch between muted palette and full neon
- **Map Editor**: Custom building layouts
- **Network Multiplayer**: Play over LAN/internet
- **Tournament Mode**: Bracket-style multiple players
- **Leaderboard**: Track win/loss ratios, best shots

## Notes

- Keep sprites simple: Wide rectangle base, square turret, rotatable cannon
- Explosion is the visual centerpiece - make it satisfying!
- Consider adding sound effects later (beep library)
- Maintain modular architecture for easy theme swapping (tanks ↔ gorillas)
- Follow TUIClassics patterns from minesweeper/solitaire

## References

- Original Gorillas.bas (MS-DOS QBASIC)
- Scorched Earth gameplay
- TUIClassics/games/minesweeper for architectural reference
- TUIClassics/games/solitaire for card rendering patterns
