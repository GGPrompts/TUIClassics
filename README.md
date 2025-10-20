# TUI Classics 🎮

> Nostalgic terminal games that bring back the classics - Minesweeper, Solitaire, and more!

Remember those Windows 95 school computer days? We're bringing that nostalgia to your terminal with beautiful, fully-interactive TUI versions of classic games.

## Games

### 💣 Minesweeper
Classic minefield puzzle game with the satisfying cascade reveals and explosive endings.

- **Controls**: Mouse click to reveal, right-click to flag
- **Difficulty levels**: Easy (8x8), Medium (16x16), Hard (30x16)
- **Features**: Timer, move counter, hint system, high scores

**Status**: 🚧 In Development

### 🎴 Solitaire (Coming Soon!)
Klondike solitaire with the iconic waterfall animation when you win.

- **Controls**: Mouse drag-and-drop or keyboard navigation
- **Features**: Draw-1 and Draw-3 modes, undo, hints, statistics

**Status**: 📋 Planned

### 🐍 Snake (Planned)
The classic snake game - eat apples, grow longer, don't crash!

**Status**: 💡 Future

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/GGPrompts/TUIClassics.git
cd TUIClassics

# Build all games
make all

# Or install to your Go bin
make install
```

### Individual Games

```bash
# Build just minesweeper
make minesweeper
./bin/minesweeper

# Or run directly
make run-minesweeper
```

### Using Go Install

```bash
# Install all games
go install github.com/GGPrompts/TUIClassics/cmd/...@latest

# Or individual games
go install github.com/GGPrompts/TUIClassics/cmd/minesweeper@latest
go install github.com/GGPrompts/TUIClassics/cmd/solitaire@latest
```

## TFE Integration

These games are designed to integrate seamlessly with [TFE (Terminal File Explorer)](https://github.com/GGPrompts/TFE).

When installed, TFE will automatically detect the games and add them to the context menu, allowing you to take quick mental breaks while working!

## Features

✅ **Beautiful TUI rendering** - Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
✅ **Mouse and keyboard support** - Play however you prefer
✅ **High score tracking** - Compete with yourself
✅ **Multiple themes** - Classic, dark, retro, and more
✅ **Satisfying animations** - Waterfall wins, cascade reveals, explosions
✅ **Cross-platform** - Works on Linux, macOS, and Windows

## Development

Built using the modular architecture from [TUITemplate](https://github.com/GGPrompts/TUITemplate), each game follows clean separation:

- `cmd/` - Minimal main.go entry points
- `games/` - Game logic and rendering
- `games/shared/` - Shared utilities (high scores, themes, animations)

### Project Structure

```
TUIClassics/
├── cmd/
│   ├── minesweeper/     # Minesweeper binary
│   ├── solitaire/       # Solitaire binary
│   └── classics/        # Launcher menu for all games
├── games/
│   ├── minesweeper/     # Minesweeper game logic
│   ├── solitaire/       # Solitaire game logic
│   └── shared/          # Shared code (scores, themes, animations)
└── docs/                # Game documentation
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
