<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <div class="r-label">Driver Inputs</div>

    <!-- Throttle -->
    <div class="flex items-center gap-3">
      <span class="r-label w-3 text-r-green">T</span>
      <div class="flex-1 relative h-7 bg-r-surface rounded-lg overflow-hidden border border-r-border">
        <div
          class="absolute inset-y-0 left-0 bg-r-green rounded-lg transition-all duration-50"
          :style="{ width: `${throttle * 100}%` }"
        />
        <div
          class="absolute inset-y-0 left-0 bg-r-green/20 rounded-lg"
          :style="{ width: '100%' }"
        />
        <span class="absolute inset-0 flex items-center justify-end pr-2 font-mono text-xs font-semibold text-white/80">
          {{ Math.round(throttle * 100) }}%
        </span>
      </div>
    </div>

    <!-- Brake -->
    <div class="flex items-center gap-3">
      <span class="r-label w-3 text-r-accent">B</span>
      <div class="flex-1 relative h-7 bg-r-surface rounded-lg overflow-hidden border border-r-border">
        <div
          class="absolute inset-y-0 left-0 bg-r-accent rounded-lg transition-all duration-50"
          :style="{ width: `${brake * 100}%` }"
        />
        <div
          class="absolute inset-y-0 left-0 bg-r-accent/10 rounded-lg"
          :style="{ width: '100%' }"
        />
        <span class="absolute inset-0 flex items-center justify-end pr-2 font-mono text-xs font-semibold text-white/80">
          {{ Math.round(brake * 100) }}%
        </span>
      </div>
    </div>

    <!-- Clutch -->
    <div class="flex items-center gap-3">
      <span class="r-label w-3 text-r-blue">C</span>
      <div class="flex-1 relative h-4 bg-r-surface rounded-lg overflow-hidden border border-r-border">
        <div
          class="absolute inset-y-0 left-0 bg-r-blue/70 rounded-lg transition-all duration-50"
          :style="{ width: `${clutch * 100}%` }"
        />
      </div>
    </div>

    <!-- Steering angle bar (centred) -->
    <div class="flex items-center gap-3">
      <span class="r-label w-3 text-r-muted">S</span>
      <div class="flex-1 relative h-4 bg-r-surface rounded-lg overflow-hidden border border-r-border">
        <!-- Center line -->
        <div class="absolute inset-y-0 left-1/2 w-px bg-r-border" />
        <!-- Steering indicator -->
        <div
          class="absolute top-1 bottom-1 rounded-sm bg-r-muted/60 transition-all duration-50"
          :style="steeringStyle"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  throttle?: number
  brake?: number
  clutch?: number
  steeringAngle?: number  // radians, ±0.5 typical
}>(), {
  throttle: 0,
  brake: 0,
  clutch: 0,
  steeringAngle: 0,
})

const steeringStyle = computed(() => {
  // Normalise ±1.0 rad to ±50% of bar
  const maxRad = 1.0
  const pct = Math.max(-1, Math.min(1, props.steeringAngle / maxRad))
  const BAR_WIDTH = 8 // percent

  if (pct >= 0) {
    // Right turn: grows from center to right
    return {
      left: '50%',
      width: `${Math.abs(pct) * 50}%`,
    }
  } else {
    // Left turn: grows from center to left
    return {
      right: `${50 + pct * 50}%`,
      left:  `${50 + pct * 50}%`,
      width: `${Math.abs(pct) * 50}%`,
    }
  }
})
</script>
