# cmd/arcade-server

Self-hosted lobby server for TUI Arcade.

## Purpose

Optional self-hosted server for private games with friends.

## Usage

```bash
# Build
make arcade-server

# Run on default port (8080)
./bin/arcade-server

# Run on custom port
./bin/arcade-server --port 9000

# Enable Tailscale auto-detection
./bin/arcade-server --tailscale
```

## Features

- WebSocket server for game lobbies
- Room management
- Chat relay
- Game state synchronization
- Player session tracking

## Deployment Options

### Local (LAN party)
```bash
./bin/arcade-server --port 8080
# Share: Players connect to 192.168.1.x:8080
```

### Tailscale (VPN mesh)
```bash
./bin/arcade-server --tailscale
# Players connect to your-machine:8080
```

### Cloud (Fly.io)
```bash
fly launch
fly deploy
# Players connect to arcade.fly.dev
```

## Architecture

```
arcade-server
├── WebSocket listener (:8080)
├── Room manager (game sessions)
├── Chat relay
└── Supabase sync (optional)
```

## Minimal main.go

```go
package main

import (
    "log"
    "github.com/GGPrompts/TUIClassics/internal/network"
)

func main() {
    server := network.NewServer(":8080")
    log.Println("Lobby server started on :8080")
    server.Start()
}
```

## Configuration

Optional `.env`:
```bash
PORT=8080
ENABLE_SUPABASE=true  # Sync game results
SUPABASE_URL=...
SUPABASE_ANON_KEY=...
```
