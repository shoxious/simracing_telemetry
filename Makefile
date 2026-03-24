.PHONY: dev frontend backend windows clean help

# ── Development ───────────────────────────────────────────────────────────────
# Run Go backend (simulate) and Nuxt dev server side-by-side
dev:
	@echo "Starting Go backend (simulate mode) on :8080..."
	@cd backend && go run . --simulate &
	@echo "Starting Nuxt dev server on :3000..."
	@cd frontend && npm run dev

# Run Go backend only (with simulation)
backend-dev:
	cd backend && go run . --simulate

# Run Nuxt dev server only (connect to existing backend)
frontend-dev:
	cd frontend && npm run dev

# ── Build ─────────────────────────────────────────────────────────────────────
frontend:
	cd frontend && npm install && npx nuxi generate

backend-mac: frontend
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../simracing-dashboard .

windows: frontend
	cd backend && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../simracing-dashboard.exe .

# ── Utility ───────────────────────────────────────────────────────────────────
clean:
	rm -rf frontend/.nuxt frontend/.output backend/static/* simracing-dashboard simracing-dashboard.exe

go-tidy:
	cd backend && go mod tidy

help:
	@echo ""
	@echo "SimRacing Dashboard - Make targets:"
	@echo "  make dev          - Start backend + frontend in dev mode"
	@echo "  make backend-dev  - Start Go backend only (simulate mode)"
	@echo "  make frontend-dev - Start Nuxt dev server only"
	@echo "  make windows      - Build Windows EXE (cross-compile from macOS)"
	@echo "  make backend-mac  - Build macOS binary"
	@echo "  make clean        - Remove build artifacts"
	@echo ""
