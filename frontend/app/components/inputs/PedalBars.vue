<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <div class="r-label">Driver Inputs</div>

    <!-- Vertical pedal bars -->
    <div class="flex items-end justify-center gap-4 h-28">

      <!-- Throttle -->
      <div class="flex flex-col items-center gap-1.5 w-14">
        <span class="font-mono font-bold text-sm tabular-nums" :class="throttle > 0.5 ? 'text-r-green' : 'text-r-muted'">
          {{ Math.round(throttle * 100) }}<span class="text-[10px]">%</span>
        </span>
        <div class="relative flex-1 w-full bg-r-surface rounded-lg overflow-hidden border border-r-border" style="height: 80px">
          <div class="absolute bottom-0 inset-x-0 bg-r-green rounded-b-lg transition-all duration-50"
            :style="{ height: `${throttle * 100}%` }" />
          <div class="absolute inset-0 bg-r-green/5 rounded-lg" />
        </div>
        <span class="r-label text-[9px] tracking-widest">THR</span>
      </div>

      <!-- Brake -->
      <div class="flex flex-col items-center gap-1.5 w-14">
        <span class="font-mono font-bold text-sm tabular-nums" :class="brake > 0.1 ? 'text-r-accent' : 'text-r-muted'">
          {{ Math.round(brake * 100) }}<span class="text-[10px]">%</span>
        </span>
        <div class="relative w-full bg-r-surface rounded-lg overflow-hidden border border-r-border" style="height: 80px">
          <div class="absolute bottom-0 inset-x-0 rounded-b-lg transition-all duration-50"
            :class="brake > 0.9 ? 'bg-r-accent animate-pulse' : 'bg-r-accent'"
            :style="{ height: `${brake * 100}%` }" />
          <div class="absolute inset-0 bg-r-accent/5 rounded-lg" />
        </div>
        <span class="r-label text-[9px] tracking-widest">BRK</span>
      </div>

      <!-- Clutch -->
      <div class="flex flex-col items-center gap-1.5 w-14">
        <span class="font-mono font-bold text-sm tabular-nums" :class="clutch > 0.05 ? 'text-r-blue' : 'text-r-muted'">
          {{ Math.round(clutch * 100) }}<span class="text-[10px]">%</span>
        </span>
        <div class="relative w-full bg-r-surface rounded-lg overflow-hidden border border-r-border" style="height: 80px">
          <div class="absolute bottom-0 inset-x-0 bg-r-blue/70 rounded-b-lg transition-all duration-50"
            :style="{ height: `${clutch * 100}%` }" />
          <div class="absolute inset-0 bg-r-blue/5 rounded-lg" />
        </div>
        <span class="r-label text-[9px] tracking-widest">CLT</span>
      </div>
    </div>

    <!-- Steering arc -->
    <div>
      <div class="flex items-center justify-between mb-1">
        <span class="r-label text-[9px] tracking-widest">STEERING</span>
        <span class="font-mono text-[10px] text-r-muted">{{ steeringDeg.toFixed(0) }}°</span>
      </div>
      <div class="relative h-5 bg-r-surface rounded-full overflow-hidden border border-r-border">
        <div class="absolute inset-y-0 left-1/2 w-px bg-r-border/60 z-10" />
        <div
          class="absolute top-1 bottom-1 rounded-full bg-gradient-to-r from-r-blue/40 to-r-blue transition-all duration-50"
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
  steeringAngle?: number
}>(), { throttle: 0, brake: 0, clutch: 0, steeringAngle: 0 })

const steeringDeg = computed(() => (props.steeringAngle * 180) / Math.PI)

const steeringStyle = computed(() => {
  const maxRad = 1.0
  const pct = Math.max(-1, Math.min(1, props.steeringAngle / maxRad))
  if (pct >= 0) {
    return { left: '50%', width: `${Math.abs(pct) * 50}%` }
  } else {
    const leftPct = 50 + pct * 50
    return { left: `${leftPct}%`, width: `${Math.abs(pct) * 50}%` }
  }
})
</script>
