package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"simracing/hub"
	"simracing/storage"
)

// NewServer creates the Gin engine with all routes configured.
// staticFiles is the embedded FS from main (contains the "static" subdirectory).
func NewServer(h *hub.Hub, ring *storage.RingBuffer, db *storage.DB, simulate bool, staticFiles embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS for Nuxt dev server (localhost:3000)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowWebSockets:  true,
	}))

	// WebSocket endpoint
	r.GET("/ws", wsHandler(h, ring, simulate))

	// REST API
	api := r.Group("/api")
	{
		api.GET("/status", statusHandler(h, simulate))
		api.GET("/laps", lapsHandler(db))
		api.GET("/history", historyHandler(db))
		api.GET("/telemetry/latest", latestHandler(ring))
	}

	// Serve the embedded Nuxt SPA (catch-all, must be last)
	staticFS, err := fs.Sub(StaticFiles, "static")
	if err != nil {
		panic("embedded static files not found – run 'npm run generate' first: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(staticFS))

	r.NoRoute(func(c *gin.Context) {
		// Try to serve the file; fall back to index.html for SPA routes
		path := c.Request.URL.Path
		f, err := staticFS.Open(path[1:]) // strip leading /
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		// Serve index.html for all unmatched routes (client-side routing)
		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	return r
}
