# cmd/arcade

Main launcher for TUI Arcade multiplayer client.

## Purpose

Entry point for multiplayer games. Connects to lobby server and launches the multiplayer UI.

## Usage

```bash
# Build
make arcade

# Run - connect to default server
./bin/arcade

# Connect to specific server
./bin/arcade --server ws://localhost:8080

# Connect to cloud server
./bin/arcade --server wss://arcade.example.com
```

## Flow

1. Launch `arcade`
2. Check for saved username, prompt if needed
3. Connect to lobby server (default: localhost:8080)
4. Show lobby UI
5. User browses/joins games
6. Launch game UI when joined
7. Return to lobby when game ends

## Configuration

Reads from `~/.config/tuiarcade/config.json`:

```json
{
  "username": "Matt",
  "player_id": "uuid-here",
  "default_server": "localhost:8080"
}
```

## Minimal main.go

```go
package main

import (
    "github.com/GGPrompts/TUIClassics/games/lobby"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    m := lobby.New()
    p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
    p.Run()
}
```
