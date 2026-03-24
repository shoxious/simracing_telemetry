<template>
  <header class="safe-top sticky top-0 z-50 bg-r-surface/95 backdrop-blur border-b border-r-border px-4 py-2 flex items-center justify-between gap-3">
    <!-- Logo & track name -->
    <div class="flex items-center gap-2 min-w-0">
      <div class="w-6 h-6 rounded bg-r-accent flex-shrink-0 flex items-center justify-center">
        <svg viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 text-white">
          <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm0 2a5 5 0 110 10A5 5 0 018 3zm0 2a3 3 0 100 6 3 3 0 000-6z"/>
        </svg>
      </div>
      <span class="font-mono font-bold text-sm text-r-text tracking-tight truncate">
        SimRacing
      </span>
      <span v-if="simulate" class="px-1.5 py-0.5 rounded text-[10px] font-mono bg-r-yellow/20 text-r-yellow border border-r-yellow/30">
        SIM
      </span>
    </div>

    <!-- Flag indicator -->
    <div v-if="connected" class="hidden sm:flex items-center gap-1.5">
      <div
        class="w-3 h-3 rounded-sm transition-colors duration-300"
        :class="flagClass"
      />
      <span class="text-xs font-mono text-r-muted">{{ flagLabel }}</span>
    </div>

    <!-- Right: FPS + connection status -->
    <div class="flex items-center gap-3 flex-shrink-0">
      <span v-if="connected" class="hidden md:block text-xs font-mono text-r-muted">
        {{ store.fps }}<span class="text-r-dim">fps</span>
      </span>

      <div class="flex items-center gap-1.5">
        <div
          class="w-2 h-2 rounded-full status-dot transition-colors duration-500"
          :class="connected ? 'bg-r-green text-r-green' : 'bg-r-accent text-r-accent'"
        />
        <span class="text-xs font-mono" :class="connected ? 'text-r-green' : 'text-r-accent'">
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
const { isGreen, isYellow, isRed } = useIRacing()

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
