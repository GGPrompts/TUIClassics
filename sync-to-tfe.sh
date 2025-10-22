#!/bin/bash
# Sync TUIClassics games to TFE submodule
# Run this after making changes to games in TUIClassics

set -e  # Exit on error

echo "📦 Building games in TUIClassics..."
cd ~/projects/TUIClassics
go build -o bin/solitaire ./cmd/solitaire
go build -o bin/minesweeper ./cmd/minesweeper
go build -o bin/classics ./cmd/classics

echo "📋 Copying source files to TFE submodule..."
cp -r games/solitaire/*.go ~/projects/TFE/games/TUIClassics/games/solitaire/
cp -r games/minesweeper/*.go ~/projects/TFE/games/TUIClassics/games/minesweeper/
cp -r games/menu/*.go ~/projects/TFE/games/TUIClassics/games/menu/

echo "🔨 Rebuilding games in TFE submodule..."
cd ~/projects/TFE/games/TUIClassics
go build -o bin/solitaire ./cmd/solitaire
go build -o bin/minesweeper ./cmd/minesweeper
go build -o bin/classics ./cmd/classics

echo ""
echo "✅ Games synced successfully!"
echo "🎮 Launch TFE and click Tools → Games Launcher to test"
