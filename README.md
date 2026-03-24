# SimRacing Live Dashboard

Real-time iRacing telemetry dashboard for desktop, tablet and mobile — served as a **single self-contained Windows EXE**.

```
┌─────────────────────────────────────────────────────┐
│              simracing-dashboard.exe                │
│                                                     │
│  iRacing SDK ──► WebSocket (60 Hz) ──► Browser     │
│  Shared Memory    REST API /api/*     Nuxt 4 SPA    │
│                   SQLite Storage      Mobile-First  │
└─────────────────────────────────────────────────────┘
```

---

## Features

| Page | What you get |
|------|-------------|
| **Dashboard** | Speed & RPM gauges, pedal trace, lap timer, fuel level, tyre temperatures |
| **Telemetry** | 30-second rolling throttle/brake/speed chart, live data table |
| **Timing** | Race standings for all cars, gap to leader, full lap history |
| **Strategy** | Fuel per lap chart, pit window alert, configurable tank size |

- 60 Hz WebSocket stream directly from iRacing shared memory
- Simulation mode – full dashboard works without iRacing (synthetic GT3 data)
- Mobile & tablet optimised (responsive Tailwind layout, safe-area aware)
- Auto-reconnect with exponential backoff
- SQLite storage – ~3–5 MB/hour, ring-buffer in RAM for live data

---

## Requirements

| Tool | Version |
|------|---------|
| [Go](https://go.dev/dl/) | ≥ 1.22 |
| [Node.js](https://nodejs.org/) | ≥ 18 |
| npm | ≥ 9 |

To **run** the EXE: Windows 10/11 (64-bit) with iRacing installed.

---

## Quick Start

### 1 – Clone

```bash
git clone https://github.com/your-user/simracing-dashboard.git
cd simracing-dashboard
```

### 2 – Build

**macOS / Linux → Windows EXE:**

```bash
./build.sh --windows
```

**Windows native:**

```bat
build.bat
```

Both scripts install npm deps, generate the Nuxt SPA, and compile the Go binary.

### 3 – Run

```bat
simracing-dashboard.exe
```

The browser opens automatically at `http://localhost:8080`.

**No iRacing? Use simulation mode:**

```bat
simracing-dashboard.exe --simulate
```

---

## Accessing from Phone / Tablet

1. PC and device on the same Wi-Fi
2. Allow port 8080 through Windows Firewall (one-time):
   ```bat
   netsh advfirewall firewall add rule name="SimRacing Dashboard" dir=in action=allow protocol=TCP localport=8080
   ```
3. Open on your device: `http://<PC-IP>:8080`
   Find your PC's IP with `ipconfig` → look for **IPv4 Address**.

---

## CLI Flags

```
simracing-dashboard.exe [options]

  --addr string    Listen address  (default: ":8080")
  --simulate       Synthetic telemetry, no iRacing required
  --db   string    SQLite database path  (default: "simracing.db")
  --open=false     Do not open browser automatically
```

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22 – single binary, zero runtime deps |
| WebSocket | gorilla/websocket |
| HTTP API | Gin |
| Database | modernc.org/sqlite (pure Go, no CGo) |
| Frontend | Nuxt 4 · Vue 3 · Pinia · Tailwind CSS |
| Gauges | Custom SVG |
| Chart | HTML5 Canvas (60 FPS) |
| iRacing SDK | Windows shared memory (`golang.org/x/sys/windows`) |

---

## License

MIT
