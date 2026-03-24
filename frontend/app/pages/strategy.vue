<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">

    <!-- Header -->
    <div class="r-card p-4">
      <div class="flex items-center justify-between mb-3">
        <span class="r-label">Race Strategy</span>
        <span class="text-xs font-mono text-r-muted">Session: {{ fmtTime(telemetry?.SessionTime) }}</span>
      </div>

      <!-- Strategy summary -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[10px] mb-1">Current Lap</div>
          <div class="font-mono font-bold text-2xl text-r-text">{{ telemetry?.Lap ?? '--' }}</div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[10px] mb-1">Fuel Left</div>
          <div class="font-mono font-bold text-2xl" :class="fuelPct < 0.15 ? 'text-r-accent' : 'text-r-green'">
            {{ telemetry?.FuelLevel?.toFixed(1) ?? '--' }}<span class="text-r-muted text-sm">L</span>
          </div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[10px] mb-1">Laps on Fuel</div>
          <div class="font-mono font-bold text-2xl" :class="fuelLaps < 5 ? 'text-r-accent animate-pulse' : 'text-r-text'">
            {{ fuelLaps }}
          </div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center border" :class="pitAlert ? 'border-r-accent/60 bg-r-accent/10' : 'border-transparent'">
          <div class="r-label text-[10px] mb-1">Pit Window</div>
          <div class="font-mono font-bold text-2xl" :class="pitAlert ? 'text-r-accent' : 'text-r-text'">
            {{ pitWindowText }}
          </div>
        </div>
      </div>
    </div>

    <!-- Fuel gauge (large) -->
    <GaugesFuelGauge
      :fuel-level="telemetry?.FuelLevel ?? 0"
      :fuel-max="55"
      :fuel-use-per-hour="telemetry?.FuelUsePerHour ?? 0"
      :laps-remaining="fuelLaps"
    />

    <!-- Fuel per lap history chart (bar chart via canvas) -->
    <div class="r-card p-4">
      <div class="r-label mb-3">Fuel Per Lap</div>
      <canvas
        ref="fuelCanvas"
        class="w-full rounded-lg bg-r-surface"
        :width="800"
        :height="120"
      />
      <div class="mt-2 flex justify-between text-[10px] font-mono text-r-muted">
        <span>Older laps →</span>
        <span>Avg: {{ avgFuelPerLap.toFixed(2) }}L/lap</span>
      </div>
    </div>

    <!-- Lap history -->
    <TimingLapTable :laps="laps" />

    <!-- Settings override -->
    <div class="r-card p-4">
      <div class="r-label mb-3">Strategy Parameters</div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="r-label text-[10px] block mb-1.5">Tank Size (L)</label>
          <input
            v-model.number="tankSize"
            type="number"
            min="10" max="100" step="0.5"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
          />
        </div>
        <div>
          <label class="r-label text-[10px] block mb-1.5">Target Lap Time (s)</label>
          <input
            v-model.number="targetLapTime"
            type="number"
            min="30" max="300"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry, laps, fmtTime, fuelLaps } = useIRacing()

const tankSize = ref(55)
const targetLapTime = ref(90)
const fuelCanvas = ref<HTMLCanvasElement | null>(null)

const fuelPct = computed(() => (telemetry.value?.FuelLevel ?? 0) / tankSize.value)

const avgFuelPerLap = computed(() => {
  const usedLaps = laps.value.filter(l => l.fuel_used > 0)
  if (!usedLaps.length) {
    // Estimate from use per hour
    const uph = telemetry.value?.FuelUsePerHour ?? 0
    return uph > 0 ? uph * (targetLapTime.value / 3600) : 0
  }
  return usedLaps.reduce((s, l) => s + l.fuel_used, 0) / usedLaps.length
})

const pitWindowText = computed(() => {
  if (fuelLaps.value <= 0) return 'PIT NOW'
  if (fuelLaps.value <= 3) return `IN ${fuelLaps.value} LAPS`
  return `${fuelLaps.value} LAPS`
})

const pitAlert = computed(() => fuelLaps.value <= 3)

// Draw fuel per lap bar chart
watch([laps, fuelCanvas], () => drawFuelChart())
onMounted(() => drawFuelChart())

function drawFuelChart() {
  const ctx = fuelCanvas.value?.getContext('2d')
  if (!ctx) return

  const W = 800
  const H = 120
  ctx.clearRect(0, 0, W, H)
  ctx.fillStyle = '#0f0f17'
  ctx.fillRect(0, 0, W, H)

  const data = [...laps.value]
    .sort((a, b) => a.lap_number - b.lap_number)
    .slice(-40)
    .map(l => l.fuel_used)
    .filter(v => v > 0)

  if (!data.length) {
    ctx.fillStyle = '#7070a0'
    ctx.font = '12px JetBrains Mono'
    ctx.textAlign = 'center'
    ctx.fillText('No lap data yet', W / 2, H / 2)
    return
  }

  const maxVal = Math.max(...data, avgFuelPerLap.value * 1.3)
  const barW = Math.max(8, (W - 20) / data.length - 3)

  data.forEach((val, i) => {
    const barH = (val / maxVal) * (H - 30)
    const x = 10 + i * (barW + 3)
    const y = H - 20 - barH

    // Bar
    const gradient = ctx.createLinearGradient(x, y, x, H - 20)
    gradient.addColorStop(0, val > avgFuelPerLap.value * 1.1 ? '#e63946' : '#4cc9f0')
    gradient.addColorStop(1, val > avgFuelPerLap.value * 1.1 ? '#e6394633' : '#4cc9f033')
    ctx.fillStyle = gradient
    ctx.fillRect(x, y, barW, barH)
  })

  // Average line
  const avgY = H - 20 - (avgFuelPerLap.value / maxVal) * (H - 30)
  ctx.strokeStyle = '#06d6a0'
  ctx.lineWidth = 1.5
  ctx.setLineDash([4, 4])
  ctx.beginPath()
  ctx.moveTo(10, avgY)
  ctx.lineTo(W - 10, avgY)
  ctx.stroke()
  ctx.setLineDash([])
}
</script>
