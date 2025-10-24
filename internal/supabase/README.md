# internal/supabase

Supabase client and operations for TUI Arcade backend.

## Purpose

Unified cloud backend for:
- Global leaderboards (single-player games)
- Multiplayer game state
- Player profiles
- Chat messages

## Files to Create

- `client.go` - Initialize Supabase client, load config
- `players.go` - GetOrCreatePlayer(), UpdateLastSeen()
- `scores.go` - SyncScore(), FetchLeaderboard()
- `games.go` - CreateGame(), JoinGame(), UpdateGameState()
- `chat.go` - SendMessage(), FetchMessages()

## Usage Example

```go
// Initialize once at app start
err := supabase.Init()

// In game logic
supabase.SyncScore("minesweeper", "easy", 42, metadata)

// Fetch leaderboard
scores, _ := supabase.FetchLeaderboard("minesweeper", "easy", 10)
```

## Configuration

Requires environment variables:
- `SUPABASE_URL`
- `SUPABASE_ANON_KEY`

See `SUPABASE_SETUP.md` in project root.
