# games/lobby

Multiplayer lobby interface for TUI Arcade.

## Purpose

Main menu for multiplayer games:
- Browse available games
- Create new game
- Join existing games
- Chat with other players
- View online players

## Layout

```
┌─ TUI ARCADE ────────────────┐  ┌─ CHAT ───────────┐
│                              │  │ <User1> hi      │
│ 🃏 Card Games                │  │ <User2> gg      │
│  ├─ Texas Hold'em    [23]   │  │ * Player joined │
│  └─ Blackjack        [12]   │  └─────────────────┘
│                              │
│ 🎯 Strategy                  │  ┌─ ONLINE ────────┐
│  ├─ Chess            [8]    │  │ 🟢 Matt         │
│  └─ Battleship       [15]   │  │ 🟢 Sarah        │
│                              │  │ 🟢 Dan          │
│ 🎮 Arcade                    │  └─────────────────┘
│  └─ Snake Arena      [31]   │
└──────────────────────────────┘
```

## Files to Create

- `model.go` - Lobby state, connection, game list
- `update.go` - Message handler (server events)
- `update_keyboard.go` - Navigation, chat input
- `view.go` - Render lobby UI
- `view_chat.go` - Chat window component
- `view_games.go` - Game browser component
- `styles.go` - Lipgloss styles

## Features

- Real-time game list updates
- Live chat
- Create game with settings
- Join game by clicking
- See player count for each game
- Filter games by type
