# TUI Classics 🎮

> Nostalgic terminal games that bring back the classics - Minesweeper, Solitaire, 2048, and more!

Remember those Windows 95 school computer days? We're bringing that nostalgia to your terminal with beautiful, fully-interactive TUI versions of classic games. Play them standalone or launch from the integrated menu!

## Games

### 💣 Minesweeper
Classic minefield puzzle game with satisfying cascade reveals and explosive animations.

- **Controls**: Mouse click to reveal, right-click to flag, or keyboard navigation
- **Difficulty levels**: Easy (8x8), Medium (16x16), Hard (30x16)
- **Features**: Timer, move counter, explosive animation, high scores
- **Status**: ✅ Complete

### 🃏 Solitaire
Klondike solitaire with the iconic waterfall animation when you win.

- **Controls**: Mouse drag-and-drop or keyboard navigation
- **Features**: Draw-1 and Draw-3 modes, undo, win animation, statistics tracking
- **Status**: ✅ Complete

### 🔢 2048
Slide numbered tiles to combine them and reach the 2048 tile!

- **Controls**: Arrow keys or vim keys (h/j/k/l)
- **Features**: Score tracking, move counter, high scores, smooth animations
- **Status**: ✅ Complete

### 🎸 Keyboard Hero
Rhythm game where you hit keys to the beat - like Guitar Hero for your keyboard!

- **Controls**: Number keys (1-3) to hit notes as they scroll down
- **Features**: Multiple songs, scoring system, combo tracking, high scores
- **Status**: ✅ Complete

### 🐍 Snake
The classic snake game - eat apples, grow longer, don't crash into yourself!

- **Controls**: Arrow keys or vim keys (h/j/k/l)
- **Features**: Progressive speed increase, score tracking, high scores
- **Status**: ✅ Complete

### 🃏 Balatro (Coming Soon!)
Poker roguelike with jokers, planets, and strategic scoring combos.

- **Status**: 🚧 In Development

## Installation

### Quick Start - Launch Menu

The easiest way to play is using the integrated launcher:

```bash
# Clone the repository
git clone https://github.com/GGPrompts/TUIClassics.git
cd TUIClassics

# Build and run the launcher
make classics
./bin/classics
```

The launcher provides a beautiful menu to select and play any game!

### From Source

```bash
# Build all games
make all

# Or install to your Go bin
make install
```

### Individual Games

You can also build and run games individually:

```bash
# Examples:
make minesweeper && ./bin/minesweeper
make solitaire && ./bin/solitaire
make 2048 && ./bin/2048
make snake && ./bin/snake
```

### Using Go Install

```bash
# Install the launcher
go install github.com/GGPrompts/TUIClassics/cmd/classics@latest

# Or install individual games
go install github.com/GGPrompts/TUIClassics/cmd/minesweeper@latest
go install github.com/GGPrompts/TUIClassics/cmd/solitaire@latest
go install github.com/GGPrompts/TUIClassics/cmd/2048@latest
go install github.com/GGPrompts/TUIClassics/cmd/snake@latest
```

## TFE Integration

These games are designed to integrate seamlessly with [TFE (Terminal File Explorer)](https://github.com/GGPrompts/TFE).

When installed, TFE will automatically detect the games and add them to the context menu, allowing you to take quick mental breaks while working!

## Features

✅ **Beautiful TUI rendering** - Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
✅ **Integrated launcher menu** - Select games from a beautiful animated menu
✅ **Mouse and keyboard support** - Play however you prefer
✅ **High score tracking** - Persistent leaderboards across all games
✅ **Statistics tracking** - Track wins, losses, streaks, and more
✅ **Satisfying animations** - Waterfall wins, cascade reveals, explosions, rhythm visualizations
✅ **Cross-platform** - Works on Linux, macOS, and Windows
✅ **TFE Integration** - Launch games directly from your terminal file explorer

## Development

Built using the modular architecture from [TUITemplate](https://github.com/GGPrompts/TUITemplate), each game follows clean separation:

- `cmd/` - Minimal main.go entry points
- `games/` - Game logic and rendering
- `games/shared/` - Shared utilities (high scores, themes, animations)

### Project Structure

```
TUIClassics/
├── cmd/
│   ├── classics/        # Launcher menu for all games
│   ├── minesweeper/     # Minesweeper binary
│   ├── solitaire/       # Solitaire binary
│   ├── 2048/            # 2048 binary
│   └── snake/           # Snake binary
├── games/
│   ├── menu/            # Main launcher menu
│   ├── minesweeper/     # Minesweeper game logic
│   ├── solitaire/       # Solitaire game logic
│   ├── 2048/            # 2048 game logic
│   ├── snake/           # Snake game logic
│   ├── hero/            # Keyboard Hero game logic
│   ├── balatro/         # Balatro game logic (WIP)
│   └── shared/          # Shared code (scores, stats, themes)
└── CLAUDE.md            # AI development notes & patterns
```

### Running Tests

```bash
make test
```

### Building

```bash
# Build all games
make all

# Clean and rebuild
make build-all

# Clean build artifacts
make clean
```

## Contributing

Contributions are welcome! Whether it's:

- 🐛 Bug fixes
- ✨ New game implementations
- 🎨 Theme additions
- 📚 Documentation improvements
- 💡 Feature suggestions

Please open an issue or submit a pull request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Architecture inspired by [TUITemplate](https://github.com/GGPrompts/TUITemplate)
- Nostalgia provided by Windows 95 💾

---

**Made with ❤️ for terminal enthusiasts who remember when these were the only games available at school**
