<template>
  <div class="r-card p-4 flex flex-col gap-3">

    <!-- Top row: Position (large) + Lap + Delta -->
    <div class="grid grid-cols-3 gap-2">

      <!-- Position — big, colored by rank -->
      <div
        class="flex flex-col items-center justify-center rounded-xl py-3 border"
        :class="posCardClass"
      >
        <div class="r-label text-[8px] mb-0.5">POS</div>
        <div class="font-mono font-black text-4xl leading-none tabular-nums" :class="posTextClass">
          {{ position > 0 ? `P${position}` : 'P--' }}
        </div>
      </div>

      <!-- Current Lap -->
      <div class="flex flex-col items-center justify-center bg-r-surface rounded-xl py-3">
        <div class="r-label text-[8px] mb-0.5">LAP</div>
        <div class="font-mono font-black text-4xl leading-none text-r-text">
          {{ lap > 0 ? lap : '--' }}
        </div>
      </div>

      <!-- Delta to best -->
      <div class="flex flex-col items-center justify-center bg-r-surface rounded-xl py-3">
        <div class="r-label text-[8px] mb-0.5">DELTA</div>
        <div
          class="font-mono font-bold text-xl leading-none tabular-nums"
          :class="delta > 0.05 ? 'text-r-accent' : delta < -0.05 ? 'text-r-green' : 'text-r-muted'"
        >
          {{ delta >= 0 ? '+' : '' }}{{ delta.toFixed(3) }}
        </div>
      </div>
    </div>

    <!-- Best lap row — prominent purple -->
    <div
      v-if="bestLapTime > 0"
      class="flex items-center justify-between px-3 py-2 rounded-xl border border-r-purple/30 bg-r-purple/5"
    >
      <div class="flex items-center gap-2">
        <svg viewBox="0 0 12 12" fill="currentColor" class="w-3 h-3 text-r-purple flex-shrink-0">
          <path d="M6 0l1.6 3.4L11 4 8.5 6.5l.6 3.5L6 8.5 2.9 10l.6-3.5L1 4l3.4-.6z"/>
        </svg>
        <span class="font-mono text-[10px] tracking-widest text-r-purple uppercase">Best Lap</span>
      </div>
      <span class="font-mono font-bold text-r-purple text-sm tabular-nums">{{ fmtLap(bestLapTime) }}</span>
    </div>

    <!-- Divider -->
    <div class="h-px bg-r-border" />

    <!-- Track progress -->
    <div>
      <div class="flex justify-between text-[9px] font-mono text-r-muted mb-1.5">
        <span>Track Progress</span>
        <span>{{ (lapDistPct * 100).toFixed(1) }}%</span>
      </div>
      <div class="h-2 bg-r-surface rounded-full overflow-hidden border border-r-border">
        <div
          class="h-full rounded-full transition-all duration-300"
          :style="{ width: `${lapDistPct * 100}%`, background: 'linear-gradient(to right,#4cc9f0,#06d6a0)' }"
        />
      </div>
    </div>

    <!-- Weather row -->
    <div class="grid grid-cols-2 gap-2">
      <div class="bg-r-surface rounded-xl p-2.5 flex items-center gap-2">
        <svg viewBox="0 0 16 16" class="w-4 h-4 text-r-yellow flex-shrink-0" fill="currentColor">
          <path d="M8 2a6 6 0 100 12A6 6 0 008 2zm0 1.5a4.5 4.5 0 110 9 4.5 4.5 0 010-9z"/>
        </svg>
        <div>
          <div class="r-label text-[9px]">Track</div>
          <div class="font-mono text-xs font-semibold text-r-text">{{ trackTemp.toFixed(0) }}°C</div>
        </div>
      </div>
      <div class="bg-r-surface rounded-xl p-2.5 flex items-center gap-2">
        <svg viewBox="0 0 16 16" class="w-4 h-4 text-r-blue flex-shrink-0" fill="currentColor">
          <path d="M12.5 8a2.5 2.5 0 01-2.5 2.5H1.5a.5.5 0 010-1H10a1.5 1.5 0 000-3H1.5a.5.5 0 010-1H10a2.5 2.5 0 012.5 2.5z"/>
        </svg>
        <div>
          <div class="r-label text-[9px]">Air</div>
          <div class="font-mono text-xs font-semibold text-r-text">{{ airTemp.toFixed(0) }}°C</div>
        </div>
      </div>
    </div>

    <!-- Session time -->
    <div class="text-center text-[10px] font-mono text-r-dim">
      Session: {{ fmtSession(sessionTime) }}
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  position?: number
  delta?: number
  lap?: number
  lapDistPct?: number
  trackTemp?: number
  airTemp?: number
  sessionTime?: number
  bestLapTime?: number
}>(), {
  position: 0, delta: 0, lap: 0, lapDistPct: 0,
  trackTemp: 0, airTemp: 0, sessionTime: 0, bestLapTime: -1,
})

// Position card background + text color based on rank
const posCardClass = computed(() => {
  if (props.position === 1) return 'border-r-gold/40   bg-r-gold/10'
  if (props.position === 2) return 'border-r-silver/40 bg-r-silver/10'
  if (props.position === 3) return 'border-r-bronze/40 bg-r-bronze/10'
  return 'border-r-border bg-r-surface'
})

const posTextClass = computed(() => {
  if (props.position === 1) return 'text-r-gold'
  if (props.position === 2) return 'text-r-silver'
  if (props.position === 3) return 'text-r-bronze'
  return 'text-r-text'
})

function fmtLap(secs: number) {
  if (!secs || secs <= 0) return '--:--.---'
  const m = Math.floor(secs / 60)
  const s = secs % 60
  return `${m}:${s.toFixed(3).padStart(6, '0')}`
}

function fmtSession(secs: number) {
  if (!secs || secs <= 0) return '--:--'
  const h = Math.floor(secs / 3600)
  const m = Math.floor((secs % 3600) / 60)
  const s = Math.floor(secs % 60)
  if (h > 0) return `${h}:${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`
  return `${m}:${String(s).padStart(2,'0')}`
}
</script>
