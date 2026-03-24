<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <div class="r-label">Session</div>

    <!-- Position + Delta row -->
    <div class="flex items-center justify-between">
      <div class="text-center">
        <div class="r-label text-[9px] mb-0.5">POS</div>
        <div class="font-mono font-black text-3xl leading-none" :class="position > 0 ? 'text-r-text' : 'text-r-dim'">
          P{{ position > 0 ? position : '--' }}
        </div>
      </div>
      <div class="text-center">
        <div class="r-label text-[9px] mb-0.5">DELTA</div>
        <div
          class="font-mono font-bold text-xl leading-none tabular-nums"
          :class="delta > 0 ? 'text-r-accent' : delta < -0.01 ? 'text-r-green' : 'text-r-muted'"
        >
          {{ delta >= 0 ? '+' : '' }}{{ delta.toFixed(3) }}
        </div>
      </div>
      <div class="text-center">
        <div class="r-label text-[9px] mb-0.5">LAP</div>
        <div class="font-mono font-bold text-2xl leading-none text-r-text">
          {{ lap > 0 ? lap : '--' }}
        </div>
      </div>
    </div>

    <!-- Divider -->
    <div class="h-px bg-r-border" />

    <!-- Track progress bar -->
    <div>
      <div class="flex justify-between text-[9px] font-mono text-r-muted mb-1">
        <span>Track</span>
        <span>{{ (lapDistPct * 100).toFixed(1) }}%</span>
      </div>
      <div class="h-2 bg-r-surface rounded-full overflow-hidden border border-r-border">
        <div
          class="h-full rounded-full transition-all duration-300"
          :style="{ width: `${lapDistPct * 100}%`, background: 'linear-gradient(to right, #4cc9f0, #06d6a0)' }"
        />
      </div>
    </div>

    <!-- Weather row -->
    <div class="grid grid-cols-2 gap-2">
      <div class="bg-r-surface rounded-xl p-2.5 flex items-center gap-2">
        <!-- Track icon -->
        <svg viewBox="0 0 16 16" class="w-4 h-4 text-r-yellow flex-shrink-0" fill="currentColor">
          <path d="M8 2a6 6 0 100 12A6 6 0 008 2zm0 1.5a4.5 4.5 0 110 9 4.5 4.5 0 010-9z"/>
        </svg>
        <div>
          <div class="r-label text-[9px]">Track</div>
          <div class="font-mono text-xs font-semibold text-r-text">{{ trackTemp.toFixed(0) }}°C</div>
        </div>
      </div>
      <div class="bg-r-surface rounded-xl p-2.5 flex items-center gap-2">
        <!-- Air icon -->
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
    <div class="text-center">
      <div class="r-label text-[9px] mb-0.5">Session Time</div>
      <div class="font-mono text-sm text-r-muted">{{ fmtTime(sessionTime) }}</div>
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
}>(), {
  position: 0, delta: 0, lap: 0, lapDistPct: 0,
  trackTemp: 0, airTemp: 0, sessionTime: 0,
})

function fmtTime(secs: number) {
  if (!secs || secs <= 0) return '--:--'
  const h = Math.floor(secs / 3600)
  const m = Math.floor((secs % 3600) / 60)
  const s = Math.floor(secs % 60)
  if (h > 0) return `${h}:${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`
  return `${m}:${String(s).padStart(2,'0')}`
}
</script>
