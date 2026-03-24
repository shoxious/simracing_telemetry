<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">

    <!-- Header card -->
    <div class="r-card px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="r-label">Relative</span>
        <span class="font-mono text-xs text-r-muted">{{ visibleCars }} cars</span>
      </div>
      <!-- Window selector -->
      <div class="flex items-center gap-1 bg-r-surface rounded-lg p-1">
        <button
          v-for="w in windows"
          :key="w.value"
          class="px-2.5 py-1 rounded-md font-mono text-xs transition-colors"
          :class="window === w.value
            ? 'bg-r-card text-r-text border border-r-border'
            : 'text-r-muted hover:text-r-text'"
          @click="window = w.value"
        >
          {{ w.label }}
        </button>
      </div>
    </div>

    <!-- Relative tower -->
    <div class="r-card overflow-hidden">
      <!-- No data state -->
      <div
        v-if="!rows.length"
        class="flex items-center justify-center h-40 text-r-muted text-xs font-mono"
      >
        Waiting for session data...
      </div>

      <div v-else>
        <div
          v-for="(row, i) in rows"
          :key="row.carIdx"
          class="relative flex items-center px-3 py-2.5 border-b border-r-border/50 last:border-b-0 transition-colors"
          :class="row.isPlayer
            ? 'bg-r-blue/10 border-l-2 !border-l-r-blue'
            : 'hover:bg-r-surface/40'"
        >
          <!-- Position badge -->
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-mono font-bold flex-shrink-0 mr-3"
            :class="posClass(row.position)"
          >{{ row.position }}</div>

          <!-- Gap bar (visual) -->
          <div class="w-20 flex-shrink-0 mr-3">
            <div v-if="row.isPlayer" class="flex items-center justify-center">
              <span class="text-[10px] font-mono font-bold text-r-blue tracking-widest">● YOU</span>
            </div>
            <div v-else class="flex flex-col items-end gap-0.5">
              <span
                class="font-mono font-bold text-sm tabular-nums leading-none"
                :class="row.gapSeconds !== null
                  ? (row.gapSeconds > 0 ? 'text-r-green' : 'text-r-accent')
                  : 'text-r-dim'"
              >{{ formatGap(row) }}</span>
              <!-- Tiny gap bar -->
              <div class="w-16 h-1 bg-r-surface rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="row.gapSeconds !== null && row.gapSeconds > 0 ? 'bg-r-green' : 'bg-r-accent'"
                  :style="{ width: gapBarWidth(row) + '%', marginLeft: row.gapSeconds !== null && row.gapSeconds < 0 ? 'auto' : '0' }"
                />
              </div>
            </div>
          </div>

          <!-- Car number badge -->
          <div class="w-10 flex-shrink-0 mr-3">
            <span class="font-mono text-xs font-bold text-r-muted bg-r-surface px-1.5 py-0.5 rounded">
              #{{ row.carNumber }}
            </span>
          </div>

          <!-- Driver name + car -->
          <div class="flex-1 min-w-0">
            <div
              class="font-mono text-sm font-semibold truncate"
              :class="row.isPlayer ? 'text-r-blue' : 'text-r-text'"
            >{{ row.userName }}</div>
            <div class="font-mono text-[10px] text-r-muted truncate">{{ row.carScreenName }}</div>
          </div>

          <!-- iRating (right side) -->
          <div class="flex-shrink-0 ml-2 text-right hidden sm:block">
            <div class="font-mono text-[10px] text-r-dim">iR</div>
            <div class="font-mono text-xs text-r-muted">{{ row.iRating > 0 ? row.iRating.toLocaleString() : '—' }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Legend -->
    <div class="flex items-center justify-center gap-6 text-[10px] font-mono text-r-muted">
      <div class="flex items-center gap-1.5">
        <div class="w-2 h-2 rounded-full bg-r-green" />
        <span>Ahead of you</span>
      </div>
      <div class="flex items-center gap-1.5">
        <div class="w-2 h-2 rounded-full bg-r-accent" />
        <span>Behind you</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry } = useIRacing()
const { driverByIdx } = useSession()

const window = ref(30) // show ±30s by default
const windows = [
  { label: '±10s', value: 10 },
  { label: '±30s', value: 30 },
  { label: 'All',  value: 9999 },
]

interface RelativeRow {
  carIdx:        number
  position:      number
  gapSeconds:    number | null
  lapDiff:       number
  isPlayer:      boolean
  carNumber:     string
  userName:      string
  carScreenName: string
  iRating:       number
}

const rows = computed<RelativeRow[]>(() => {
  const t = telemetry.value
  if (!t?.CarIdxPosition?.length) return []

  const playerIdx  = t.PlayerCarIdx ?? 0
  const playerEst  = t.CarIdxEstTime?.[playerIdx]  ?? 0
  const playerLap  = t.CarIdxLap?.[playerIdx]      ?? 0

  const all: RelativeRow[] = t.CarIdxPosition
    .map((pos, idx) => {
      if (pos <= 0) return null

      const driver = driverByIdx(idx)
      if (driver?.isSpectator) return null

      const est     = t.CarIdxEstTime?.[idx]  ?? 0
      const lap     = t.CarIdxLap?.[idx]      ?? 0
      const lapDiff = lap - playerLap

      let gapSeconds: number | null = null
      if (idx === playerIdx) {
        gapSeconds = 0
      } else if (lapDiff === 0) {
        // same lap: positive = car ahead of player (lower est time means further along track)
        gapSeconds = playerEst - est
      } else {
        // different lap – show lap diff, no second gap
        gapSeconds = null
      }

      return {
        carIdx:        idx,
        position:      pos,
        gapSeconds,
        lapDiff,
        isPlayer:      idx === playerIdx,
        carNumber:     driver?.carNumber ?? String(idx),
        userName:      driver?.userName  ?? `Car #${idx}`,
        carScreenName: driver?.carScreenName ?? '',
        iRating:       driver?.iRating ?? 0,
      } satisfies RelativeRow
    })
    .filter((r): r is RelativeRow => r !== null)
    .sort((a, b) => a.position - b.position)

  // Filter by time window around player
  const playerRow = all.find(r => r.isPlayer)
  if (!playerRow || window.value === 9999) return all

  return all.filter(r => {
    if (r.isPlayer) return true
    if (r.lapDiff !== 0) return Math.abs(r.lapDiff) <= 1
    return r.gapSeconds !== null && Math.abs(r.gapSeconds) <= window.value
  })
})

const visibleCars = computed(() => rows.value.length)

function formatGap(row: RelativeRow): string {
  if (row.isPlayer) return ''
  if (row.lapDiff !== 0) {
    return row.lapDiff > 0 ? `+${row.lapDiff}L` : `${row.lapDiff}L`
  }
  if (row.gapSeconds === null) return '—'
  const sign = row.gapSeconds > 0 ? '+' : ''
  return `${sign}${row.gapSeconds.toFixed(2)}s`
}

function gapBarWidth(row: RelativeRow): number {
  if (row.gapSeconds === null) return 0
  return Math.min(100, (Math.abs(row.gapSeconds) / window.value) * 100)
}

function posClass(pos: number): string {
  if (pos === 1) return 'bg-r-gold/20 text-r-gold border border-r-gold/40'
  if (pos === 2) return 'bg-r-silver/20 text-r-silver border border-r-silver/40'
  if (pos === 3) return 'bg-r-bronze/20 text-r-bronze border border-r-bronze/40'
  return 'bg-r-surface text-r-muted border border-r-border'
}
</script>
