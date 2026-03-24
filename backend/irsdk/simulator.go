package irsdk

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// lapZone defines the target driving conditions for a section of the track
type lapZone struct {
	startPct float64
	endPct   float64
	speed    float64 // target speed at zone end (m/s)
	throttle float64 // 0.0-1.0
	brake    float64 // 0.0-1.0
	gear     int32
}

// fictitious 4.2km GT3 circuit with 7 turns (~88s lap time)
var trackZones = []lapZone{
	// T/F straight - full beans
	{0.00, 0.12, 74.0, 1.00, 0.00, 5},
	// T1 heavy brake (Lesmo style)
	{0.12, 0.16, 19.0, 0.00, 0.95, 2},
	// T1 exit
	{0.16, 0.24, 44.0, 0.75, 0.00, 3},
	// Mid straight
	{0.24, 0.37, 66.0, 1.00, 0.00, 5},
	// T2 medium brake
	{0.37, 0.42, 26.0, 0.00, 0.70, 2},
	// T2-T3 flowing
	{0.42, 0.54, 52.0, 0.80, 0.00, 4},
	// Back straight
	{0.54, 0.68, 82.0, 1.00, 0.00, 6},
	// T4 - hardest brake of the lap
	{0.68, 0.74, 14.0, 0.00, 1.00, 1},
	// Hairpin exit
	{0.74, 0.82, 38.0, 0.65, 0.00, 3},
	// Chicane run
	{0.82, 0.90, 55.0, 0.90, 0.00, 4},
	// Return to grid
	{0.90, 1.00, 74.0, 1.00, 0.00, 5},
}

// gear ratios × final drive / wheel factor → RPM per m/s
var gearRPMFactor = [7]float64{0, 210, 160, 122, 97, 82, 69}

// simCar represents one car on track (player + AI opponents)
type simCar struct {
	offset    float64 // initial lap dist offset (0-1)
	lapOffset float64 // random lap time variance (seconds)
	carIdx    int
}

// Simulator generates realistic iRacing telemetry without a connected game.
type Simulator struct {
	mu sync.Mutex

	startTime  time.Time
	lapStart   time.Time
	lap        int32
	sessionUpd int32

	// per-car simulation (6 cars total)
	cars []simCar

	// slowly evolving state
	tireWarmup  float64 // 0→1 over first 3 laps
	fuel        float32
	bestLapTime float64
	lastLapTime float64

	// noise source
	rng *rand.Rand
}

const (
	lapDurationSec = 88.0 // base lap time
	fuelPerLap     = 2.8  // liters
	maxFuel        = 55.0 // liters
)

// NewSimulator creates a ready-to-use telemetry simulator.
func NewSimulator() Reader {
	s := &Simulator{
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		fuel:        maxFuel,
		bestLapTime: -1,
		lastLapTime: -1,
	}
	// Spread 6 cars around the track
	for i := 0; i < 6; i++ {
		s.cars = append(s.cars, simCar{
			offset:    float64(i) / 6.0,
			lapOffset: s.rng.Float64()*4 - 2, // ±2s lap time variance
			carIdx:    i,
		})
	}
	return s
}

func (s *Simulator) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.startTime = time.Now()
	s.lapStart = time.Now()
	return nil
}

func (s *Simulator) IsConnected() bool { return !s.startTime.IsZero() }

func (s *Simulator) SessionYAML() (string, error) {
	yaml := `---
WeekendInfo:
 TrackName: Monza Circuit GP
 TrackConfigName: Grand Prix
 TrackLength: 4.20 km
 TrackDisplayName: Autodromo Nazionale Monza
 EventType: Race
 Category: Road

DriverInfo:
 DriverCarIdx: 0
 Drivers:
 - CarIdx: 0
   UserName: SimRacer Pro
   CarScreenName: Mercedes-AMG GT3 Evo
   CarNumber: "99"
   IRating: 4523
   LicString: A 4.23
 - CarIdx: 1
   UserName: Max Hoffman
   CarScreenName: BMW M4 GT3
   CarNumber: "7"
   IRating: 4102
 - CarIdx: 2
   UserName: Sarah Chen
   CarScreenName: Audi R8 LMS GT3
   CarNumber: "44"
   IRating: 3850
 - CarIdx: 3
   UserName: Luca Ferrari
   CarScreenName: Ferrari 296 GT3
   CarNumber: "16"
   IRating: 3720
 - CarIdx: 4
   UserName: Alex Müller
   CarScreenName: Lamborghini Huracán GT3
   CarNumber: "63"
   IRating: 3680
 - CarIdx: 5
   UserName: Tom Bradley
   CarScreenName: Porsche 992 GT3 R
   CarNumber: "92"
   IRating: 3510

SessionInfo:
 Sessions:
 - SessionNum: 0
   SessionName: RACE
   SessionType: Race
   SessionLaps: 30
   SessionTime: 2700.0000 sec
`
	return yaml, nil
}

func (s *Simulator) SessionUpdateCount() int32 { return s.sessionUpd }

func (s *Simulator) Close() {}

func (s *Simulator) ReadFrame() (*TelemetryFrame, error) {
	if s.startTime.IsZero() {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(s.startTime).Seconds()
	lapElapsed := now.Sub(s.lapStart).Seconds()

	// Compute lap distance for player car (carIdx 0)
	lapDist := math.Mod(elapsed/lapDurationSec, 1.0)

	// Lap completion detection
	currentLap := int32(elapsed / lapDurationSec)
	if currentLap > s.lap {
		lt := lapDurationSec + s.rng.Float64()*2 - 1 // slight variation
		s.lastLapTime = lt
		if s.bestLapTime < 0 || lt < s.bestLapTime {
			s.bestLapTime = lt
		}
		s.lap = currentLap
		s.lapStart = now
		lapElapsed = 0
		s.fuel -= fuelPerLap
		if s.fuel < 0 {
			s.fuel = 0
		}
	}

	// Tire warmup ramp
	s.tireWarmup = math.Min(elapsed/(lapDurationSec*3), 1.0)

	// Get interpolated track state for player
	speed, throttle, brake, gear := s.trackState(lapDist)

	// Add subtle noise to inputs
	noise := func(base, scale float64) float32 {
		return float32(math.Max(0, math.Min(1, base+s.rng.Float64()*scale-scale/2)))
	}
	noiseF := func(base, scale float64) float32 {
		return float32(base + s.rng.Float64()*scale - scale/2)
	}

	actualThrottle := noise(throttle, 0.04)
	actualBrake := noise(brake, 0.03)
	actualSpeed := float32(speed + s.rng.Float64()*0.5 - 0.25)
	actualSteering := noiseF(steeringAtDist(lapDist), 0.05)

	rpm := float32(float64(actualSpeed) * gearRPMFactor[gear])
	rpm = float32(math.Min(float64(rpm), 8600) + s.rng.Float64()*50)

	// Tire temps: cold at start, reach operating temp over 3 laps
	baseTireCold := 55.0
	baseTireHot := 94.0
	tireBase := baseTireCold + s.tireWarmup*(baseTireHot-baseTireCold)
	tireVar := 8.0 // spread across tire width

	tf := func(base, side float64) float32 {
		return float32(base + side + s.rng.Float64()*2)
	}

	// Build per-car arrays (6 cars)
	carIdxPos := make([]int32, 64)
	carIdxLapDist := make([]float32, 64)
	carIdxEstTime := make([]float32, 64)
	carIdxLap := make([]int32, 64)

	// Sort cars by race position
	positions := s.carPositions(elapsed)
	for i, cp := range positions {
		carIdxPos[cp.carIdx] = int32(i + 1)
		carIdxLapDist[cp.carIdx] = float32(cp.lapDist)
		carIdxEstTime[cp.carIdx] = float32(cp.estTime)
		carIdxLap[cp.carIdx] = int32(cp.lap)
	}

	playerPos := carIdxPos[0]

	_ = playerPos // used implicitly via carIdxPos

	return &TelemetryFrame{
		Speed:             actualSpeed,
		RPM:               rpm,
		Gear:              gear,
		SteeringAngle:     actualSteering,
		Throttle:          actualThrottle,
		Brake:             actualBrake,
		Clutch:            boolToFloat32(brake > 0.1),
		FuelLevel:         s.fuel,
		FuelLevelPct:      s.fuel / maxFuel,
		FuelUsePerHour:    float32(fuelPerLap * 3600 / lapDurationSec),
		Lap:               s.lap + 1,
		LapDistPct:        float32(lapDist),
		LapCurrentLapTime: lapElapsed,
		LapLastLapTime:    s.lastLapTime,
		LapBestLapTime:    s.bestLapTime,
		SessionTime:       elapsed,
		SessionFlags:      0x00040000, // green flag
		IsOnTrack:         true,
		PlayerCarIdx:      0,
		TrackTemp:         float32(32 + s.rng.Float64()*2),
		AirTemp:           float32(24 + s.rng.Float64()),
		// Tire temps - asymmetric due to track layout
		LFtempCL: tf(tireBase-tireVar*0.4, -2),
		LFtempCM: tf(tireBase, 0),
		LFtempCR: tf(tireBase+tireVar*0.3, 2),
		RFtempCL: tf(tireBase+tireVar*0.2, -1),
		RFtempCM: tf(tireBase+tireVar*0.5, 0),
		RFtempCR: tf(tireBase-tireVar*0.1, 1),
		LRtempCL: tf(tireBase-tireVar*0.3, -2),
		LRtempCM: tf(tireBase+tireVar*0.1, 0),
		LRtempCR: tf(tireBase+tireVar*0.4, 2),
		RRtempCL: tf(tireBase+tireVar*0.3, -1),
		RRtempCM: tf(tireBase+tireVar*0.6, 0),
		RRtempCR: tf(tireBase-tireVar*0.2, 1),
		// Per-car
		CarIdxPosition: carIdxPos,
		CarIdxLapDist:  carIdxLapDist,
		CarIdxEstTime:  carIdxEstTime,
		CarIdxLap:      carIdxLap,
	}, nil
}

// trackState interpolates speed/throttle/brake for a given lap position
func (s *Simulator) trackState(lapDist float64) (speed, throttle, brake float64, gear int32) {
	// Find surrounding zones
	var z0, z1 lapZone
	for i, z := range trackZones {
		if lapDist >= z.startPct && lapDist <= z.endPct {
			z0 = z
			if i+1 < len(trackZones) {
				z1 = trackZones[i+1]
			} else {
				z1 = trackZones[0]
			}
			break
		}
	}

	// Interpolate within zone
	if z0.endPct == z0.startPct {
		return z0.speed, z0.throttle, z0.brake, z0.gear
	}
	t := (lapDist - z0.startPct) / (z0.endPct - z0.startPct)
	t = smoothstep(t)

	speed = z0.speed + t*(z1.speed-z0.speed)
	throttle = z0.throttle + t*(z1.throttle-z0.throttle)
	brake = z0.brake + t*(z1.brake-z0.brake)
	gear = z0.gear

	// Gear selection based on speed
	switch {
	case speed < 20:
		gear = 1
	case speed < 33:
		gear = 2
	case speed < 48:
		gear = 3
	case speed < 62:
		gear = 4
	case speed < 74:
		gear = 5
	default:
		gear = 6
	}

	return
}

// steeringAtDist returns approximate steering angle (radians) for track position
func steeringAtDist(lapDist float64) float64 {
	// Sinusoidal steering approximation
	return math.Sin(lapDist*2*math.Pi*4) * 0.3
}

type carState struct {
	carIdx  int
	lapDist float64
	lap     int32
	estTime float32
}

// carPositions computes all car positions sorted by race order
func (s *Simulator) carPositions(elapsed float64) []carState {
	cars := make([]carState, len(s.cars))
	for i, c := range s.cars {
		totalDist := elapsed/lapDurationSec + c.offset
		lapNum := int32(totalDist)
		lapDist := totalDist - float64(lapNum)
		// Estimate gap to player in seconds
		playerTotalDist := elapsed / lapDurationSec
		gapLaps := playerTotalDist - totalDist
		estTime := float32(gapLaps * lapDurationSec)

		cars[i] = carState{
			carIdx:  c.carIdx,
			lapDist: lapDist,
			lap:     lapNum,
			estTime: estTime,
		}
	}

	// Sort by total distance covered (descending = race leader first)
	for i := 0; i < len(cars); i++ {
		for j := i + 1; j < len(cars); j++ {
			di := float64(cars[i].lap) + cars[i].lapDist
			dj := float64(cars[j].lap) + cars[j].lapDist
			if dj > di {
				cars[i], cars[j] = cars[j], cars[i]
			}
		}
	}
	return cars
}

func smoothstep(t float64) float64 {
	return t * t * (3 - 2*t)
}

func boolToFloat32(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// SimulatorSessionYAML returns an example session string for testing
func SimulatorSessionYAML() string {
	return fmt.Sprintf("SimRacing Dashboard v1.0 - Simulation Mode\nLap Duration: %.0fs", lapDurationSec)
}
