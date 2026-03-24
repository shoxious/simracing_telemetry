<template>
  <div class="r-card p-4 flex flex-col gap-3">

    <!-- Header row: label + lap number -->
    <div class="flex items-center justify-between">
      <div class="r-label">Lap Times</div>
      <div class="px-3 py-1 bg-r-surface rounded-full border border-r-border">
        <span class="r-label mr-1">LAP</span>
        <span class="font-mono font-bold text-r-text">{{ lap > 0 ? lap : '--' }}</span>
      </div>
    </div>

    <!-- Current lap prominent display -->
    <div class="flex items-end gap-3">
      <div class="flex-1">
        <div class="r-label text-[9px] mb-1 tracking-widest">CURRENT</div>
        <div
          class="font-mono text-4xl font-black tabular-nums tracking-tight transition-colors duration-200 leading-none"
          :class="deltaClass"
        >{{ fmtTime(currentLapTime) }}</div>
      </div>
      <!-- Delta bubble -->
      <div
        v-if="deltaToBest !== 0 && bestLapTime > 0"
        class="px-3 py-1.5 rounded-xl font-mono font-bold text-base tabular-nums mb-1 border"
        :class="deltaToBest < 0
          ? 'text-r-green border-r-green/30 bg-r-green/10'
          : deltaToBest < 0.3
            ? 'text-r-yellow border-r-yellow/30 bg-r-yellow/10'
            : 'text-r-accent border-r-accent/30 bg-r-accent/10'"
      >
        {{ deltaToBest > 0 ? '+' : '' }}{{ deltaToBest.toFixed(3) }}
      </div>
    </div>

    <!-- Divider -->
    <div class="h-px bg-r-border" />

    <!-- Last / Best side by side -->
    <div class="grid grid-cols-2 gap-2">
      <!-- Last lap -->
      <div class="bg-r-surface rounded-xl p-3">
        <div class="r-label text-[9px] mb-1.5 tracking-widest">LAST</div>
        <div class="font-mono text-lg font-semibold text-r-text tabular-nums leading-none">
          {{ lastLapTime > 0 ? fmtTime(lastLapTime) : '--:--.---' }}
        </div>
        <!-- Delta last vs best -->
        <div
          v-if="lastLapTime > 0 && bestLapTime > 0"
          class="mt-1 font-mono text-[10px] tabular-nums"
          :class="lastLapTime <= bestLapTime ? 'text-r-green' : 'text-r-muted'"
        >
          {{ lastLapTime <= bestLapTime ? '▲ BEST' : '+' + (lastLapTime - bestLapTime).toFixed(3) }}
        </div>
      </div>

      <!-- Best lap -->
      <div class="bg-r-surface rounded-xl p-3 border border-r-purple/30">
        <div class="r-label text-[9px] mb-1.5 tracking-widest text-r-purple">BEST</div>
        <div class="font-mono text-lg font-bold tabular-nums leading-none" style="color: #c084fc">
          {{ bestLapTime > 0 ? fmtTime(bestLapTime) : '--:--.---' }}
        </div>
        <div class="mt-1 font-mono text-[10px] text-r-purple/60">Personal</div>
      </div>
    </div>

    <!-- Lap progress bar -->
    <div>
      <div class="flex justify-between text-[9px] font-mono text-r-muted mb-1">
        <span>Lap progress</span>
        <span>{{ (lapDistPct * 100).toFixed(1) }}%</span>
      </div>
      <div class="h-1.5 bg-r-surface rounded-full overflow-hidden border border-r-border">
        <div
          class="h-full rounded-full transition-all duration-300"
          :class="deltaClass.includes('green') ? 'bg-r-green' : deltaClass.includes('yellow') ? 'bg-r-yellow' : 'bg-r-accent'"
          :style="{ width: `${lapDistPct * 100}%` }"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  lap?: number
  currentLapTime?: number
  lastLapTime?: number
  bestLapTime?: number
  lapDistPct?: number
}>(), {
  lap: 0,
  currentLapTime: 0,
  lastLapTime: -1,
  bestLapTime: -1,
  lapDistPct: 0,
})

const { fmtTime } = useIRacing()

// Delta to best projection
const deltaToBest = computed(() => {
  const best = props.bestLapTime
  const cur  = props.currentLapTime
  const pct  = props.lapDistPct
  if (best <= 0 || pct < 0.02) return 0
  const projected = cur / pct
  return projected - best
})

const deltaClass = computed(() => {
  if (props.bestLapTime <= 0) return 'text-r-text'
  if (deltaToBest.value < -0.1) return 'text-r-green glow-green'
  if (deltaToBest.value < 0.3)  return 'text-r-yellow'
  return 'text-r-accent'
})
</script>
