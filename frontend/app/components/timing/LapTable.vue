<template>
  <div class="r-card overflow-hidden">
    <div class="px-4 py-3 border-b border-r-border flex items-center justify-between">
      <span class="r-label">Lap History</span>
      <span class="text-xs font-mono text-r-muted">{{ laps.length }} laps</span>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-sm font-mono">
        <thead>
          <tr class="border-b border-r-border">
            <th class="px-3 py-2 text-left r-label text-[10px]">LAP</th>
            <th class="px-3 py-2 text-right r-label text-[10px]">TIME</th>
            <th class="px-3 py-2 text-right r-label text-[10px]">FUEL</th>
            <th class="px-3 py-2 text-right r-label text-[10px] hidden sm:table-cell">MAX V</th>
            <th class="px-3 py-2 text-right r-label text-[10px] hidden md:table-cell">THR%</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="lap in displayLaps"
            :key="lap.lap_number"
            class="border-b border-r-border/50 transition-colors hover:bg-r-surface/60"
            :class="lap.isBest ? 'bg-r-purple/10' : ''"
          >
            <td class="px-3 py-2.5 text-r-muted">{{ lap.lap_number }}</td>
            <td class="px-3 py-2.5 text-right tabular-nums" :class="lap.timeClass">
              {{ fmtTime(lap.lap_time) }}
            </td>
            <td class="px-3 py-2.5 text-right tabular-nums text-r-text">
              {{ lap.fuel_used?.toFixed(2) ?? '--' }}L
            </td>
            <td class="px-3 py-2.5 text-right tabular-nums text-r-muted hidden sm:table-cell">
              {{ lap.max_speed ? Math.round(lap.max_speed) : '--' }}
            </td>
            <td class="px-3 py-2.5 text-right tabular-nums text-r-muted hidden md:table-cell">
              {{ lap.avg_throttle ? Math.round(lap.avg_throttle * 100) + '%' : '--' }}
            </td>
          </tr>
          <tr v-if="!laps.length">
            <td colspan="5" class="px-4 py-8 text-center text-r-muted text-xs">
              No laps recorded yet
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LapRecord } from '~/stores/iracing'

const props = defineProps<{ laps: LapRecord[] }>()
const { fmtTime } = useIRacing()

const displayLaps = computed(() => {
  const sorted = [...props.laps].sort((a, b) => b.lap_number - a.lap_number)
  const bestTime = Math.min(...sorted.filter(l => l.lap_time > 0).map(l => l.lap_time))

  return sorted.map(lap => ({
    ...lap,
    isBest: lap.lap_time === bestTime && bestTime > 0,
    timeClass: lap.lap_time === bestTime && bestTime > 0
      ? 'text-purple-400 font-bold'
      : 'text-r-text',
  }))
})
</script>
