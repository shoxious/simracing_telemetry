package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite" // pure-Go SQLite, no CGo needed for Windows cross-compilation
)

const schema = `
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA cache_size=-65536;
PRAGMA temp_store=MEMORY;

CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,
    started_at INTEGER NOT NULL,
    track      TEXT,
    car        TEXT,
    session_type TEXT
);

CREATE TABLE IF NOT EXISTS laps (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id   TEXT NOT NULL,
    lap_number   INTEGER NOT NULL,
    lap_time     REAL,
    fuel_used    REAL,
    max_speed    REAL,
    avg_throttle REAL,
    completed_at INTEGER NOT NULL,
    UNIQUE(session_id, lap_number)
);

CREATE TABLE IF NOT EXISTS telemetry_1hz (
    ts           INTEGER NOT NULL,
    session_id   TEXT NOT NULL,
    speed        REAL,
    rpm          REAL,
    gear         INTEGER,
    throttle     REAL,
    brake        REAL,
    fuel_level   REAL,
    lap_dist_pct REAL,
    lap          INTEGER
);

CREATE INDEX IF NOT EXISTS idx_tel_session ON telemetry_1hz(session_id, ts DESC);
CREATE INDEX IF NOT EXISTS idx_laps_session ON laps(session_id, lap_number);
`

// DB wraps the SQLite connection and exposes storage operations.
type DB struct {
	db        *sql.DB
	SessionID string
}

// Open opens (or creates) the SQLite database at path and applies the schema.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite is single-writer
	if _, err = db.Exec(schema); err != nil {
		return nil, err
	}

	sessionID := time.Now().Format("20060102-150405")
	_, _ = db.Exec(`INSERT OR IGNORE INTO sessions(id, started_at) VALUES (?,?)`,
		sessionID, time.Now().UnixMilli())

	return &DB{db: db, SessionID: sessionID}, nil
}

// InsertTelemetry saves a downsampled 1Hz telemetry snapshot.
func (d *DB) InsertTelemetry(ts int64, speed, rpm float32, gear int32, throttle, brake, fuel float32, lapDist float32, lap int32) {
	_, err := d.db.Exec(`
		INSERT INTO telemetry_1hz(ts,session_id,speed,rpm,gear,throttle,brake,fuel_level,lap_dist_pct,lap)
		VALUES(?,?,?,?,?,?,?,?,?,?)`,
		ts, d.SessionID, speed, rpm, gear, throttle, brake, fuel, lapDist, lap)
	if err != nil {
		log.Printf("storage: insert telemetry: %v", err)
	}
}

// InsertLap saves a completed lap record.
func (d *DB) InsertLap(lapNumber int32, lapTime, fuelUsed, maxSpeed, avgThrottle float64) {
	_, err := d.db.Exec(`
		INSERT OR REPLACE INTO laps(session_id,lap_number,lap_time,fuel_used,max_speed,avg_throttle,completed_at)
		VALUES(?,?,?,?,?,?,?)`,
		d.SessionID, lapNumber, lapTime, fuelUsed, maxSpeed, avgThrottle,
		time.Now().UnixMilli())
	if err != nil {
		log.Printf("storage: insert lap: %v", err)
	}
}

// LapRecord represents a stored lap.
type LapRecord struct {
	LapNumber   int32   `json:"lap_number"`
	LapTime     float64 `json:"lap_time"`
	FuelUsed    float64 `json:"fuel_used"`
	MaxSpeed    float64 `json:"max_speed"`
	AvgThrottle float64 `json:"avg_throttle"`
	CompletedAt int64   `json:"completed_at"`
}

// GetLaps returns the last n laps for the current session.
func (d *DB) GetLaps(limit int) ([]LapRecord, error) {
	rows, err := d.db.QueryContext(context.Background(), `
		SELECT lap_number,lap_time,fuel_used,max_speed,avg_throttle,completed_at
		FROM laps WHERE session_id=? ORDER BY lap_number DESC LIMIT ?`,
		d.SessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var laps []LapRecord
	for rows.Next() {
		var l LapRecord
		if err := rows.Scan(&l.LapNumber, &l.LapTime, &l.FuelUsed, &l.MaxSpeed, &l.AvgThrottle, &l.CompletedAt); err != nil {
			continue
		}
		laps = append(laps, l)
	}
	return laps, nil
}

// GetTelemetryHistory returns telemetry samples for the last `seconds` seconds.
func (d *DB) GetTelemetryHistory(seconds int) ([]map[string]interface{}, error) {
	since := time.Now().Add(-time.Duration(seconds) * time.Second).UnixMilli()
	rows, err := d.db.QueryContext(context.Background(), `
		SELECT ts,speed,rpm,gear,throttle,brake,fuel_level,lap_dist_pct,lap
		FROM telemetry_1hz WHERE session_id=? AND ts>=? ORDER BY ts ASC`,
		d.SessionID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var ts int64
		var speed, rpm, throttle, brake, fuel, lapDist float64
		var gear, lap int32
		if err := rows.Scan(&ts, &speed, &rpm, &gear, &throttle, &brake, &fuel, &lapDist, &lap); err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"ts": ts, "speed": speed, "rpm": rpm, "gear": gear,
			"throttle": throttle, "brake": brake, "fuel": fuel,
			"lapDist": lapDist, "lap": lap,
		})
	}
	return result, nil
}

// Close closes the database connection.
func (d *DB) Close() {
	d.db.Close()
}

// Downsampler writes 1Hz snapshots to SQLite from the ring buffer.
type Downsampler struct {
	db      *DB
	ring    *RingBuffer
	lastLap int32

	// Lap stats accumulator (reset each lap)
	lapFuelStart  float32
	lapMaxSpeed   float32
	lapThrottleAcc float64
	lapFrameCount  int
}

// NewDownsampler creates a downsampler.
func NewDownsampler(db *DB, ring *RingBuffer) *Downsampler {
	return &Downsampler{db: db, ring: ring, lapFuelStart: -1}
}

// Run starts the 1Hz sampling loop. Block until ctx is cancelled.
func (ds *Downsampler) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ds.sample()
		}
	}
}

func (ds *Downsampler) sample() {
	f := ds.ring.Latest()
	if f == nil {
		return
	}

	now := time.Now().UnixMilli()
	ds.db.InsertTelemetry(now, f.Speed, f.RPM, f.Gear, f.Throttle, f.Brake, f.FuelLevel, f.LapDistPct, f.Lap)

	// Initialise lap fuel tracking
	if ds.lapFuelStart < 0 {
		ds.lapFuelStart = f.FuelLevel
	}

	// Accumulate lap stats
	if float64(f.Speed) > float64(ds.lapMaxSpeed) {
		ds.lapMaxSpeed = f.Speed
	}
	ds.lapThrottleAcc += float64(f.Throttle)
	ds.lapFrameCount++

	// Lap change detection
	if f.Lap > ds.lastLap && ds.lastLap > 0 {
		fuelUsed := float64(ds.lapFuelStart - f.FuelLevel)
		avgThrottle := 0.0
		if ds.lapFrameCount > 0 {
			avgThrottle = ds.lapThrottleAcc / float64(ds.lapFrameCount)
		}
		ds.db.InsertLap(
			ds.lastLap,
			f.LapLastLapTime,
			fuelUsed,
			float64(ds.lapMaxSpeed)*3.6, // convert to km/h
			avgThrottle,
		)

		// Reset accumulators
		ds.lapFuelStart = f.FuelLevel
		ds.lapMaxSpeed = 0
		ds.lapThrottleAcc = 0
		ds.lapFrameCount = 0
	}
	ds.lastLap = f.Lap
}
