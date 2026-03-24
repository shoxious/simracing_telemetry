package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/mdp/qrterminal/v3"

	"simracing/api"
	"simracing/hub"
	"simracing/irsdk"
	"simracing/storage"
)

//go:embed all:static
var staticFiles embed.FS

// shutdownCh is closed / sent-to from the console close handler (Windows)
// or a SIGINT/SIGTERM signal to trigger graceful shutdown.
var shutdownCh = make(chan struct{}, 2)

func main() {
	var (
		addr     = flag.String("addr", ":8080", "HTTP listen address")
		simulate = flag.Bool("simulate", false, "Simulation mode (no iRacing required)")
		dbPath   = flag.String("db", "simracing.db", "SQLite database path")
		open     = flag.Bool("open", true, "Automatically open browser on start")
		mdnsName = flag.String("name", "simracing", "mDNS hostname  → <name>.local")
	)
	flag.Parse()

	// Set console title and register close-window → shutdown handler (Windows only).
	initConsole()

	port := portFromAddr(*addr)

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Printf( "║   SimRacing Live Dashboard  %-11s║\n", version)
	fmt.Println("╚════════════════════════════════════════╝")

	// Check GitHub for a newer release before starting the server.
	checkForUpdate()
	fmt.Println()

	if *simulate {
		fmt.Println("► Mode: SIMULATION (synthetic telemetry)")
	} else {
		fmt.Println("► Mode: LIVE (iRacing shared memory)")
	}
	fmt.Printf("► Local:    http://localhost%s\n", *addr)
	fmt.Printf("► Network:  http://%s.local:%d\n", *mdnsName, port)
	fmt.Printf("► Database: %s\n", *dbPath)
	fmt.Println()

	// --- Storage ---
	db, err := storage.Open(*dbPath)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	defer db.Close()

	// 120s ring buffer at 60 Hz
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

	// --- mDNS / Bonjour registration ---
	mdnsServer, mdnsErr := startMDNS(*mdnsName, port)
	if mdnsErr != nil {
		fmt.Printf("► mDNS: unavailable (%v) – use IP address instead\n", mdnsErr)
	} else {
		fmt.Printf("► mDNS: registered as http://%s.local:%d\n", *mdnsName, port)
	}

	// --- QR Code for the first non-loopback IPv4 address ---
	if networkIP := firstNetworkIP(); networkIP != "" {
		networkURL := fmt.Sprintf("http://%s:%d", networkIP, port)
		fmt.Printf("\n► Scan to open on your phone / tablet:\n   %s\n\n", networkURL)
		qrterminal.GenerateWithConfig(networkURL, qrterminal.Config{
			Level:     qrterminal.M,
			Writer:    os.Stdout,
			BlackChar: qrterminal.BLACK,
			WhiteChar: qrterminal.WHITE,
			QuietZone: 1,
		})
		fmt.Println()
	}

	// --- Auto-open browser ---
	if *open {
		time.Sleep(500 * time.Millisecond)
		openBrowser(fmt.Sprintf("http://localhost%s", *addr))
	}

	printNetworkAddresses(port)
	fmt.Println("► Dashboard running. Press Ctrl+C to stop (or close this window).")

	// --- Graceful Shutdown ---
	// Wait for SIGINT / SIGTERM  OR  console-window close (Windows).
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-osSignal:
	case <-shutdownCh:
	}

	fmt.Println("\n► Shutting down...")
	if mdnsServer != nil {
		mdnsServer.Shutdown()
	}
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	reader.Close()
	fmt.Println("► Goodbye!")
}

// startMDNS registers the service via mDNS so it is reachable as <name>.local.
func startMDNS(name string, port int) (*zeroconf.Server, error) {
	return zeroconf.Register(
		"SimRacing Dashboard", // instance name
		"_http._tcp",          // service type
		"local.",              // domain
		port,
		[]string{"path=/"}, // TXT records
		nil,                // bind to all interfaces
	)
}

// portFromAddr extracts the integer port from ":8080" or "0.0.0.0:8080".
func portFromAddr(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) == 0 {
		return 8080
	}
	p, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 8080
	}
	return p
}

// firstNetworkIP returns the first non-loopback IPv4 address or "".
func firstNetworkIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""
}

// printNetworkAddresses prints all non-loopback IPv4 addresses.
func printNetworkAddresses(port int) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			fmt.Printf("► Network: http://%s:%d\n", ip.String(), port)
		}
	}
}

// wsMessage is the envelope for all WebSocket frames.
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
			reconnectBackoff = time.Second

			ring.Push(frame)

			msg := wsMessage{Type: "telemetry", TS: time.Now().UnixMilli(), Data: frame}
			if b, err := json.Marshal(msg); err == nil {
				h.Broadcast(b)
			}

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
