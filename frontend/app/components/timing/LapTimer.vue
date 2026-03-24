<template>
  <div class="r-card p-4 flex flex-col gap-3">

    <!-- Header row: label + lap badge -->
    <div class="flex items-center justify-between">
      <div class="r-label">Lap Times</div>
      <div class="px-3 py-1 bg-r-surface rounded-full border border-r-border">
        <span class="r-label mr-1">LAP</span>
        <span class="font-mono font-bold text-r-text">{{ lap > 0 ? lap : '--' }}</span>
      </div>
    </div>

    <!-- FASTEST LAP banner — flashes when a new personal best is set -->
    <Transition name="best-fade">
      <div
        v-if="showBestBanner"
        class="flex items-center justify-center gap-2 py-2 px-3 rounded-xl border border-r-purple/50 bg-r-purple/10 animate-pulse"
      >
        <svg viewBox="0 0 16 16" fill="currentColor" class="w-3.5 h-3.5 text-r-purple flex-shrink-0">
          <path d="M8 1l2.2 4.4L15 6.3l-3.5 3.4.8 4.8L8 12.1l-4.3 2.4.8-4.8L1 6.3l4.8-.9z"/>
        </svg>
        <span class="font-mono font-bold text-sm text-r-purple tracking-widest">FASTEST LAP</span>
        <svg viewBox="0 0 16 16" fill="currentColor" class="w-3.5 h-3.5 text-r-purple flex-shrink-0">
          <path d="M8 1l2.2 4.4L15 6.3l-3.5 3.4.8 4.8L8 12.1l-4.3 2.4.8-4.8L1 6.3l4.8-.9z"/>
        </svg>
      </div>
    </Transition>

    <!-- Current lap prominent display -->
    <div class="flex items-end gap-3">
      <div class="flex-1">
        <div class="r-label text-[9px] mb-1 tracking-widest">CURRENT</div>
        <div
          class="font-mono text-4xl font-black tabular-nums tracking-tight transition-colors duration-200 leading-none"
          :class="currentClass"
        >{{ fmtTime(currentLapTime) }}</div>
      </div>
      <!-- Delta bubble -->
      <div
        v-if="deltaToBest !== 0 && bestLapTime > 0"
        class="px-3 py-1.5 rounded-xl font-mono font-bold text-base tabular-nums mb-1 border transition-colors duration-200"
        :class="deltaToBest < 0
          ? 'text-r-green border-r-green/40 bg-r-green/10'
          : deltaToBest < 0.5
            ? 'text-r-yellow border-r-yellow/40 bg-r-yellow/10'
            : 'text-r-accent border-r-accent/40 bg-r-accent/10'"
      >
        {{ deltaToBest > 0 ? '+' : '' }}{{ deltaToBest.toFixed(3) }}
      </div>
    </div>

    <!-- Divider -->
    <div class="h-px bg-r-border" />

    <!-- Last / Best side by side -->
    <div class="grid grid-cols-2 gap-2">

      <!-- Last lap -->
      <div class="bg-r-surface rounded-xl p-3 border transition-colors duration-300"
        :class="isLastBest ? 'border-r-purple/40' : 'border-transparent'">
        <div class="r-label text-[9px] mb-1.5 tracking-widest">LAST</div>
        <div
          class="font-mono text-lg font-semibold tabular-nums leading-none transition-colors duration-200"
          :class="isLastBest ? 'text-r-purple' : 'text-r-text'"
        >
          {{ lastLapTime > 0 ? fmtTime(lastLapTime) : '--:--.---' }}
        </div>
        <div class="mt-1 font-mono text-[10px] tabular-nums"
          :class="isLastBest ? 'text-r-purple' : lastLapTime > 0 && bestLapTime > 0 ? 'text-r-muted' : 'text-r-dim'">
          {{ isLastBest ? '★ NEW BEST' : lastLapTime > 0 && bestLapTime > 0 ? '+' + (lastLapTime - bestLapTime).toFixed(3) : '---' }}
        </div>
      </div>

      <!-- Personal best lap -->
      <div class="bg-r-surface rounded-xl p-3 border border-r-purple/30">
        <div class="r-label text-[9px] mb-1.5 tracking-widest text-r-purple">BEST</div>
        <div class="font-mono text-lg font-bold tabular-nums leading-none text-r-purple">
          {{ bestLapTime > 0 ? fmtTime(bestLapTime) : '--:--.---' }}
        </div>
        <div class="mt-1 font-mono text-[10px] text-r-purple/50">Personal</div>
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
          :class="deltaToBest < 0 ? 'bg-r-green' : deltaToBest < 0.5 ? 'bg-r-yellow' : 'bg-r-accent'"
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

// True when last lap equals or beats best lap
const isLastBest = computed(() =>
  props.lastLapTime > 0 && props.bestLapTime > 0 &&
  props.lastLapTime <= props.bestLapTime
)

// Show FASTEST LAP banner for 8s after a new best
const showBestBanner = ref(false)
let bannerTimer: ReturnType<typeof setTimeout> | null = null

watch(isLastBest, (val) => {
  if (!val) return
  showBestBanner.value = true
  if (bannerTimer) clearTimeout(bannerTimer)
  bannerTimer = setTimeout(() => { showBestBanner.value = false }, 8000)
})

// Delta projection to best
const deltaToBest = computed(() => {
  const best = props.bestLapTime
  const cur  = props.currentLapTime
  const pct  = props.lapDistPct
  if (best <= 0 || pct < 0.02) return 0
  return (cur / pct) - best
})

const currentClass = computed(() => {
  if (props.bestLapTime <= 0) return 'text-r-text'
  if (deltaToBest.value < -0.1) return 'text-r-green glow-green'
  if (deltaToBest.value < 0.5)  return 'text-r-yellow'
  return 'text-r-accent'
})
</script>

<style scoped>
.best-fade-enter-active,
.best-fade-leave-active { transition: opacity 0.4s ease, transform 0.4s ease; }
.best-fade-enter-from,
.best-fade-leave-to     { opacity: 0; transform: scaleY(0.8); }
</style>
