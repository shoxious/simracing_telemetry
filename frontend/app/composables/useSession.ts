/**
 * useSession – parses the iRacing session YAML and exposes structured data.
 * The YAML is already pushed by the backend on every session update;
 * here we extract the fields the UI needs.
 */

export interface DriverEntry {
  carIdx:        number
  userName:      string
  carNumber:     string
  carScreenName: string
  teamName:      string
  iRating:       number
  isSpectator:   boolean
}

export interface ParsedSession {
  trackDisplayName:      string
  trackDisplayShortName: string
  trackLength:           string
  trackCity:             string
  trackCountry:          string
  drivers:               DriverEntry[]
}

const EMPTY: ParsedSession = {
  trackDisplayName: '', trackDisplayShortName: '',
  trackLength: '', trackCity: '', trackCountry: '',
  drivers: [],
}

/** Pull a single scalar value from the YAML string. */
function scalar(yaml: string, key: string): string {
  return yaml.match(new RegExp(key + ':\\s*(.+)'))?.[1]?.trim() ?? ''
}

function parseDrivers(yaml: string): DriverEntry[] {
  // Each driver block starts with "- CarIdx:" (indented)
  const blocks = yaml.split(/\n\s+- CarIdx:/)
  const drivers: DriverEntry[] = []

  for (let i = 1; i < blocks.length; i++) {
    const b = blocks[i]
    const carIdx = parseInt(b.match(/^\s*(\d+)/)?.[1] ?? '-1')
    if (carIdx < 0) continue

    drivers.push({
      carIdx,
      userName:      scalar(b, 'UserName'),
      carNumber:     scalar(b, 'CarNumber').replace(/"/g, ''),
      carScreenName: scalar(b, 'CarScreenName'),
      teamName:      scalar(b, 'TeamName'),
      iRating:       parseInt(scalar(b, 'IRating') || '0'),
      isSpectator:   scalar(b, 'IsSpectator') === '1',
    })
  }
  return drivers
}

function parseYaml(yaml: string): ParsedSession {
  return {
    trackDisplayName:      scalar(yaml, 'TrackDisplayName'),
    trackDisplayShortName: scalar(yaml, 'TrackDisplayShortName'),
    trackLength:           scalar(yaml, 'TrackLength'),
    trackCity:             scalar(yaml, 'TrackCity'),
    trackCountry:          scalar(yaml, 'TrackCountry'),
    drivers:               parseDrivers(yaml),
  }
}

export function useSession() {
  const store = useIRacingStore()

  const session = computed<ParsedSession>(() =>
    store.sessionYaml ? parseYaml(store.sessionYaml) : EMPTY
  )

  const trackName = computed(() =>
    session.value.trackDisplayShortName ||
    session.value.trackDisplayName ||
    ''
  )

  const trackLength = computed(() => session.value.trackLength)

  /** Look up a driver by car index. */
  function driverByIdx(idx: number): DriverEntry | undefined {
    return session.value.drivers.find(d => d.carIdx === idx)
  }

  /** The player's own driver entry. */
  const playerDriver = computed(() => {
    const idx = store.telemetry?.PlayerCarIdx ?? -1
    return idx >= 0 ? driverByIdx(idx) : undefined
  })

  return { session, trackName, trackLength, driverByIdx, playerDriver }
}
