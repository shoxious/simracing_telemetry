package storage

import (
	"sync"

	"simracing/irsdk"
)

// RingBuffer holds the last N TelemetryFrames in memory with O(1) writes.
// Safe for concurrent use.
type RingBuffer struct {
	mu   sync.RWMutex
	buf  []*irsdk.TelemetryFrame
	head int // next write position
	size int // number of valid entries
	cap  int
}

// NewRingBuffer creates a ring buffer with the given capacity.
// Default: 7200 frames = 120 seconds at 60Hz (~3-5 MB)
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		buf: make([]*irsdk.TelemetryFrame, capacity),
		cap: capacity,
	}
}

// Push adds a frame, overwriting the oldest entry when full.
func (r *RingBuffer) Push(f *irsdk.TelemetryFrame) {
	r.mu.Lock()
	r.buf[r.head] = f
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
	r.mu.Unlock()
}

// Latest returns the most recently pushed frame, or nil if empty.
func (r *RingBuffer) Latest() *irsdk.TelemetryFrame {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.size == 0 {
		return nil
	}
	idx := (r.head - 1 + r.cap) % r.cap
	return r.buf[idx]
}

// Snapshot returns up to n most-recent frames as a slice (oldest first).
// Returns a copy of the pointers; callers must not mutate the frames.
func (r *RingBuffer) Snapshot(n int) []*irsdk.TelemetryFrame {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.size == 0 {
		return nil
	}
	if n > r.size {
		n = r.size
	}

	out := make([]*irsdk.TelemetryFrame, n)
	for i := 0; i < n; i++ {
		// oldest of the n requested frames first
		idx := (r.head - n + i + r.cap*2) % r.cap
		out[i] = r.buf[idx]
	}
	return out
}

// Len returns the number of valid entries.
func (r *RingBuffer) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}
