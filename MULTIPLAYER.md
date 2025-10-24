# TUI Arcade - Multiplayer Implementation Plan

> **Branch**: `multiplayer`
> **Status**: 🚧 In Development
> **Goal**: Transform TUIClassics into a unified gaming platform with both single-player and multiplayer games, powered by Supabase

---

## Vision

A terminal-based gaming platform that combines:
- **Single-player games** (Minesweeper, 2048, Snake, Solitaire, Hero) with global leaderboards
- **Multiplayer games** (Texas Hold'em, Chess, Battleship, Tic-Tac-Toe, etc.) with lobby system
- **Unified backend** (One Supabase database for everything)
- **Dual launchers** (`classics` for solo, `arcade` for multiplayer)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────┐
│  TUI ARCADE - Unified Gaming Platform          │
├─────────────────────────────────────────────────┤
│                                                 │
│  SINGLE-PLAYER (cmd/classics)                   │
│  ├─ Games: Minesweeper, 2048, Snake, etc.      │
│  ├─ Local scores (JSON)                        │
│  └─ Global leaderboards (Supabase)             │
│                                                 │
│  MULTIPLAYER (cmd/arcade)                       │
│  ├─ Lobby system (WebSocket)                   │
│  ├─ Games: Hold'em, Chess, Battleship, etc.    │
│  └─ Real-time sync (Supabase)                  │
│                                                 │
│  BACKEND (Supabase)                             │
│  ├─ Players table (unified identity)           │
│  ├─ Scores table (global leaderboards)         │
│  ├─ Games table (active multiplayer sessions)  │
│  └─ Chat table (lobby + game chat)             │
└─────────────────────────────────────────────────┘
```

---

## Phase 1: Supabase Foundation (Day 1 - Session 1)

**Time estimate**: 1-2 hours
**Claude usage**: 10-15%

### 1.1 Supabase Setup (Manual - 15 minutes)

**Steps:**
1. Go to https://supabase.com
2. Create free account
3. Create new project: "tuiarcade"
4. Wait for database provisioning (~2 minutes)
5. Copy API keys from Settings > API

**Save these values:**
```bash
SUPABASE_URL=https://YOUR_PROJECT.supabase.co
SUPABASE_ANON_KEY=your_anon_key_here
```

### 1.2 Database Schema (Run in Supabase SQL Editor)

```sql
-- Enable UUID extension
create extension if not exists "uuid-ossp";

-- Players table (unified identity for all games)
create table players (
  id uuid primary key default uuid_generate_v4(),
  username text unique not null,
  created_at timestamp with time zone default now(),
  last_seen timestamp with time zone default now()
);

-- Scores table (global leaderboards for single-player games)
create table scores (
  id uuid primary key default uuid_generate_v4(),
  player_id uuid references players(id) on delete cascade,
  game text not null,              -- 'minesweeper', 'snake', '2048', etc.
  category text,                   -- 'easy', 'medium', 'hard' (for difficulty-based games)
  score int not null,              -- points or time in seconds (lower is better for time-based)
  metadata jsonb,                  -- { "moves": 42, "tiles": [2048, 1024], "time": 120 }
  created_at timestamp with time zone default now()
);

-- Create index for fast leaderboard queries
create index idx_scores_leaderboard on scores(game, category, score desc);
create index idx_scores_player on scores(player_id, created_at desc);

-- Multiplayer games table (active game sessions)
create table games (
  id uuid primary key default uuid_generate_v4(),
  game_type text not null,         -- 'holdem', 'chess', 'battleship', etc.
  host_id uuid references players(id) on delete cascade,
  status text not null default 'waiting',  -- 'waiting', 'playing', 'finished'
  max_players int not null,
  current_players int not null default 1,
  game_state jsonb,                -- Game-specific state
  created_at timestamp with time zone default now(),
  started_at timestamp with time zone,
  finished_at timestamp with time zone
);

-- Game participants (many-to-many)
create table game_participants (
  game_id uuid references games(id) on delete cascade,
  player_id uuid references players(id) on delete cascade,
  joined_at timestamp with time zone default now(),
  left_at timestamp with time zone,
  primary key (game_id, player_id)
);

-- Chat messages (lobby + in-game)
create table chat_messages (
  id uuid primary key default uuid_generate_v4(),
  game_id uuid references games(id) on delete cascade,  -- null for lobby chat
  player_id uuid references players(id) on delete cascade,
  message text not null,
  created_at timestamp with time zone default now()
);

-- Achievements (optional - can add later)
create table achievements (
  id uuid primary key default uuid_generate_v4(),
  player_id uuid references players(id) on delete cascade,
  game text not null,
  achievement_type text not null,  -- 'first_win', 'speed_demon', 'perfect_game', etc.
  unlocked_at timestamp with time zone default now(),
  unique(player_id, game, achievement_type)
);

-- Enable Row Level Security (optional - for future multi-tenancy)
alter table players enable row level security;
alter table scores enable row level security;
alter table games enable row level security;
alter table chat_messages enable row level security;

-- Create public read policies (everyone can see scores/games)
create policy "Public can view scores" on scores for select using (true);
create policy "Public can view games" on games for select using (true);
create policy "Public can view chat" on chat_messages for select using (true);

-- Players can insert their own scores (TODO: add auth later)
create policy "Players can insert scores" on scores for insert with check (true);
create policy "Players can create games" on games for insert with check (true);
create policy "Players can send chat" on chat_messages for insert with check (true);
```

### 1.3 Code Structure (Create directories)

```
internal/
├── supabase/
│   ├── client.go       (Supabase client initialization)
│   ├── players.go      (Player management)
│   ├── scores.go       (Leaderboard operations)
│   ├── games.go        (Multiplayer game sessions)
│   └── chat.go         (Chat operations)
└── network/
    ├── server.go       (WebSocket server for multiplayer)
    ├── client.go       (WebSocket client)
    └── protocol.go     (Message types)

cmd/
├── arcade/
│   └── main.go         (Multiplayer lobby launcher)
└── arcade-server/
    └── main.go         (Self-hosted server option)

games/
├── lobby/              (Lobby UI)
│   ├── model.go
│   ├── update.go
│   ├── view.go
│   └── styles.go
└── holdem/             (First multiplayer game)
    ├── model.go
    ├── update.go
    ├── view.go
    └── poker.go        (Poker logic)
```

### 1.4 Dependencies

```bash
go get github.com/supabase-community/supabase-go
go get github.com/gorilla/websocket
```

---

## Phase 2: Supabase Integration (Day 1 - Session 2)

**Time estimate**: 1 hour
**Claude usage**: 10%

### 2.1 Create `internal/supabase` Package

**Goal**: Reusable Supabase client for all games

**Files to create:**
- `internal/supabase/client.go` - Initialize client, config
- `internal/supabase/players.go` - GetOrCreatePlayer, UpdateLastSeen
- `internal/supabase/scores.go` - SyncScore, FetchLeaderboard

### 2.2 Integrate with Minesweeper (Test case)

**Modify:**
- `games/minesweeper/model.go` - Add Supabase sync on win
- `games/minesweeper/view_stats.go` - Show local + global leaderboards

**Test:**
- Play a game
- Win and verify score syncs to Supabase
- View stats and see global leaderboard

---

## Phase 3: Expand to All Single-Player Games (Day 1 - Session 3)

**Time estimate**: 1 hour
**Claude usage**: 10%

### 3.1 Add Supabase to Remaining Games

**Games to update:**
- 2048
- Snake
- Solitaire
- Hero

**Pattern (same for all):**
```go
// In each game's model.go - on game end:
if gameWon {
    // 1. Save local (instant)
    localScores, _ := stats.Load()
    stats.Update[Game]Score(localScores, score)
    stats.Save(localScores)

    // 2. Sync to Supabase (async)
    go supabase.SyncScore("game_name", category, score, metadata)
}
```

### 3.2 Update Stats Views

Each game's `view_stats.go`:
- Fetch global leaderboard
- Display side-by-side with personal bests
- Cache results (don't fetch every render)

---

## Phase 4: Multiplayer Lobby Infrastructure (Day 1 - Session 4)

**Time estimate**: 2-3 hours
**Claude usage**: 15-20%

### 4.1 WebSocket Server

**Create**: `internal/network/server.go`

Features:
- Accept WebSocket connections
- Room management (create/join/leave)
- Message routing (chat, game events)
- Player sessions

### 4.2 WebSocket Client

**Create**: `internal/network/client.go`

Features:
- Connect to server
- Send/receive messages
- Reconnection logic
- Heartbeat

### 4.3 Protocol Definition

**Create**: `internal/network/protocol.go`

```go
type MessageType string

const (
    MsgJoinLobby    MessageType = "join_lobby"
    MsgCreateGame   MessageType = "create_game"
    MsgJoinGame     MessageType = "join_game"
    MsgLeaveGame    MessageType = "leave_game"
    MsgGameAction   MessageType = "game_action"
    MsgChatMessage  MessageType = "chat"
    MsgGameState    MessageType = "game_state"
)

type Message struct {
    Type    MessageType     `json:"type"`
    Payload json.RawMessage `json:"payload"`
}
```

### 4.4 Lobby UI

**Create**: `games/lobby/`

Views:
- Game browser (list available games)
- Chat window
- Player list
- Create game dialog

---

## Phase 5: First Multiplayer Game - Tic-Tac-Toe (Day 1 - Session 5)

**Time estimate**: 1 hour
**Claude usage**: 10%

### 5.1 Why Tic-Tac-Toe First?

- Simplest multiplayer logic
- Validates entire flow
- 3x3 grid (you've done grids!)
- Turn-based (easy sync)
- Can implement in <1 hour

### 5.2 Implementation

**Create**: `games/tictactoe/`

**Server-side (in lobby server):**
- Track board state (3x3 array)
- Validate moves
- Check win conditions
- Broadcast state to both players

**Client-side:**
- Render board
- Send moves to server
- Receive updates
- Display winner

---

## Phase 6: Additional Multiplayer Games (Day 2+)

### Quick Wins (30 min each)
- **Rock Paper Scissors** - Simultaneous input, instant result
- **Connect Four** - Reuse tic-tac-toe, bigger grid, gravity
- **Wordle Race** - Both players same word, first to solve wins

### Medium Complexity (2-3 hours each)
- **Battleship** - Grid placement, hidden state, turn-based shooting
- **Checkers** - Movement logic, jump validation, piece promotion
- **Blackjack** - Card rendering from Solitaire, dealer AI

### Flagship Games (4-6 hours each)
- **Texas Hold'em** - Full poker logic, betting rounds, pot management
- **Chess** - Complete chess rules, move validation, check/checkmate
- **Snake Arena** - Multiplayer snake, real-time, collision detection

---

## Testing Checklist

### Supabase Integration
- [ ] Scores sync to database
- [ ] Leaderboards load correctly
- [ ] Works offline (degrades gracefully)
- [ ] No duplicate scores
- [ ] Handles network errors

### Multiplayer Lobby
- [ ] Can connect to server
- [ ] See other players
- [ ] Chat works
- [ ] Create/join games
- [ ] Leave games cleanly

### Multiplayer Games
- [ ] Game state syncs correctly
- [ ] No race conditions
- [ ] Handles disconnections
- [ ] Spectators work (optional)
- [ ] Reconnection recovers state

---

## Branch Strategy

### Main Branch
- Stable single-player games
- Can release anytime
- No multiplayer code

### Multiplayer Branch
- All features from main (merge regularly)
- Plus multiplayer infrastructure
- Plus Supabase integration
- Experimental/WIP

### Merging Back
```bash
# Regular sync (main -> multiplayer)
git checkout multiplayer
git merge main

# When Supabase is stable (multiplayer -> main)
git checkout main
git cherry-pick <supabase-commits>  # Just leaderboards

# When multiplayer is production-ready
git checkout main
git merge multiplayer
```

---

## Configuration

### Environment Variables

Create `.env` (gitignored):
```bash
SUPABASE_URL=https://YOUR_PROJECT.supabase.co
SUPABASE_ANON_KEY=your_anon_key
ARCADE_SERVER_URL=ws://localhost:8080  # or wss://arcade.example.com
```

### User Config

`~/.config/tuiarcade/config.json`:
```json
{
  "username": "Matt",
  "player_id": "uuid-here",
  "preferences": {
    "enable_global_leaderboards": true,
    "enable_multiplayer": true,
    "default_server": "localhost:8080"
  }
}
```

---

## Deployment Options

### Option 1: Central Server (Recommended for MVP)
- Deploy lobby server to Fly.io/Railway
- Free tier: $0-5/month
- All players connect to same server
- Easy matchmaking

### Option 2: Self-Hosted (P2P Alternative)
- Users run `arcade-server` locally
- Share IP/port with friends
- Zero hosting cost
- NAT traversal required

### Option 3: Hybrid
- Discovery server on Fly.io
- Game servers self-hosted
- Best of both worlds

---

## Success Metrics

### Phase 1 Complete When:
- [x] Supabase database created
- [ ] Minesweeper syncs scores
- [ ] Global leaderboard visible in-game

### Phase 2 Complete When:
- [ ] All single-player games have leaderboards
- [ ] Offline mode works
- [ ] User has username/profile

### Phase 3 Complete When:
- [ ] Lobby server running
- [ ] Client can connect
- [ ] Chat works
- [ ] Tic-tac-toe playable

### Phase 4 Complete When:
- [ ] 3+ multiplayer games
- [ ] Stable networking
- [ ] Can host private games

---

## Post-Lunch Quick Start Guide

When you return from lunch:

### 1. Verify branch
```bash
git status
# Should show: On branch multiplayer
```

### 2. Set up Supabase (15 minutes)
- Follow Section 1.1 above
- Run SQL from Section 1.2
- Save API keys

### 3. Ask Claude to:
```
"Create the internal/supabase package with client initialization and score sync functionality"
```

### 4. Test with Minesweeper:
```
"Integrate Supabase with Minesweeper - sync scores on win and show global leaderboard"
```

### 5. Iterate from there!

---

**Ready to build something awesome! 🚀**

---

## Resources

- [Supabase Docs](https://supabase.com/docs)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Bubbletea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [TUITemplate Architecture](https://github.com/GGPrompts/TUITemplate/blob/main/ARCHITECTURE.md)

---

**Last Updated**: 2025-10-24 (Day 1 - Pre-lunch setup)
