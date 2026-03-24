package hub

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
	maxMessageSize = 8192
	sendBufSize    = 256 // per-client send channel buffer
)

// Client is a single WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub manages all active WebSocket clients and broadcasts messages to all of them.
type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]struct{}
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// New creates and initialises a Hub (call Run() in a goroutine).
func New() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte, 512),
		register:   make(chan *Client, 32),
		unregister: make(chan *Client, 32),
	}
}

// Run processes hub events; must be started in its own goroutine.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			// Two-phase: send, then remove slow clients (avoids RLock→Lock upgrade deadlock)
			var slow []*Client
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					slow = append(slow, c)
				}
			}
			h.mu.RUnlock()
			if len(slow) > 0 {
				h.mu.Lock()
				for _, c := range slow {
					if _, ok := h.clients[c]; ok {
						delete(h.clients, c)
						close(c.send)
					}
				}
				h.mu.Unlock()
			}
		}
	}
}

// Broadcast enqueues a message to be sent to all clients.
// Non-blocking: drops the message if the broadcast channel is full.
func (h *Hub) Broadcast(msg []byte) {
	select {
	case h.broadcast <- msg:
	default:
		log.Println("hub: broadcast channel full, dropping frame")
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ServeClient registers the client and starts its read/write pumps.
// conn must already be upgraded to WebSocket.
func (h *Hub) ServeClient(conn *websocket.Conn, initialMsg []byte) {
	c := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, sendBufSize),
	}
	h.register <- c

	// Send the initial status/session message immediately
	if initialMsg != nil {
		c.send <- initialMsg
	}

	go c.writePump()
	c.readPump() // blocks until connection closes
}

// readPump discards incoming messages and handles pong responses.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// We don't process client→server messages in this version,
		// just drain them so the connection stays healthy.
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

// writePump serialises writes to the WebSocket connection and sends pings.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
