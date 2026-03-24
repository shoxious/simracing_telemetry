package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"simracing/hub"
	"simracing/storage"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (local use only)
	},
}

// --- WebSocket ---

type statusMsg struct {
	Type      string `json:"type"`
	Connected bool   `json:"connected"`
	Simulate  bool   `json:"simulate"`
	TS        int64  `json:"ts"`
}

func wsHandler(h *hub.Hub, ring *storage.RingBuffer, simulate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// Build initial message: status + latest telemetry frame if available
		initial := mustMarshal(statusMsg{
			Type:      "status",
			Connected: true,
			Simulate:  simulate,
			TS:        time.Now().UnixMilli(),
		})

		h.ServeClient(conn, initial)
	}
}

// --- REST Handlers ---

func statusHandler(h *hub.Hub, simulate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"connected":    true,
			"simulate":     simulate,
			"client_count": h.ClientCount(),
			"ts":           time.Now().UnixMilli(),
		})
	}
}

func lapsHandler(db *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 50
		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}
		laps, err := db.GetLaps(limit)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if laps == nil {
			laps = []storage.LapRecord{}
		}
		c.JSON(200, laps)
	}
}

func historyHandler(db *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		secs := 60
		if s := c.Query("seconds"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 600 {
				secs = v
			}
		}
		data, err := db.GetTelemetryHistory(secs)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if data == nil {
			data = []map[string]interface{}{}
		}
		c.JSON(200, data)
	}
}

func latestHandler(ring *storage.RingBuffer) gin.HandlerFunc {
	return func(c *gin.Context) {
		f := ring.Latest()
		if f == nil {
			c.JSON(204, nil)
			return
		}
		c.JSON(200, f)
	}
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
