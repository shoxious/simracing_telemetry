package irsdk

// iRacing SDK variable types (mirrors irsdk_VarType enum)
const (
	VarChar     = 0
	VarBool     = 1
	VarInt      = 2
	VarBitField = 3
	VarFloat    = 4
	VarDouble   = 5
)

// irHeader mirrors the irsdk_header C struct (112 bytes, must match byte-for-byte)
type irHeader struct {
	Ver               int32
	Status            int32
	TickRate          int32
	SessionInfoUpdate int32
	SessionInfoLen    int32
	SessionInfoOffset int32
	NumVars           int32
	VarHeaderOffset   int32
	NumBuf            int32
	BufLen            int32
	Pad               [2]int32
	VarBuf            [4]irVarBuf
}

// irVarBuf mirrors irsdk_bufHeader (16 bytes each)
type irVarBuf struct {
	TickCount int32
	BufOffset int32
	Pad       [2]int32
}

// irVarHeader mirrors irsdk_varHeader (144 bytes each)
type irVarHeader struct {
	Type        int32
	Offset      int32
	Count       int32
	CountAsTime uint8
	Pad         [3]byte
	Name        [32]byte
	Desc        [64]byte
	Unit        [32]byte
}

// varInfo holds the offset and type info for a telemetry variable
type varInfo struct {
	offset int32
	typ    int32
	count  int32
}

// TelemetryFrame contains all live data extracted from iRacing
type TelemetryFrame struct {
	// Motion
	Speed         float32 `json:"Speed"`         // m/s
	RPM           float32 `json:"RPM"`
	Gear          int32   `json:"Gear"`           // -1=R, 0=N, 1-6=gears
	SteeringAngle float32 `json:"SteeringAngle"`  // radians

	// Driver inputs
	Throttle float32 `json:"Throttle"` // 0.0-1.0
	Brake    float32 `json:"Brake"`    // 0.0-1.0
	Clutch   float32 `json:"Clutch"`   // 0.0-1.0

	// Fuel
	FuelLevel      float32 `json:"FuelLevel"`      // liters
	FuelLevelPct   float32 `json:"FuelLevelPct"`   // 0.0-1.0
	FuelUsePerHour float32 `json:"FuelUsePerHour"` // liters/hour

	// Lap info
	Lap               int32   `json:"Lap"`
	LapDistPct        float32 `json:"LapDistPct"`        // 0.0-1.0
	LapCurrentLapTime float64 `json:"LapCurrentLapTime"` // seconds
	LapLastLapTime    float64 `json:"LapLastLapTime"`
	LapBestLapTime    float64 `json:"LapBestLapTime"`

	// Session
	SessionTime  float64 `json:"SessionTime"` // seconds elapsed
	SessionFlags int32   `json:"SessionFlags"`
	IsOnTrack    bool    `json:"IsOnTrack"`
	PlayerCarIdx int32   `json:"PlayerCarIdx"`

	// Environment
	TrackTemp float32 `json:"TrackTemp"` // Celsius
	AirTemp   float32 `json:"AirTemp"`

	// Tire temps (CL=left, CM=middle, CR=right strip of tire)
	LFtempCL, LFtempCM, LFtempCR float32
	RFtempCL, RFtempCM, RFtempCR float32
	LRtempCL, LRtempCM, LRtempCR float32
	RRtempCL, RRtempCM, RRtempCR float32

	// Per-car arrays (64 entries each, indexed by car idx)
	CarIdxPosition []int32   `json:"CarIdxPosition"`
	CarIdxLapDist  []float32 `json:"CarIdxLapDist"`
	CarIdxEstTime  []float32 `json:"CarIdxEstTime"`
	CarIdxLap      []int32   `json:"CarIdxLap"`
}

// Reader is the interface implemented by both the real iRacing reader and the simulator
type Reader interface {
	Connect() error
	IsConnected() bool
	ReadFrame() (*TelemetryFrame, error)
	SessionYAML() (string, error)
	SessionUpdateCount() int32
	Close()
}

// nullTermStr converts a null-terminated byte slice to a Go string
func nullTermStr(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
