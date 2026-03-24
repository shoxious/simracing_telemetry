<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">
    <!-- Lap Timer (top, always visible) -->
    <TimingLapTimer
      :lap="telemetry?.Lap ?? 0"
      :current-lap-time="telemetry?.LapCurrentLapTime ?? 0"
      :last-lap-time="telemetry?.LapLastLapTime ?? -1"
      :best-lap-time="telemetry?.LapBestLapTime ?? -1"
      :lap-dist-pct="telemetry?.LapDistPct ?? 0"
    />

    <!-- Timing Tower -->
    <div class="r-card overflow-hidden">
      <div class="px-4 py-3 border-b border-r-border flex items-center justify-between">
        <span class="r-label">Race Standings</span>
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 rounded-full bg-r-green animate-pulse" />
          <span class="text-xs font-mono text-r-green">LIVE</span>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-sm font-mono">
          <thead>
            <tr class="border-b border-r-border">
              <th class="px-3 py-2 text-left r-label text-[10px]">POS</th>
              <th class="px-3 py-2 text-left r-label text-[10px]">DRIVER</th>
              <th class="px-3 py-2 text-right r-label text-[10px]">LAP</th>
              <th class="px-3 py-2 text-right r-label text-[10px]">DIST</th>
              <th class="px-3 py-2 text-right r-label text-[10px]">GAP</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="car in standings"
              :key="car.carIdx"
              class="border-b border-r-border/50 transition-colors"
              :class="[
                car.isPlayer ? 'bg-r-blue/10 border-l-2 border-r-blue' : 'hover:bg-r-surface/40',
              ]"
            >
              <td class="px-3 py-3">
                <div
                  class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold"
                  :class="posClass(car.position)"
                >{{ car.position }}</div>
              </td>
              <td class="px-3 py-3">
                <div class="flex items-center gap-2">
                  <span class="text-r-text" :class="car.isPlayer ? 'font-bold text-r-blue' : ''">
                    {{ car.name }}
                  </span>
                  <span v-if="car.isPlayer" class="px-1 py-0.5 rounded text-[9px] bg-r-blue/20 text-r-blue border border-r-blue/30">YOU</span>
                </div>
                <div class="text-[10px] text-r-muted">{{ car.car }}</div>
              </td>
              <td class="px-3 py-3 text-right tabular-nums text-r-muted">{{ car.lap }}</td>
              <td class="px-3 py-3 text-right tabular-nums text-r-muted">
                {{ (car.lapDist * 100).toFixed(0) }}%
              </td>
              <td class="px-3 py-3 text-right tabular-nums" :class="car.isPlayer ? 'text-r-blue' : 'text-r-muted'">
                {{ car.gap }}
              </td>
            </tr>
            <tr v-if="!standings.length">
              <td colspan="5" class="px-4 py-10 text-center text-r-muted text-xs">
                Waiting for session data...
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Lap History Table -->
    <TimingLapTable :laps="laps" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry, laps, fmtTime, position } = useIRacing()

// Parse session YAML to get driver names
const store = useIRacingStore()
const driverNames = computed(() => {
  const yaml = store.sessionYaml
  if (!yaml) return []
  // Quick regex extraction for names
  const matches = [...yaml.matchAll(/UserName: (.+)/g)]
  return matches.map(m => m[1].trim())
})

const carNames = computed(() => {
  const yaml = store.sessionYaml
  if (!yaml) return []
  const matches = [...yaml.matchAll(/CarScreenName: (.+)/g)]
  return matches.map(m => m[1].trim())
})

const standings = computed(() => {
  const t = telemetry.value
  if (!t?.CarIdxPosition?.length) return []

  const playerIdx = t.PlayerCarIdx ?? 0

  return t.CarIdxPosition
    .map((pos, idx) => ({
      carIdx: idx,
      position: pos,
      lapDist: t.CarIdxLapDist?.[idx] ?? 0,
      lap: t.CarIdxLap?.[idx] ?? 0,
      estTime: t.CarIdxEstTime?.[idx] ?? 0,
      isPlayer: idx === playerIdx,
      name: driverNames.value[idx] ?? `Car #${idx}`,
      car: carNames.value[idx] ?? '',
    }))
    .filter(c => c.position > 0)
    .sort((a, b) => a.position - b.position)
    .map(c => ({
      ...c,
      gap: c.isPlayer ? '---' : formatGap(c.estTime),
    }))
})

function formatGap(estTime: number): string {
  if (estTime === 0) return 'L +1'
  const sign = estTime < 0 ? '+' : '-'
  return `${sign}${Math.abs(estTime).toFixed(3)}`
}

function posClass(pos: number): string {
  if (pos === 1) return 'bg-r-yellow text-black'
  if (pos === 2) return 'bg-gray-400 text-black'
  if (pos === 3) return 'bg-amber-600 text-white'
  return 'bg-r-surface text-r-muted border border-r-border'
}
</script>
