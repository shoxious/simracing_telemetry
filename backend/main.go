package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"simracing/api"
	"simracing/hub"
	"simracing/irsdk"
	"simracing/storage"
)

//go:embed static
var staticFiles embed.FS

func main() {
	var (
		addr     = flag.String("addr", ":8080", "HTTP listen address")
		simulate = flag.Bool("simulate", false, "Simulation mode (no iRacing required)")
		dbPath   = flag.String("db", "simracing.db", "SQLite database path")
		open     = flag.Bool("open", true, "Automatically open browser on start")
	)
	flag.Parse()

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║     SimRacing Live Dashboard v1.0      ║")
	fmt.Println("╚════════════════════════════════════════╝")

	if *simulate {
		fmt.Println("► Mode: SIMULATION (synthetic telemetry)")
	} else {
		fmt.Println("► Mode: LIVE (iRacing shared memory)")
	}
	fmt.Printf("► Address: http://localhost%s\n", *addr)
	fmt.Printf("► Database: %s\n", *dbPath)
	fmt.Println()

	// --- Storage ---
	db, err := storage.Open(*dbPath)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer db.Close()

	// 120s ring buffer at 60Hz
	ring := storage.NewRingBuffer(7200)

	// --- WebSocket Hub ---
	wsHub := hub.New()
	go wsHub.Run()

	// --- iRacing Reader ---
	var reader irsdk.Reader
	if *simulate {
		reader = irsdk.NewSimulator()
	} else {
		reader = irsdk.NewWindowsReader()
	}

	// --- Background Workers ---
	ctx, cancel := context.WithCancel(context.Background())

	go telemetryPump(ctx, reader, ring, wsHub, *simulate)
	go storage.NewDownsampler(db, ring).Run(ctx)

	// --- HTTP Server ---
	engine := api.NewServer(wsHub, ring, db, *simulate, staticFiles)
	srv := &http.Server{
		Addr:         *addr,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Auto-open browser
	if *open {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost" + *addr)
	}

	fmt.Println("► Dashboard running. Press Ctrl+C to stop.")

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n► Shutting down...")
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	reader.Close()
	fmt.Println("► Goodbye!")
}

// wsMessage is the envelope for all WebSocket frames
type wsMessage struct {
	Type string `json:"type"`
	TS   int64  `json:"ts"`
	Data any    `json:"data,omitempty"`
	YAML string `json:"yaml,omitempty"`
}

func telemetryPump(ctx context.Context, r irsdk.Reader, ring *storage.RingBuffer, h *hub.Hub, simulate bool) {
	ticker := time.NewTicker(time.Second / 60) // 60 Hz
	defer ticker.Stop()

	var lastSessionUpd int32
	reconnectBackoff := time.Second

	// Initial connect attempt
	if err := r.Connect(); err != nil {
		log.Printf("pump: initial connect failed (%v) – will retry", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			frame, err := r.ReadFrame()
			if err != nil {
				if !simulate {
					log.Printf("pump: read error: %v – retrying in %v", err, reconnectBackoff)
					time.Sleep(reconnectBackoff)
					if reconnectBackoff < 30*time.Second {
						reconnectBackoff *= 2
					}
				}
				continue
			}
			reconnectBackoff = time.Second // reset on success

			ring.Push(frame)

			// Marshal and broadcast telemetry
			msg := wsMessage{Type: "telemetry", TS: time.Now().UnixMilli(), Data: frame}
			if b, err := json.Marshal(msg); err == nil {
				h.Broadcast(b)
			}

			// Broadcast session YAML when it changes
			upd := r.SessionUpdateCount()
			if upd != lastSessionUpd {
				lastSessionUpd = upd
				yaml, _ := r.SessionYAML()
				if yaml != "" {
					sessionMsg := wsMessage{Type: "session", TS: time.Now().UnixMilli(), YAML: yaml}
					if b, err := json.Marshal(sessionMsg); err == nil {
						h.Broadcast(b)
					}
				}
			}
		}
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	exec.Command(cmd, args...).Start()
}
