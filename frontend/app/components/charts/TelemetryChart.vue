<template>
  <div class="r-card p-4 flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <span class="r-label">Telemetry Trace</span>
      <div class="flex items-center gap-3 text-[10px] font-mono">
        <span class="text-r-green">— Throttle</span>
        <span class="text-r-accent">— Brake</span>
        <span class="text-r-blue">— Speed</span>
      </div>
    </div>

    <canvas
      ref="canvas"
      class="w-full rounded-lg bg-r-surface"
      :width="canvasW"
      :height="canvasH"
      style="image-rendering: pixelated"
    />

    <!-- Window selector -->
    <div class="flex gap-2 justify-end">
      <button
        v-for="w in [10, 30, 60]"
        :key="w"
        @click="windowSec = w"
        class="px-2.5 py-1 rounded text-[10px] font-mono border transition-colors"
        :class="windowSec === w
          ? 'bg-r-accent border-r-accent text-white'
          : 'border-r-border text-r-muted hover:text-r-text'"
      >{{ w }}s</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useIRacingStore } from '~/stores/iracing'

const store = useIRacingStore()
const canvas = ref<HTMLCanvasElement | null>(null)
const windowSec = ref(30)
const canvasW = 800
const canvasH = 180

// Ring buffers (60 Hz × 60s max = 3600 points)
const MAX_POINTS = 3600
const throttleRing = new Float32Array(MAX_POINTS)
const brakeRing    = new Float32Array(MAX_POINTS)
const speedRing    = new Float32Array(MAX_POINTS)
let head = 0
let count = 0

// Watch telemetry and fill ring buffers
watch(() => store.telemetry, (t) => {
  if (!t) return
  throttleRing[head] = t.Throttle
  brakeRing[head]    = t.Brake
  speedRing[head]    = (t.Speed * 3.6) / 350 // normalised 0-1 against 350 km/h
  head = (head + 1) % MAX_POINTS
  if (count < MAX_POINTS) count++
})

let rafId = 0

onMounted(() => {
  draw()
})

onUnmounted(() => {
  cancelAnimationFrame(rafId)
})

function draw() {
  rafId = requestAnimationFrame(draw)

  const ctx = canvas.value?.getContext('2d')
  if (!ctx) return

  const W = canvasW
  const H = canvasH
  ctx.clearRect(0, 0, W, H)

  // Background
  ctx.fillStyle = '#0f0f17'
  ctx.fillRect(0, 0, W, H)

  // Grid lines
  ctx.strokeStyle = '#1a1a2e'
  ctx.lineWidth = 1
  for (let y of [0.25, 0.5, 0.75]) {
    ctx.beginPath()
    ctx.moveTo(0, y * H)
    ctx.lineTo(W, y * H)
    ctx.stroke()
  }
  // Center line (50%)
  ctx.strokeStyle = '#22223a'
  ctx.setLineDash([4, 4])
  ctx.beginPath()
  ctx.moveTo(0, 0.5 * H)
  ctx.lineTo(W, 0.5 * H)
  ctx.stroke()
  ctx.setLineDash([])

  const points = Math.min(count, windowSec.value * 60)
  if (points < 2) return

  // Draw trace
  function drawTrace(ring: Float32Array, color: string, alpha: number) {
    ctx.beginPath()
    ctx.strokeStyle = color
    ctx.globalAlpha = alpha
    ctx.lineWidth = 1.5
    ctx.lineJoin = 'round'

    for (let i = 0; i < points; i++) {
      const idx = (head - points + i + MAX_POINTS) % MAX_POINTS
      const x = (i / (points - 1)) * W
      const y = H - ring[idx] * H * 0.9 - H * 0.05
      if (i === 0) ctx.moveTo(x, y)
      else ctx.lineTo(x, y)
    }
    ctx.stroke()
    ctx.globalAlpha = 1
  }

  // Speed (blue, background)
  drawTrace(speedRing, '#4cc9f0', 0.6)
  // Throttle (green)
  drawTrace(throttleRing, '#06d6a0', 0.9)
  // Brake (red)
  drawTrace(brakeRing, '#e63946', 0.9)

  // Current values on right edge
  const lastIdx = (head - 1 + MAX_POINTS) % MAX_POINTS
  const labels = [
    { val: throttleRing[lastIdx], color: '#06d6a0', label: `${Math.round(throttleRing[lastIdx] * 100)}%` },
    { val: brakeRing[lastIdx],    color: '#e63946', label: `${Math.round(brakeRing[lastIdx] * 100)}%` },
    { val: speedRing[lastIdx],    color: '#4cc9f0', label: `${Math.round(speedRing[lastIdx] * 350)}` },
  ]
  ctx.font = '11px JetBrains Mono, monospace'
  ctx.textAlign = 'right'
  labels.forEach(l => {
    const y = H - l.val * H * 0.9 - H * 0.05
    ctx.fillStyle = l.color
    ctx.fillText(l.label, W - 4, Math.max(12, Math.min(H - 4, y - 2)))
  })
}
</script>
