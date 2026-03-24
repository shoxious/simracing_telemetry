<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <span class="r-label">Fuel</span>
      <div class="flex items-center gap-3">
        <span class="r-value text-base" :class="fuelColor">
          {{ fuelLevel.toFixed(1) }}<span class="text-r-muted text-xs ml-0.5">L</span>
        </span>
        <span class="text-xs font-mono text-r-muted">
          ~{{ lapsRemaining }}<span class="text-r-dim">lap</span>
        </span>
      </div>
    </div>

    <!-- Fuel bar -->
    <div class="relative h-4 bg-r-surface rounded-full overflow-hidden border border-r-border">
      <!-- Segmented tick marks -->
      <div
        v-for="tick in fuelTicks"
        :key="tick"
        class="absolute top-0 bottom-0 w-px bg-r-bg/40"
        :style="{ left: `${tick}%` }"
      />

      <!-- Fuel fill -->
      <div
        class="absolute inset-y-0 left-0 rounded-full transition-all duration-300"
        :class="fuelBarClass"
        :style="{ width: `${fuelPct * 100}%` }"
      />

      <!-- Low fuel flash overlay -->
      <div
        v-if="fuelPct < 0.12"
        class="absolute inset-0 bg-r-accent/20 animate-pulse-fast rounded-full"
      />
    </div>

    <!-- Usage stats -->
    <div class="grid grid-cols-2 gap-2 text-center">
      <div class="bg-r-surface rounded-lg p-2">
        <div class="r-label text-[10px]">Use/Lap</div>
        <div class="r-value text-sm">{{ usePerLap.toFixed(2) }}<span class="text-r-muted text-xs">L</span></div>
      </div>
      <div class="bg-r-surface rounded-lg p-2">
        <div class="r-label text-[10px]">Use/Hr</div>
        <div class="r-value text-sm">{{ usePerHour.toFixed(1) }}<span class="text-r-muted text-xs">L</span></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  fuelLevel?: number
  fuelMax?: number
  fuelUsePerHour?: number
  lapsRemaining?: number
}>(), {
  fuelLevel: 0,
  fuelMax: 55,
  fuelUsePerHour: 0,
  lapsRemaining: 0,
})

const fuelPct = computed(() => Math.min(props.fuelLevel / props.fuelMax, 1))
const usePerHour = computed(() => props.fuelUsePerHour)
const usePerLap  = computed(() => {
  // Assume ~90s lap time
  return props.fuelUsePerHour > 0 ? props.fuelUsePerHour * (90 / 3600) : 0
})

const fuelColor = computed(() => {
  if (fuelPct.value < 0.12) return 'text-r-accent glow-red'
  if (fuelPct.value < 0.25) return 'text-r-yellow glow-yellow'
  return 'text-r-green'
})

const fuelBarClass = computed(() => {
  if (fuelPct.value < 0.12) return 'bg-r-accent'
  if (fuelPct.value < 0.25) return 'bg-r-yellow'
  return 'bg-r-green'
})

// Tick marks at 25% intervals
const fuelTicks = [25, 50, 75]
</script>
