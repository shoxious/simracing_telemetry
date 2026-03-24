#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# SimRacing Dashboard – Build Script
# Usage:
#   ./build.sh              → local dev binary (simulate mode only on macOS)
#   ./build.sh --windows    → cross-compile Windows EXE from macOS
# ─────────────────────────────────────────────────────────────────────────────
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$SCRIPT_DIR/frontend"
BACKEND_DIR="$SCRIPT_DIR/backend"
STATIC_DIR="$BACKEND_DIR/static"

TARGET="${1:-}"

echo ""
echo "╔════════════════════════════════════════╗"
echo "║     SimRacing Dashboard – Build        ║"
echo "╚════════════════════════════════════════╝"
echo ""

# ── Step 1: Install frontend dependencies ────────────────────────────────────
echo "▶ Step 1: Installing Nuxt dependencies..."
cd "$FRONTEND_DIR"
npm install --silent
echo "  ✓ Done"

# ── Step 2: Build Nuxt static SPA ─────────────────────────────────────────────
echo ""
echo "▶ Step 2: Building Nuxt frontend (static mode)..."
npx nuxi generate
echo "  ✓ Output → backend/static/"

# Verify output
if [ ! -f "$STATIC_DIR/index.html" ]; then
  echo "  ✗ ERROR: index.html not found in backend/static/"
  exit 1
fi

# ── Step 3: Download Go dependencies ──────────────────────────────────────────
echo ""
echo "▶ Step 3: Downloading Go modules..."
cd "$BACKEND_DIR"
go mod download
echo "  ✓ Done"

# ── Step 4: Compile Go binary ──────────────────────────────────────────────────
echo ""
if [ "$TARGET" = "--windows" ]; then
  echo "▶ Step 4: Cross-compiling for Windows (amd64)..."
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o "$SCRIPT_DIR/simracing-dashboard.exe" .
  echo "  ✓ Built: simracing-dashboard.exe"
  echo ""
  echo "  Transfer simracing-dashboard.exe to your Windows machine"
  echo "  and run: simracing-dashboard.exe"
  echo "  (make sure iRacing is running first)"
else
  echo "▶ Step 4: Building for host platform (macOS/Linux)..."
  CGO_ENABLED=0 \
    go build -ldflags="-s -w" \
    -o "$SCRIPT_DIR/simracing-dashboard" .
  echo "  ✓ Built: simracing-dashboard"
  echo ""
  echo "  Run with simulation mode:"
  echo "  ./simracing-dashboard --simulate"
fi

echo ""
echo "╔════════════════════════════════════════╗"
echo "║  Build complete!                       ║"
echo "╚════════════════════════════════════════╝"
