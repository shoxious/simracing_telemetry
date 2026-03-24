<template>
  <div class="r-card p-3 flex flex-col items-center">
    <svg :viewBox="`0 0 ${SIZE} ${SIZE}`" class="w-full max-w-[220px]" :style="`height: ${SIZE}px`" role="img" :aria-label="`Speed: ${Math.round(speed)} km/h`">
      <!-- Background arc track -->
      <path
        :d="bgArc"
        fill="none"
        stroke="#1a1a2e"
        stroke-width="14"
        stroke-linecap="round"
      />

      <!-- Speed arc (colored gradient zones) -->
      <!-- Zone 1: 0-200 km/h (blue) -->
      <path
        :d="arcForRange(0, Math.min(speedKmh, 200), 0, maxSpeed)"
        fill="none"
        stroke="#4cc9f0"
        stroke-width="10"
        stroke-linecap="round"
        class="transition-all duration-75"
      />
      <!-- Zone 2: 200-280 km/h (yellow) -->
      <path
        v-if="speedKmh > 200"
        :d="arcForRange(200, Math.min(speedKmh, 280), 0, maxSpeed)"
        fill="none"
        stroke="#ffd166"
        stroke-width="10"
        stroke-linecap="round"
        class="transition-all duration-75"
      />
      <!-- Zone 3: 280+ km/h (red) -->
      <path
        v-if="speedKmh > 280"
        :d="arcForRange(280, Math.min(speedKmh, maxSpeed), 0, maxSpeed)"
        fill="none"
        stroke="#e63946"
        stroke-width="10"
        stroke-linecap="round"
        class="transition-all duration-75"
      />

      <!-- Tick marks at 50 km/h intervals -->
      <line
        v-for="tick in speedTicks"
        :key="tick.spd"
        :x1="tick.x1" :y1="tick.y1"
        :x2="tick.x2" :y2="tick.y2"
        stroke="#3a3a55"
        :stroke-width="tick.major ? 2 : 1"
      />

      <!-- Center digital speed -->
      <text
        :x="CX"
        :y="CY - 6"
        text-anchor="middle"
        dominant-baseline="middle"
        :font-size="speedKmh >= 100 ? 42 : 46"
        font-family="JetBrains Mono, monospace"
        font-weight="600"
        :fill="speedColor"
        class="tabular-nums"
      >{{ Math.round(speedKmh) }}</text>

      <!-- Unit label -->
      <text
        :x="CX"
        :y="CY + 30"
        text-anchor="middle"
        font-size="12"
        font-family="JetBrains Mono, monospace"
        fill="#7070a0"
        letter-spacing="3"
      >KM/H</text>

      <!-- Speed limit indicator (optional) -->
      <circle v-if="speedKmh > 280" :cx="CX" :cy="CY + 55" r="5" fill="#e63946" opacity="0.8" />
    </svg>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  speedKmh?: number
  maxSpeed?: number
}>(), {
  speedKmh: 0,
  maxSpeed: 350,
})

const SIZE = 200
const CX = SIZE / 2
const CY = SIZE / 2
const R = 82
const START_ANGLE = 225 // degrees from 12 o'clock, clockwise
const SWEEP = 270

function polarToXY(angleDeg: number, r = R) {
  const rad = ((angleDeg - 90) * Math.PI) / 180
  return { x: CX + r * Math.cos(rad), y: CY + r * Math.sin(rad) }
}

function arcPath(startVal: number, endVal: number, min: number, max: number, r = R): string {
  if (endVal <= startVal) return ''
  const startAngle = START_ANGLE + ((startVal - min) / (max - min)) * SWEEP
  const endAngle   = START_ANGLE + ((endVal   - min) / (max - min)) * SWEEP
  const s = polarToXY(startAngle, r)
  const e = polarToXY(endAngle, r)
  const large = endAngle - startAngle > 180 ? 1 : 0
  return `M ${s.x.toFixed(2)} ${s.y.toFixed(2)} A ${r} ${r} 0 ${large} 1 ${e.x.toFixed(2)} ${e.y.toFixed(2)}`
}

const bgArc = computed(() => arcPath(0, props.maxSpeed, 0, props.maxSpeed))

function arcForRange(from: number, to: number, min: number, max: number) {
  return arcPath(from, to, min, max)
}

const speedTicks = computed(() => {
  const ticks = []
  const step = 50
  for (let spd = 0; spd <= props.maxSpeed; spd += step) {
    const angle = START_ANGLE + (spd / props.maxSpeed) * SWEEP
    const rOuter = R + 12
    const rInner = R + 5
    const outer = polarToXY(angle, rOuter)
    const inner = polarToXY(angle, rInner)
    ticks.push({ spd, x1: inner.x, y1: inner.y, x2: outer.x, y2: outer.y, major: spd % 100 === 0 })
  }
  return ticks
})

const speed = computed(() => props.speedKmh)

const speedColor = computed(() => {
  if (props.speedKmh > 280) return '#e63946'
  if (props.speedKmh > 200) return '#ffd166'
  return '#e8e8f0'
})
</script>
