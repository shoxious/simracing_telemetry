<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <!-- Current lap time (prominent) -->
    <div class="text-center">
      <div class="r-label mb-1">Current Lap</div>
      <div
        class="font-mono text-3xl font-bold tabular-nums tracking-tight transition-colors duration-200"
        :class="deltaClass"
      >{{ fmtTime(currentLapTime) }}</div>
      <div v-if="deltaToBest !== 0" class="mt-0.5 font-mono text-sm" :class="deltaClass">
        {{ deltaToBest > 0 ? '+' : '' }}{{ deltaToBest.toFixed(3) }}
      </div>
    </div>

    <!-- Lap number indicator -->
    <div class="flex justify-center">
      <div class="px-3 py-1 bg-r-surface rounded-full border border-r-border">
        <span class="r-label mr-1">LAP</span>
        <span class="font-mono font-bold text-r-text">{{ lap }}</span>
      </div>
    </div>

    <!-- Last / Best row -->
    <div class="grid grid-cols-2 gap-2">
      <div class="bg-r-surface rounded-xl p-3 text-center">
        <div class="r-label text-[10px] mb-1">Last Lap</div>
        <div class="font-mono text-sm font-semibold text-r-text tabular-nums">
          {{ fmtTime(lastLapTime) }}
        </div>
      </div>
      <div class="bg-r-surface rounded-xl p-3 text-center border border-r-purple/30">
        <div class="r-label text-[10px] mb-1 text-r-purple">Best Lap</div>
        <div class="font-mono text-sm font-bold tabular-nums" style="color: #c084fc">
          {{ fmtTime(bestLapTime) }}
        </div>
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
