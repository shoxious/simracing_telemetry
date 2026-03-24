# SimRacing Live Dashboard

Echtzeit iRacing Telemetrie-Dashboard – läuft als **einzelne Windows EXE** die Backend, WebSocket-Server und das gesamte Web-Frontend enthält.

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

## Inhaltsverzeichnis

1. [Voraussetzungen](#1-voraussetzungen)
2. [Projekt klonen / herunterladen](#2-projekt-klonen--herunterladen)
3. [Lokale Entwicklung (macOS / Linux)](#3-lokale-entwicklung-macos--linux)
4. [Windows EXE bauen](#4-windows-exe-bauen)
5. [Produktionsbetrieb auf Windows](#5-produktionsbetrieb-auf-windows)
6. [Alle Start-Optionen (CLI Flags)](#6-alle-start-optionen-cli-flags)
7. [Projektstruktur](#7-projektstruktur)
8. [Seiten & Features](#8-seiten--features)
9. [Datenspeicher-Strategie](#9-datenspeicher-strategie)
10. [Häufige Probleme](#10-häufige-probleme)

---

## 1. Voraussetzungen

### Pflicht

| Tool | Version | Download |
|------|---------|----------|
| **Go** | ≥ 1.22 | https://go.dev/dl/ |
| **Node.js** | ≥ 18 | https://nodejs.org/ |
| **npm** | ≥ 9 | (kommt mit Node.js) |

### Nur für Windows EXE (Production)

- Windows 10/11 (64-bit) zum Ausführen
- iRacing muss auf demselben Windows-PC installiert sein

### Installation prüfen

```bash
go version     # go version go1.22.x ...
node --version # v18.x.x oder höher
npm --version  # 9.x.x oder höher
```

---

## 2. Projekt klonen / herunterladen

```bash
git clone https://github.com/dein-user/simracing-dashboard.git
cd simracing-dashboard
```

Oder als ZIP herunterladen und entpacken.

---

## 3. Lokale Entwicklung (macOS / Linux)

Da iRacing nur unter Windows läuft, gibt es einen **Simulation-Modus** der realistische GT3-Telemetrie synthetisch erzeugt. Damit lässt sich das gesamte Dashboard auf macOS entwickeln und testen.

### Option A – Alles auf einmal (empfohlen)

```bash
make dev
```

Startet:
- Go Backend auf `http://localhost:8080` (Simulation)
- Nuxt Dev-Server auf `http://localhost:3000` (Hot-Reload)

> Der Nuxt Dev-Server verbindet sich automatisch mit dem Go Backend auf Port 8080 für WebSocket-Daten.

### Option B – Schritt für Schritt

**Schritt 1 – Frontend-Abhängigkeiten installieren:**

```bash
cd frontend
npm install
```

**Schritt 2 – Backend starten (Simulation):**

```bash
cd backend
go run . --simulate
```

Das Backend läuft jetzt auf `http://localhost:8080`.

**Schritt 3 – Frontend Dev-Server starten (neues Terminal):**

```bash
cd frontend
npm run dev
```

Das Frontend läuft jetzt auf `http://localhost:3000` mit Hot-Reload.

**Schritt 4 – Browser öffnen:**

```
http://localhost:3000
```

---

## 4. Windows EXE bauen

Die EXE enthält Frontend + Backend in einer einzigen Datei (~15–20 MB).

### Von macOS / Linux aus cross-compilieren

```bash
# Alles in einem Schritt:
./build.sh --windows

# Oder manuell:
cd frontend
npm install
npx nuxi generate           # Baut das Frontend nach backend/static/

cd ../backend
go mod download             # Go-Abhängigkeiten laden
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
  go build -ldflags="-s -w" -o ../simracing-dashboard.exe .
```

**Ergebnis:** `simracing-dashboard.exe` im Projektroot.

> **Wichtig:** `CGO_ENABLED=0` ist zwingend nötig – nur so funktioniert die Cross-Kompilierung, weil `modernc.org/sqlite` reines Go ist und kein CGo benötigt.

### Auf Windows direkt bauen

```bat
build.bat
```

Das Skript führt automatisch aus:
1. `npm install` + `npx nuxi generate` (Frontend)
2. `go mod download` (Go-Module)
3. `go build` (EXE kompilieren)

---

## 5. Produktionsbetrieb auf Windows

**Voraussetzung:** iRacing ist gestartet und die Strecke geladen.

### Starten

```bat
simracing-dashboard.exe
```

Der Browser öffnet sich automatisch auf `http://localhost:8080`.
Die EXE gibt im Terminalfenster folgendes aus:

```
╔════════════════════════════════════════╗
║     SimRacing Live Dashboard v1.0      ║
╚════════════════════════════════════════╝
► Mode: LIVE (iRacing shared memory)
► Address: http://localhost:8080
► Database: simracing.db
► Dashboard running. Press Ctrl+C to stop.
```

### Vom Handy / Tablet zugreifen

1. PC und Gerät im selben WLAN
2. Windows-Firewall: Port 8080 freigeben (einmalig)
3. Auf dem Handy aufrufen: `http://[PC-IP-Adresse]:8080`

Die IP-Adresse des PCs herausfinden:
```bat
ipconfig
```
→ `IPv4-Adresse` unter dem WLAN-Adapter.

### Mit eigenem Port starten

```bat
simracing-dashboard.exe --addr :9090
```

→ Dashboard erreichbar unter `http://localhost:9090`

### Als Windows-Dienst (Autostart)

Mit dem kostenlosen Tool **NSSM** (Non-Sucking Service Manager):

```bat
nssm install SimRacingDashboard "C:\Pfad\simracing-dashboard.exe"
nssm set SimRacingDashboard AppParameters "--addr :8080"
nssm start SimRacingDashboard
```

---

## 6. Alle Start-Optionen (CLI Flags)

```
simracing-dashboard.exe [Optionen]

Optionen:
  --addr string     HTTP-Adresse und Port (Standard: ":8080")
  --simulate        Simulation-Modus – kein iRacing benötigt
  --db string       Pfad zur SQLite-Datenbank (Standard: "simracing.db")
  --open            Browser automatisch öffnen (Standard: true)
  --open=false      Browser NICHT automatisch öffnen
```

**Beispiele:**

```bat
# Standard (iRacing live)
simracing-dashboard.exe

# Simulation zum Testen
simracing-dashboard.exe --simulate

# Anderen Port verwenden
simracing-dashboard.exe --addr :9090

# Eigene Datenbank-Datei
simracing-dashboard.exe --db "C:\iRacing\session_2026.db"

# Kein automatischer Browser-Start
simracing-dashboard.exe --open=false

# Alles kombiniert
simracing-dashboard.exe --simulate --addr :9090 --db test.db
```

---

## 7. Projektstruktur

```
simracing/
│
├── backend/                     # Go Backend
│   ├── main.go                  # Einstiegspunkt, embedded Frontend, Server-Start
│   ├── go.mod                   # Go-Abhängigkeiten
│   │
│   ├── irsdk/                   # iRacing SDK Integration
│   │   ├── types.go             # Shared-Memory-Structs (C-kompatibel)
│   │   ├── reader_windows.go    # Echte iRacing-Daten (nur Windows)
│   │   ├── reader_other.go      # Stub für macOS/Linux
│   │   └── simulator.go         # Synthetische GT3-Telemetrie (Monza)
│   │
│   ├── hub/
│   │   └── hub.go               # WebSocket Broadcast Hub
│   │
│   ├── storage/
│   │   ├── ring.go              # In-Memory Ring Buffer (120s @ 60Hz)
│   │   └── sqlite.go            # SQLite Downsampler + Lap Records
│   │
│   ├── api/
│   │   ├── server.go            # Gin Router, SPA Fallback
│   │   └── handlers.go          # WebSocket + REST Endpoints
│   │
│   └── static/                  # ← Nuxt Build-Output (wird embedded)
│       ├── index.html
│       ├── _nuxt/
│       └── ...
│
├── frontend/                    # Nuxt 4 Frontend
│   ├── nuxt.config.ts           # SPA-Modus, Output → backend/static/
│   ├── tailwind.config.ts       # Racing-Dark-Theme
│   ├── package.json
│   │
│   └── app/                     # Nuxt 4 App-Verzeichnis
│       ├── app.vue
│       ├── pages/
│       │   ├── index.vue        # Dashboard (Speed, RPM, Fuel, Tires)
│       │   ├── telemetry.vue    # Detailierte Traces + Datentabelle
│       │   ├── timing.vue       # Renntabelle + Rundenhistorie
│       │   └── strategy.vue     # Kraftstoffstrategie + Pit-Window
│       │
│       ├── components/
│       │   ├── gauges/          # SVG-Rundanzeigen (Speed, RPM, Fuel)
│       │   ├── charts/          # Canvas Rolling Chart (Throttle/Brake)
│       │   ├── tires/           # Reifentemperatur-Widget
│       │   ├── timing/          # Rundenzeit + Runden-Tabelle
│       │   └── inputs/          # Pedal + Lenkrad Anzeige
│       │
│       ├── composables/
│       │   ├── useWebSocket.ts  # WS-Verbindung mit Auto-Reconnect
│       │   └── useIRacing.ts    # Reaktive Telemetrie-Daten
│       │
│       ├── stores/
│       │   └── iracing.ts       # Pinia Store (Telemetrie, Session, Status)
│       │
│       └── assets/css/
│           └── main.css         # Tailwind + Racing-Dark-Theme
│
├── build.sh                     # Build-Script (macOS/Linux)
├── build.bat                    # Build-Script (Windows)
├── Makefile                     # Make-Targets für Entwicklung
└── README.md
```

---

## 8. Seiten & Features

### Dashboard (`/`)

Die Hauptseite mit allen wichtigen Werten auf einen Blick:

- **Speed-Gauge** – SVG-Rundanzeige 0–350 km/h (Farbzonen: blau → gelb → rot)
- **RPM-Gauge** – Drehzahl mit 5 Schaltlichtern und Redline-Flash
- **Pedal-Bars** – Throttle (grün), Brake (rot), Clutch (blau), Lenkwinkel
- **Rundenzeit** – Aktuelle Runde, letzte Runde, beste Runde + Delta
- **Kraftstoff** – Füllstand mit Runden-Reichweite und Verbrauch/Stunde
- **Reifentemperaturen** – 4 Reifen × 3 Zonen (innen/mitte/außen), farbkodiert

### Telemetry (`/telemetry`)

- **Rolling Chart** – Canvas-basierter 30-Sekunden-Trace (Throttle/Brake/Speed), wählbar 10s/30s/60s
- **Detail-Gauges** – Speed, RPM, Gang, Position als Zahlendisplay
- **Live-Datentabelle** – Alle 16+ Telemetrie-Werte in Echtzeit

### Timing (`/timing`)

- **Renntabelle** – Alle Fahrzeuge sortiert nach Position mit Gap-Anzeige
- **Eigenes Auto** hervorgehoben (blau)
- **Rundenhistorie** – Tabelle aller gefahrenen Runden aus SQLite

### Strategy (`/strategy`)

- **Pit-Window-Warnung** – Blinkt wenn ≤ 3 Runden Kraftstoff verbleiben
- **Kraftstoff-Gauge** – Großformatig mit Runden-Reichweite
- **Fuel-per-Lap Diagramm** – Canvas Bar-Chart aller vergangenen Runden
- **Parameter-Override** – Tankgröße und Soll-Rundenzeit anpassbar

---

## 9. Datenspeicher-Strategie

iRacing liefert bis zu 100+ Variablen bei 60 Hz. Naives Speichern wäre ~86 MB/Stunde. Das Dashboard verwendet eine zweistufige Strategie:

```
iRacing (60 Hz)
      │
      ▼
Ring Buffer          ← In-Memory, letzte 120 Sekunden bei 60 Hz (~5 MB)
(7.200 Frames)         Für: Live WebSocket, /api/history, /api/telemetry/latest
      │
      │ 1× pro Sekunde
      ▼
SQLite (1 Hz)        ← ~3 MB/Stunde für 20 Key-Variablen
  telemetry_1hz        Für: Rundenhistorie, Strategy-Charts, /api/laps
  laps
  sessions
```

**SQLite-Optimierungen:**
- WAL-Modus (kein Write-Lock beim Lesen)
- `synchronous=NORMAL` (Kompromiss aus Sicherheit und Speed)
- 64 MB Page-Cache im RAM
- Automatische Lap-Erkennung (Runden werden beim Lap-Wechsel gespeichert)

Die Datenbank `simracing.db` liegt neben der EXE und wächst pro Stunde um ~3–5 MB.

---

## 10. Häufige Probleme

### Backend / Go

**`iRacing not running (OpenFileMapping)`**
→ iRacing muss gestartet sein und eine Session geladen haben (Strecke befahren).
→ Zum Testen: `--simulate` Flag verwenden.

**`go: command not found`**
→ Go ist nicht installiert. Download: https://go.dev/dl/
→ Nach Installation Terminal neu starten oder PATH prüfen: `export PATH=$PATH:/usr/local/go/bin`

**Port 8080 bereits belegt**
→ Anderen Port verwenden: `simracing-dashboard.exe --addr :9090`
→ Oder belegen des Ports prüfen: `netstat -ano | findstr :8080` (Windows)

**`embedded static files not found`**
→ Das Frontend wurde noch nicht gebaut. Ausführen:
```bash
cd frontend && npx nuxi generate
```
→ Dann Backend neu bauen/starten.

**Windows Firewall blockiert Zugriff vom Handy**
→ Firewall-Regel hinzufügen (einmalig, als Administrator):
```bat
netsh advfirewall firewall add rule name="SimRacing Dashboard" dir=in action=allow protocol=TCP localport=8080
```

### Frontend / Node

**`nuxi: command not found`**
→ Dependencies noch nicht installiert:
```bash
cd frontend && npm install
```

**Seite lädt, aber keine Daten (grauer Screen)**
→ WebSocket-Verbindung prüfen: Browser DevTools → Network → WS
→ Backend läuft? `http://localhost:8080/api/status` im Browser aufrufen
→ Im Dev-Modus: Backend muss auf Port 8080 laufen, Nuxt auf 3000

**Hot-Reload funktioniert nicht**
→ `npm run dev` neu starten
→ `.nuxt` Verzeichnis löschen: `rm -rf frontend/.nuxt`

### Simulation

**Simulation sieht unrealistisch aus**
→ Der Simulator modelliert einen fiktiven GT3-Kurs (~88s Rundenzeit).
→ Reifentemperaturen brauchen 3 Runden um auf Betriebstemperatur zu kommen (realistisch).
→ Kraftstoffverbrauch: ~2.8L/Runde bei einem vollen Tank von 55L.

---

## Schnellreferenz

```bash
# Entwicklung (macOS, Simulation)
make dev

# Nur Backend
cd backend && go run . --simulate

# Nur Frontend Dev-Server
cd frontend && npm run dev

# Frontend bauen (Ausgabe → backend/static/)
cd frontend && npx nuxi generate

# Windows EXE von macOS
./build.sh --windows

# Windows EXE auf Windows
build.bat

# Datenbank zurücksetzen
rm simracing.db

# Go-Module aufräumen
cd backend && go mod tidy
```

---

## Technologie-Stack

| Komponente | Technologie | Grund |
|-----------|-------------|-------|
| Backend Language | Go 1.22 | Single-binary EXE, natives Windows |
| HTTP/WS Server | Gin + gorilla/websocket | Performance, stabile APIs |
| SQLite | modernc.org/sqlite | **Reines Go** – ermöglicht CGO_ENABLED=0 für Cross-Compilation |
| iRacing SDK | Windows Shared Memory | Offizielle Methode, `golang.org/x/sys/windows` |
| Frontend | Nuxt 4 (Vue 3, SPA) | SSR=false → statisch embedded im Go-Binary |
| State Management | Pinia | Vue 3 Standard |
| Styling | Tailwind CSS | Mobile-first, kein Build-Overhead |
| Gauges | SVG (custom) | Keine externe Abhängigkeit, animierbar |
| Telemetry Chart | HTML5 Canvas | 60 FPS ohne Framework-Overhead |
