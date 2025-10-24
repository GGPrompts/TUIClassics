# 🚀 Post-Lunch Quick Start

Welcome back! Here's exactly what to do to get started.

---

## Step 1: Verify You're on the Right Branch (5 seconds)

```bash
cd ~/projects/TUIClassics
git status
```

Should show: `On branch multiplayer`

---

## Step 2: Set Up Supabase (15 minutes)

Follow **SUPABASE_SETUP.md** step-by-step:

1. Create account at supabase.com
2. Create project named "tuiarcade"
3. Copy API keys
4. Create `.env` file
5. Run database schema in SQL Editor
6. Verify tables exist

**Don't skip this!** Everything else depends on Supabase being set up.

---

## Step 3: Install Go Dependencies (30 seconds)

```bash
go get github.com/supabase-community/supabase-go
go get github.com/gorilla/websocket
```

---

## Step 4: Build Supabase Package (Claude - 30 minutes)

Ask Claude:

```
"Create the internal/supabase package with the following:

1. client.go - Initialize Supabase client using SUPABASE_URL and SUPABASE_ANON_KEY from .env
2. players.go - GetOrCreatePlayer(username) function that returns player UUID
3. scores.go - SyncScore() and FetchLeaderboard() functions

Use the supabase-go library. Reference the database schema from MULTIPLAYER.md Section 1.2."
```

Expected files created:
- `internal/supabase/client.go`
- `internal/supabase/players.go`
- `internal/supabase/scores.go`

---

## Step 5: Test with Minesweeper (Claude - 30 minutes)

Ask Claude:

```
"Integrate Supabase with the Minesweeper game:

1. When a game is won, sync the score to Supabase (async, non-blocking)
2. Modify view_stats.go to show both local scores AND global leaderboard
3. Handle offline gracefully - if Supabase is unavailable, just skip sync

Files to modify:
- games/minesweeper/model.go (add sync on win)
- games/minesweeper/view_stats.go (add global leaderboard)

Reference the existing stats package in internal/stats/ for how scores are currently saved locally."
```

---

## Step 6: Test It! (5 minutes)

```bash
# Build minesweeper
make minesweeper

# Play a game
./bin/minesweeper

# Win a game (Easy mode is fastest)
# Check the stats view (press 'S' from menu)
# You should see:
#   - Your local best times (left side)
#   - Global leaderboard (right side)
```

---

## Step 7: Verify in Supabase (2 minutes)

1. Go to your Supabase dashboard
2. Click **Table Editor**
3. Click **scores** table
4. You should see your score entry!

```
game         | category | score | created_at
-------------|----------|-------|------------------
minesweeper  | easy     | 42    | 2025-10-24 15:30
```

---

## Step 8: Expand to Other Games (Claude - 45 minutes)

Ask Claude:

```
"Apply the same Supabase integration to these games:
- 2048 (games/2048/)
- Snake (games/snake/)
- Solitaire (games/solitaire/)
- Hero (games/hero/)

For each game:
1. Add sync on game end
2. Update stats view to show global leaderboard

Reuse the pattern from Minesweeper."
```

---

## Step 9: Build Lobby Infrastructure (Claude - 1 hour)

Ask Claude:

```
"Create the multiplayer lobby infrastructure:

1. internal/network/protocol.go - Define message types for lobby/games
2. internal/network/server.go - WebSocket server with room management
3. internal/network/client.go - WebSocket client for connecting to server

Reference the README files in each directory for architecture details.

Start simple - just lobby join, chat, and room listing. No games yet."
```

---

## Step 10: Build Lobby UI (Claude - 45 minutes)

Ask Claude:

```
"Create the lobby client UI in games/lobby/:

Layout:
- Left side: Game browser (list available games)
- Right top: Chat window
- Right bottom: Online players list

Features:
- Connect to WebSocket server
- Display available games
- Send/receive chat messages
- Show connected players

Use Bubbletea. Reference games/menu/ for similar UI patterns."
```

---

## Step 11: Create Launcher (Claude - 15 minutes)

Ask Claude:

```
"Create cmd/arcade/main.go - the multiplayer client launcher.

On launch:
1. Check for username in ~/.config/tuiarcade/config.json
2. Prompt for username if first time
3. Connect to localhost:8080
4. Launch lobby UI
5. Handle connection errors gracefully

Keep it minimal - just launch the lobby."
```

---

## Step 12: First Multiplayer Game (Claude - 1 hour)

Ask Claude:

```
"Implement Tic-Tac-Toe as the first multiplayer game:

Server-side (add to network/server.go):
- Track 3x3 board state
- Validate moves
- Check win conditions
- Broadcast state to both players

Client-side (create games/tictactoe/):
- Render 3x3 board
- Click to place X/O
- Receive game state updates
- Show winner

Make it playable between two clients!"
```

---

## Progress Checkpoints

After each step, verify it works before moving on:

- [ ] Step 2: Supabase tables visible in dashboard
- [ ] Step 4: Supabase package compiles
- [ ] Step 6: Minesweeper score appears in Supabase
- [ ] Step 6: Global leaderboard shows in stats view
- [ ] Step 8: All games sync scores
- [ ] Step 9: WebSocket server starts
- [ ] Step 10: Lobby UI renders
- [ ] Step 11: Can launch arcade client
- [ ] Step 12: Two clients can play tic-tac-toe!

---

## If You Get Stuck

### Supabase connection errors
- Verify `.env` file exists with correct keys
- Check `SUPABASE_URL` has no trailing slash
- Verify you ran the database schema SQL

### WebSocket connection fails
- Make sure server is running: `./bin/arcade-server`
- Check firewall isn't blocking port 8080
- Verify client connects to `ws://localhost:8080` (not `wss://`)

### Compilation errors
- Run `go mod tidy` to sync dependencies
- Check imports are correct
- Verify all files are on `multiplayer` branch

---

## Estimated Timeline

| Phase | Time | Claude % |
|-------|------|----------|
| Steps 1-3 (Setup) | 20 min | 0% |
| Step 4 (Supabase pkg) | 30 min | 10% |
| Step 5-6 (Minesweeper) | 35 min | 10% |
| Step 7-8 (Other games) | 45 min | 10% |
| Step 9 (Network) | 60 min | 15% |
| Step 10 (Lobby UI) | 45 min | 10% |
| Step 11 (Launcher) | 15 min | 5% |
| Step 12 (Tic-tac-toe) | 60 min | 15% |
| **Total** | **~5 hours** | **~75%** |

You have 50% usage - prioritize Steps 1-10 today, save Step 12 for when you want to test multiplayer end-to-end.

---

## Alternative: Shorter Session (2-3 hours)

If you want a quicker win:

**Just do Steps 1-6** (Supabase + Minesweeper leaderboard)
- Takes ~1.5 hours
- Uses ~20% Claude
- Gets global leaderboards working
- Immediate visible progress

Then continue multiplayer tomorrow with full Claude budget.

---

## End of Day Goal

**Minimum viable success:**
- ✅ Minesweeper has global leaderboard
- ✅ Scores sync to Supabase
- ✅ Can see other players' scores

**Stretch goal:**
- ✅ All games have leaderboards
- ✅ Lobby server runs
- ✅ Lobby client connects

**Dream goal:**
- ✅ Can play tic-tac-toe multiplayer!

---

**Ready to build! Enjoy your lunch, then let's make this happen. 🚀**
