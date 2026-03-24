//go:build windows

package irsdk

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	memMapFileName = "Local\\IRSDKMemMapFileName"
	dataValidEvent = "Local\\IRSDKDataValidEvent"
)

type windowsReader struct {
	mu          sync.RWMutex
	hMemMapFile windows.Handle
	hDataEvent  windows.Handle
	pSharedMem  uintptr
	connected   bool
	varMap      map[string]varInfo
	header      *irHeader
	sessionUpd  int32
	sessionYAML string
}

// NewWindowsReader creates a reader that reads from iRacing's shared memory.
func NewWindowsReader() Reader {
	return &windowsReader{
		varMap: make(map[string]varInfo),
	}
}

func (r *windowsReader) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	namePtr, err := windows.UTF16PtrFromString(memMapFileName)
	if err != nil {
		return fmt.Errorf("utf16 conversion: %w", err)
	}

	h, err := windows.OpenFileMapping(windows.FILE_MAP_READ, false, namePtr)
	if err != nil {
		return fmt.Errorf("iRacing not running (OpenFileMapping): %w", err)
	}

	mem, err := windows.MapViewOfFile(h, windows.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		windows.CloseHandle(h)
		return fmt.Errorf("MapViewOfFile: %w", err)
	}

	// Optional event handle for efficient polling
	eventNamePtr, _ := windows.UTF16PtrFromString(dataValidEvent)
	eventH, _ := windows.OpenEvent(windows.SYNCHRONIZE, false, eventNamePtr)

	r.hMemMapFile = h
	r.pSharedMem = mem
	r.hDataEvent = eventH
	r.header = (*irHeader)(unsafe.Pointer(mem)) //nolint:unsafeptr

	r.buildVarMap()
	r.connected = true
	return nil
}

func (r *windowsReader) buildVarMap() {
	h := r.header
	if h.NumVars <= 0 || h.VarHeaderOffset <= 0 {
		return
	}
	varHeaderBase := r.pSharedMem + uintptr(h.VarHeaderOffset)
	for i := int32(0); i < h.NumVars; i++ {
		vh := (*irVarHeader)(unsafe.Pointer(varHeaderBase + uintptr(i)*144)) //nolint:unsafeptr
		name := nullTermStr(vh.Name[:])
		r.varMap[name] = varInfo{
			offset: vh.Offset,
			typ:    vh.Type,
			count:  vh.Count,
		}
	}
}

func (r *windowsReader) IsConnected() bool {
	if r.pSharedMem == 0 {
		return false
	}
	return r.header.Status == 1
}

func (r *windowsReader) freshestBuf() *irVarBuf {
	best := &r.header.VarBuf[0]
	for i := 1; i < 4; i++ {
		if r.header.VarBuf[i].TickCount > best.TickCount {
			best = &r.header.VarBuf[i]
		}
	}
	return best
}

func (r *windowsReader) rf32(base uintptr, name string) float32 {
	vi, ok := r.varMap[name]
	if !ok || vi.typ != VarFloat {
		return 0
	}
	return *(*float32)(unsafe.Pointer(base + uintptr(vi.offset))) //nolint:unsafeptr
}

func (r *windowsReader) rf64(base uintptr, name string) float64 {
	vi, ok := r.varMap[name]
	if !ok || vi.typ != VarDouble {
		return 0
	}
	return *(*float64)(unsafe.Pointer(base + uintptr(vi.offset))) //nolint:unsafeptr
}

func (r *windowsReader) ri32(base uintptr, name string) int32 {
	vi, ok := r.varMap[name]
	if !ok || (vi.typ != VarInt && vi.typ != VarBitField) {
		return 0
	}
	return *(*int32)(unsafe.Pointer(base + uintptr(vi.offset))) //nolint:unsafeptr
}

func (r *windowsReader) rbool(base uintptr, name string) bool {
	vi, ok := r.varMap[name]
	if !ok || vi.typ != VarBool {
		return false
	}
	return *(*bool)(unsafe.Pointer(base + uintptr(vi.offset))) //nolint:unsafeptr
}

func (r *windowsReader) ri32arr(base uintptr, name string, max int) []int32 {
	vi, ok := r.varMap[name]
	if !ok {
		return nil
	}
	count := int(vi.count)
	if count > max {
		count = max
	}
	out := make([]int32, count)
	for i := range out {
		out[i] = *(*int32)(unsafe.Pointer(base + uintptr(vi.offset) + uintptr(i*4))) //nolint:unsafeptr
	}
	return out
}

func (r *windowsReader) rf32arr(base uintptr, name string, max int) []float32 {
	vi, ok := r.varMap[name]
	if !ok {
		return nil
	}
	count := int(vi.count)
	if count > max {
		count = max
	}
	out := make([]float32, count)
	for i := range out {
		out[i] = *(*float32)(unsafe.Pointer(base + uintptr(vi.offset) + uintptr(i*4))) //nolint:unsafeptr
	}
	return out
}

func (r *windowsReader) ReadFrame() (*TelemetryFrame, error) {
	if !r.IsConnected() {
		if err := r.Connect(); err != nil {
			return nil, err
		}
	}

	// Wait for fresh data (33ms timeout = ~30Hz minimum)
	if r.hDataEvent != 0 {
		windows.WaitForSingleObject(r.hDataEvent, 33)
	} else {
		time.Sleep(16 * time.Millisecond)
	}

	buf := r.freshestBuf()
	base := r.pSharedMem + uintptr(buf.BufOffset)

	frame := &TelemetryFrame{
		Speed:             r.rf32(base, "Speed"),
		RPM:               r.rf32(base, "RPM"),
		Gear:              r.ri32(base, "Gear"),
		SteeringAngle:     r.rf32(base, "SteeringWheelAngle"),
		Throttle:          r.rf32(base, "Throttle"),
		Brake:             r.rf32(base, "Brake"),
		Clutch:            r.rf32(base, "Clutch"),
		FuelLevel:         r.rf32(base, "FuelLevel"),
		FuelLevelPct:      r.rf32(base, "FuelLevelPct"),
		FuelUsePerHour:    r.rf32(base, "FuelUsePerHour"),
		Lap:               r.ri32(base, "Lap"),
		LapDistPct:        r.rf32(base, "LapDistPct"),
		LapCurrentLapTime: r.rf64(base, "LapCurrentLapTime"),
		LapLastLapTime:    r.rf64(base, "LapLastLapTime"),
		LapBestLapTime:    r.rf64(base, "LapBestLapTime"),
		SessionTime:       r.rf64(base, "SessionTime"),
		SessionFlags:      r.ri32(base, "SessionFlags"),
		IsOnTrack:         r.rbool(base, "IsOnTrack"),
		PlayerCarIdx:      r.ri32(base, "PlayerCarIdx"),
		TrackTemp:         r.rf32(base, "TrackTemp"),
		AirTemp:           r.rf32(base, "AirTemp"),
		LFtempCL:          r.rf32(base, "LFtempCL"),
		LFtempCM:          r.rf32(base, "LFtempCM"),
		LFtempCR:          r.rf32(base, "LFtempCR"),
		RFtempCL:          r.rf32(base, "RFtempCL"),
		RFtempCM:          r.rf32(base, "RFtempCM"),
		RFtempCR:          r.rf32(base, "RFtempCR"),
		LRtempCL:          r.rf32(base, "LRtempCL"),
		LRtempCM:          r.rf32(base, "LRtempCM"),
		LRtempCR:          r.rf32(base, "LRtempCR"),
		RRtempCL:          r.rf32(base, "RRtempCL"),
		RRtempCM:          r.rf32(base, "RRtempCM"),
		RRtempCR:          r.rf32(base, "RRtempCR"),
		CarIdxPosition:    r.ri32arr(base, "CarIdxPosition", 64),
		CarIdxLapDist:     r.rf32arr(base, "CarIdxLapDist", 64),
		CarIdxEstTime:     r.rf32arr(base, "CarIdxEstTime", 64),
		CarIdxLap:         r.ri32arr(base, "CarIdxLap", 64),
	}

	// Refresh session YAML when it changes
	newUpd := r.header.SessionInfoUpdate
	if newUpd != r.sessionUpd {
		r.sessionUpd = newUpd
		offset := r.header.SessionInfoOffset
		length := r.header.SessionInfoLen
		if offset > 0 && length > 0 {
			yamlBytes := make([]byte, length)
			for i := int32(0); i < length; i++ {
				yamlBytes[i] = *(*byte)(unsafe.Pointer(r.pSharedMem + uintptr(offset) + uintptr(i))) //nolint:unsafeptr
			}
			r.sessionYAML = string(yamlBytes)
		}
	}

	return frame, nil
}

func (r *windowsReader) SessionYAML() (string, error) {
	return r.sessionYAML, nil
}

func (r *windowsReader) SessionUpdateCount() int32 {
	return r.sessionUpd
}

func (r *windowsReader) Close() {
	if r.pSharedMem != 0 {
		windows.UnmapViewOfFile(r.pSharedMem)
		r.pSharedMem = 0
	}
	if r.hMemMapFile != 0 {
		windows.CloseHandle(r.hMemMapFile)
		r.hMemMapFile = 0
	}
	if r.hDataEvent != 0 {
		windows.CloseHandle(r.hDataEvent)
		r.hDataEvent = 0
	}
	r.connected = false
}
