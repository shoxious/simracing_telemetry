import { storeToRefs } from 'pinia'
import { useIRacingStore } from '~/stores/iracing'

/** Convenience composable that exposes reactive iRacing data */
export function useIRacing() {
  const store = useIRacingStore()
  const { telemetry, connected, simulate, sessionYaml, laps, lastUpdated } = storeToRefs(store)

  const speedKmh       = computed(() => store.speedKmh)
  const speedMph       = computed(() => store.speedMph)
  const rpmPct         = computed(() => store.rpmPct)
  const gearLabel      = computed(() => store.gearLabel)
  const fuelLaps       = computed(() => store.fuelLapsRemaining)
  const position       = computed(() => store.playerPosition)
  const deltaToBest    = computed(() => store.deltaToBest)
  const isGreen        = computed(() => store.isGreenFlag)
  const isYellow       = computed(() => store.isYellowFlag)
  const isRed          = computed(() => store.isRedFlag)

  /** Format seconds as M:SS.mmm */
  function fmtTime(secs: number | null | undefined): string {
    if (!secs || secs <= 0) return '--:--.---'
    const m = Math.floor(secs / 60)
    const s = secs % 60
    return `${m}:${s.toFixed(3).padStart(6, '0')}`
  }

  /** Format delta as ±0.000 */
  function fmtDelta(delta: number): string {
    const sign = delta >= 0 ? '+' : ''
    return `${sign}${delta.toFixed(3)}`
  }

  return {
    telemetry, connected, simulate, sessionYaml, laps, lastUpdated,
    speedKmh, speedMph, rpmPct, gearLabel, fuelLaps, position,
    deltaToBest, isGreen, isYellow, isRed,
    fmtTime, fmtDelta,
    store,
  }
}
