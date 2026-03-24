<template>
  <div class="r-card p-3 flex flex-col items-center">
    <svg :viewBox="`0 0 ${SIZE} ${SIZE}`" class="w-full max-w-[220px]" style="height: auto" role="img" :aria-label="`RPM: ${Math.round(rpm)}, Gear: ${gear}`">
      <!-- Shift lights row at top -->
      <g>
        <circle
          v-for="(dot, i) in shiftDots"
          :key="i"
          :cx="dot.cx"
          :cy="22"
          r="5"
          :fill="dot.lit ? dot.color : '#1a1a2e'"
          :class="dot.lit && shiftNow ? 'shift-flash' : ''"
        />
      </g>

      <!-- Background arc -->
      <path :d="bgArc" fill="none" stroke="#1a1a2e" stroke-width="14" stroke-linecap="round" />

      <!-- Normal RPM zone (teal) -->
      <path
        :d="arcForRange(0, Math.min(rpm, YELLOW_RPM), 0, MAX_RPM)"
        fill="none"
        stroke="#06d6a0"
        stroke-width="10"
        stroke-linecap="round"
        class="transition-all duration-75"
      />
      <!-- Yellow zone -->
      <path
        v-if="rpm > YELLOW_RPM"
        :d="arcForRange(YELLOW_RPM, Math.min(rpm, RED_RPM), 0, MAX_RPM)"
        fill="none"
        stroke="#ffd166"
        stroke-width="10"
        stroke-linecap="round"
        class="transition-all duration-75"
      />
      <!-- Redline zone -->
      <path
        v-if="rpm > RED_RPM"
        :d="arcForRange(RED_RPM, Math.min(rpm, MAX_RPM), 0, MAX_RPM)"
        fill="none"
        stroke="#e63946"
        stroke-width="14"
        stroke-linecap="round"
        :class="shiftNow ? 'shift-flash' : 'transition-all duration-75'"
      />

      <!-- Gear number (large, centred) -->
      <text
        :x="CX"
        :y="CY - 4"
        text-anchor="middle"
        dominant-baseline="middle"
        font-size="64"
        font-family="JetBrains Mono, monospace"
        font-weight="700"
        :fill="shiftNow ? '#e63946' : '#e8e8f0'"
        class="tabular-nums transition-colors duration-100"
      >{{ gear }}</text>

      <!-- RPM number (small, below gear) -->
      <text
        :x="CX"
        :y="CY + 38"
        text-anchor="middle"
        font-size="13"
        font-family="JetBrains Mono, monospace"
        :fill="rpmColor"
        class="tabular-nums"
      >{{ rpmDisplay }}</text>

      <!-- "RPM" label -->
      <text
        :x="CX"
        :y="CY + 55"
        text-anchor="middle"
        font-size="9"
        font-family="JetBrains Mono, monospace"
        fill="#7070a0"
        letter-spacing="3"
      >RPM</text>
    </svg>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  rpm?: number
  gear?: string | number
  maxRpm?: number
}>(), {
  rpm: 0,
  gear: 'N',
  maxRpm: 8600,
})

const MAX_RPM   = computed(() => props.maxRpm)
const YELLOW_RPM = computed(() => MAX_RPM.value * 0.80) // 80% = yellow
const RED_RPM    = computed(() => MAX_RPM.value * 0.90) // 90% = redline

const SIZE = 200
const CX = SIZE / 2
const CY = SIZE / 2 + 4
const R = 82
const START_ANGLE = 225
const SWEEP = 270

function polarToXY(angleDeg: number, r = R) {
  const rad = ((angleDeg - 90) * Math.PI) / 180
  return { x: CX + r * Math.cos(rad), y: CY + r * Math.sin(rad) }
}

function arcPath(startVal: number, endVal: number, min: number, max: number): string {
  if (endVal <= startVal) return ''
  const startAngle = START_ANGLE + ((startVal - min) / (max - min)) * SWEEP
  const endAngle   = START_ANGLE + ((endVal   - min) / (max - min)) * SWEEP
  const s = polarToXY(startAngle)
  const e = polarToXY(endAngle)
  const large = endAngle - startAngle > 180 ? 1 : 0
  return `M ${s.x.toFixed(2)} ${s.y.toFixed(2)} A ${R} ${R} 0 ${large} 1 ${e.x.toFixed(2)} ${e.y.toFixed(2)}`
}

const bgArc = computed(() => arcPath(0, MAX_RPM.value, 0, MAX_RPM.value))
function arcForRange(from: number, to: number, min: number, max: number) {
  return arcPath(from, to, min, max)
}

const rpm = computed(() => props.rpm ?? 0)
const shiftNow = computed(() => rpm.value >= RED_RPM.value)

const rpmDisplay = computed(() => {
  const k = rpm.value / 1000
  return k.toFixed(1) + 'k'
})

const rpmColor = computed(() => {
  if (rpm.value > RED_RPM.value)    return '#e63946'
  if (rpm.value > YELLOW_RPM.value) return '#ffd166'
  return '#7070a0'
})

// 5 shift light dots across the top of the gauge
const shiftDots = computed(() => {
  const dots = []
  const startX = CX - 40
  const step = 20
  const thresholds = [0.65, 0.72, 0.80, 0.88, 0.95]
  const colors = ['#06d6a0', '#06d6a0', '#ffd166', '#e63946', '#e63946']
  const rpmPct = rpm.value / MAX_RPM.value

  for (let i = 0; i < 5; i++) {
    dots.push({
      cx: startX + i * step,
      lit: rpmPct >= thresholds[i],
      color: colors[i],
    })
  }
  return dots
})
</script>
