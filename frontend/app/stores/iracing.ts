import { defineStore } from 'pinia'

export interface TelemetryFrame {
  Speed: number
  RPM: number
  Gear: number
  SteeringAngle: number
  Throttle: number
  Brake: number
  Clutch: number
  FuelLevel: number
  FuelLevelPct: number
  FuelUsePerHour: number
  Lap: number
  LapDistPct: number
  LapCurrentLapTime: number
  LapLastLapTime: number
  LapBestLapTime: number
  SessionTime: number
  SessionFlags: number
  IsOnTrack: boolean
  PlayerCarIdx: number
  TrackTemp: number
  AirTemp: number
  LFtempCL: number; LFtempCM: number; LFtempCR: number
  RFtempCL: number; RFtempCM: number; RFtempCR: number
  LRtempCL: number; LRtempCM: number; LRtempCR: number
  RRtempCL: number; RRtempCM: number; RRtempCR: number
  CarIdxPosition: number[]
  CarIdxLapDist: number[]
  CarIdxEstTime: number[]
  CarIdxLap: number[]
}

export interface LapRecord {
  lap_number: number
  lap_time: number
  fuel_used: number
  max_speed: number
  avg_throttle: number
  completed_at: number
}

export const useIRacingStore = defineStore('iracing', {
  state: () => ({
    connected:   false,
    simulate:    false,
    telemetry:   null as TelemetryFrame | null,
    sessionYaml: null as string | null,
    laps:        [] as LapRecord[],
    lastUpdated: 0,

    // FPS tracking
    _frameCount: 0,
    _fps:        0,
  }),

  getters: {
    fps: (state) => state._fps,

    speedKmh: (state): number => {
      return (state.telemetry?.Speed ?? 0) * 3.6
    },

    speedMph: (state): number => {
      return (state.telemetry?.Speed ?? 0) * 2.237
    },

    rpmPct: (state): number => {
      const rpm = state.telemetry?.RPM ?? 0
      return Math.min(rpm / 8600, 1)
    },

    gearLabel: (state): string => {
      const g = state.telemetry?.Gear ?? 0
      if (g === 0) return 'N'
      if (g === -1) return 'R'
      return String(g)
    },

    fuelLapsRemaining: (state): number => {
      const level = state.telemetry?.FuelLevel ?? 0
      const usePerHour = state.telemetry?.FuelUsePerHour ?? 0
      if (usePerHour <= 0) return 99
      // Convert use-per-hour to use-per-lap assuming ~90s laps
      const usePerLap = usePerHour * (90 / 3600)
      if (usePerLap <= 0) return 99
      return Math.floor(level / usePerLap)
    },

    playerPosition: (state): number => {
      const idx = state.telemetry?.PlayerCarIdx ?? 0
      return state.telemetry?.CarIdxPosition?.[idx] ?? 0
    },

    isRedFlag: (state): boolean => {
      return Boolean(state.telemetry?.SessionFlags && (state.telemetry.SessionFlags & 0x0010))
    },

    isYellowFlag: (state): boolean => {
      return Boolean(state.telemetry?.SessionFlags && (state.telemetry.SessionFlags & 0x0002))
    },

    isGreenFlag: (state): boolean => {
      return Boolean(state.telemetry?.SessionFlags && (state.telemetry.SessionFlags & 0x00040000))
    },

    deltaToBest: (state): number => {
      const cur = state.telemetry?.LapCurrentLapTime ?? 0
      const best = state.telemetry?.LapBestLapTime ?? 0
      if (best <= 0) return 0
      // Very rough estimate: scale by lap progress
      const pct = state.telemetry?.LapDistPct ?? 0
      if (pct <= 0.01) return 0
      const projected = cur / pct
      return projected - best
    },
  },

  actions: {
    setTelemetry(data: TelemetryFrame) {
      this.telemetry = data
      this.lastUpdated = Date.now()
      this._frameCount++
    },

    setStatus(connected: boolean, simulate: boolean) {
      this.connected = connected
      this.simulate = simulate
    },

    setSession(yaml: string) {
      this.sessionYaml = yaml
    },

    setLaps(laps: LapRecord[]) {
      this.laps = laps
    },

    tickFPS() {
      this._fps = this._frameCount
      this._frameCount = 0
    },
  },
})
