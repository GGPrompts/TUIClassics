# internal/network

WebSocket server and client for multiplayer games.

## Purpose

Real-time networking for multiplayer lobby and games.

## Files to Create

- `protocol.go` - Message types and structures
- `server.go` - WebSocket server (for lobby/game hosting)
- `client.go` - WebSocket client (for connecting to games)
- `rooms.go` - Room/lobby management

## Architecture

```
Client 1                    Server                    Client 2
   |                          |                          |
   |--[Connect]-------------->|                          |
   |<--[Welcome]--------------|                          |
   |                          |<--[Connect]--------------|
   |                          |--[Welcome]-------------->|
   |--[CreateGame]----------->|                          |
   |<--[GameCreated]----------|                          |
   |                          |--[GameList]------------->|
   |                          |<--[JoinGame]-------------|
   |<--[PlayerJoined]---------|--[PlayerJoined]-------->|
   |--[GameAction]----------->|                          |
   |                          |--[GameState]------------>|
```

## Message Protocol

All messages are JSON:

```json
{
  "type": "game_action",
  "payload": {
    "action": "move",
    "data": { "x": 5, "y": 3 }
  }
}
```

## Usage Example

### Server
```go
server := network.NewServer(":8080")
server.Start()
```

### Client
```go
client := network.NewClient("ws://localhost:8080")
client.Connect()
client.Send(Message{Type: "join_lobby"})
```
