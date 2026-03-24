<template>
  <header class="safe-top sticky top-0 z-50 bg-r-surface/95 backdrop-blur border-b border-r-border px-3 py-2 flex items-center justify-between gap-2">

    <!-- Left: Logo -->
    <div class="flex items-center gap-2 min-w-0 flex-shrink-0">
      <div class="w-6 h-6 rounded bg-r-accent flex-shrink-0 flex items-center justify-center">
        <svg viewBox="0 0 16 16" fill="currentColor" class="w-3.5 h-3.5 text-white">
          <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 2a5 5 0 110 10A5 5 0 018 3zm0 2a3 3 0 100 6 3 3 0 000-6z"/>
        </svg>
      </div>
      <span class="font-mono font-bold text-sm text-r-text tracking-tight hidden sm:block">SimRacing</span>
      <span v-if="simulate" class="px-1.5 py-0.5 rounded text-[9px] font-mono bg-r-yellow/20 text-r-yellow border border-r-yellow/30">SIM</span>
    </div>

    <!-- Center: Live race data (position + best lap + flag) — only when connected -->
    <div v-if="connected" class="flex items-center gap-2 flex-1 justify-center min-w-0">

      <!-- Position badge -->
      <div
        class="flex items-center gap-1 px-2.5 py-1 rounded-lg border font-mono font-black text-sm leading-none"
        :class="posClass"
      >
        {{ posLabel }}
      </div>

      <!-- Divider -->
      <div class="w-px h-5 bg-r-border hidden sm:block" />

      <!-- Best lap -->
      <div class="hidden sm:flex items-center gap-1.5">
        <!-- Purple dot = fastest personal lap indicator -->
        <svg viewBox="0 0 10 10" class="w-2.5 h-2.5 flex-shrink-0" :class="hasBest ? 'text-r-purple' : 'text-r-dim'" fill="currentColor">
          <circle cx="5" cy="5" r="5"/>
        </svg>
        <span class="font-mono text-xs" :class="hasBest ? 'text-r-purple' : 'text-r-dim'">
          {{ hasBest ? fmtTime(bestLap) : '--:--.---' }}
        </span>
        <span class="font-mono text-[9px] text-r-muted uppercase tracking-widest">BEST</span>
      </div>

      <!-- Divider -->
      <div class="w-px h-5 bg-r-border hidden md:block" />

      <!-- Flag indicator -->
      <div class="hidden md:flex items-center gap-1.5">
        <div class="w-2.5 h-2.5 rounded-sm" :class="flagClass" />
        <span class="text-[10px] font-mono text-r-muted">{{ flagLabel }}</span>
      </div>
    </div>

    <!-- Right: FPS + status -->
    <div class="flex items-center gap-2 flex-shrink-0">
      <span v-if="connected" class="hidden lg:block text-xs font-mono text-r-dim">
        {{ store.fps }}<span class="text-r-dim/60">fps</span>
      </span>
      <div class="flex items-center gap-1.5">
        <div
          class="w-2 h-2 rounded-full status-dot transition-colors duration-500"
          :class="connected ? 'bg-r-green text-r-green' : 'bg-r-accent text-r-accent'"
        />
        <span class="text-xs font-mono font-semibold" :class="connected ? 'text-r-green' : 'text-r-accent'">
          {{ connected ? 'LIVE' : 'OFFLINE' }}
        </span>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useIRacingStore } from '~/stores/iracing'

const store = useIRacingStore()
const { connected, simulate } = storeToRefs(store)
const { isGreen, isYellow, isRed, position, fmtTime } = useIRacing()

const bestLap  = computed(() => store.telemetry?.LapBestLapTime ?? -1)
const hasBest  = computed(() => bestLap.value > 0)

// Position label + styling
const posLabel = computed(() => {
  const p = position.value
  return p > 0 ? `P${p}` : 'P--'
})

const posClass = computed(() => {
  const p = position.value
  if (p === 1) return 'text-r-gold  border-r-gold/40  bg-r-gold/10'
  if (p === 2) return 'text-r-silver border-r-silver/40 bg-r-silver/10'
  if (p === 3) return 'text-r-bronze border-r-bronze/40 bg-r-bronze/10'
  return 'text-r-text border-r-border bg-r-card'
})

// Flag
const flagClass = computed(() => {
  if (isRed.value)    return 'bg-red-500'
  if (isYellow.value) return 'bg-r-yellow'
  if (isGreen.value)  return 'bg-r-green'
  return 'bg-r-dim'
})
const flagLabel = computed(() => {
  if (isRed.value)    return 'Red Flag'
  if (isYellow.value) return 'Yellow Flag'
  if (isGreen.value)  return 'Green Flag'
  return 'No Flag'
})
</script>
