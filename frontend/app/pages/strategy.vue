<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">

    <!-- ── Pit Stop Calculator ─────────────────────────────────────────────── -->
    <div class="r-card p-4">
      <div class="r-label mb-3 flex items-center gap-2">
        <svg viewBox="0 0 16 16" fill="currentColor" class="w-3.5 h-3.5 text-r-accent">
          <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm.5 3.5v4.25l3 1.75-.75 1.25-3.5-2.25V4.5h1.25z"/>
        </svg>
        Pit Stop Calculator
      </div>

      <!-- Inputs row -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-4">
        <div>
          <label class="r-label text-[9px] block mb-1.5">Race Laps Total</label>
          <input
            v-model.number="totalLaps"
            type="number" min="1" max="500"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
            placeholder="e.g. 50"
          />
        </div>
        <div>
          <label class="r-label text-[9px] block mb-1.5">Tank Size (L)</label>
          <input
            v-model.number="tankSize"
            type="number" min="10" max="120" step="0.5"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
          />
        </div>
        <div>
          <label class="r-label text-[9px] block mb-1.5">Target Lap Time (s)</label>
          <input
            v-model.number="targetLapTime"
            type="number" min="30" max="600"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
          />
        </div>
        <div>
          <label class="r-label text-[9px] block mb-1.5">Safety Margin (L)</label>
          <input
            v-model.number="safetyMargin"
            type="number" min="0" max="5" step="0.1"
            class="w-full bg-r-surface border border-r-border rounded-lg px-3 py-2 font-mono text-sm text-r-text focus:border-r-blue outline-none"
          />
        </div>
      </div>

      <!-- Live stats row -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-4">
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[9px] mb-1">Current Lap</div>
          <div class="font-mono font-bold text-xl text-r-text">{{ telemetry?.Lap ?? '--' }}</div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[9px] mb-1">Fuel Left</div>
          <div class="font-mono font-bold text-xl" :class="fuelPct < 0.15 ? 'text-r-accent' : 'text-r-green'">
            {{ telemetry?.FuelLevel?.toFixed(1) ?? '--' }}<span class="text-r-muted text-sm">L</span>
          </div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[9px] mb-1">Avg Fuel/Lap</div>
          <div class="font-mono font-bold text-xl text-r-text">
            {{ avgFuelPerLap > 0 ? avgFuelPerLap.toFixed(2) : '--' }}<span class="text-r-muted text-sm">L</span>
          </div>
        </div>
        <div class="bg-r-surface rounded-xl p-3 text-center">
          <div class="r-label text-[9px] mb-1">Laps on Tank</div>
          <div class="font-mono font-bold text-xl" :class="lapsOnCurrentFuel < 5 ? 'text-r-accent animate-pulse' : 'text-r-text'">
            {{ lapsOnCurrentFuel > 0 ? lapsOnCurrentFuel.toFixed(1) : '--' }}
          </div>
        </div>
      </div>

      <!-- Fuel bar -->
      <GaugesFuelGauge
        :fuel-level="telemetry?.FuelLevel ?? 0"
        :fuel-max="tankSize"
        :fuel-use-per-hour="telemetry?.FuelUsePerHour ?? 0"
        :laps-remaining="fuelLaps"
      />

      <!-- ── Strategy result ─────────────────────────────────────────────── -->
      <div v-if="plan.length" class="mt-4 space-y-2">
        <div class="r-label text-[9px] mb-2">Pit Strategy</div>
        <div
          v-for="(stop, i) in plan"
          :key="i"
          class="flex items-center gap-3 px-4 py-3 rounded-xl border transition-colors"
          :class="stop.urgent
            ? 'border-r-accent/50 bg-r-accent/10'
            : 'border-r-border bg-r-surface'"
        >
          <!-- Stop number -->
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center font-mono font-bold text-sm flex-shrink-0"
            :class="stop.urgent ? 'bg-r-accent/20 text-r-accent' : 'bg-r-blue/20 text-r-blue'"
          >
            {{ i + 1 }}
          </div>

          <div class="flex-1 min-w-0">
            <div class="font-mono text-sm text-r-text font-semibold">
              Pit after lap {{ stop.pitAfterLap }}
              <span v-if="stop.urgent" class="ml-2 text-[10px] text-r-accent font-bold tracking-widest">⚠ URGENT</span>
            </div>
            <div class="font-mono text-[11px] text-r-muted mt-0.5">
              Add {{ stop.fuelToAdd.toFixed(1) }}L → {{ stop.fuelAfter.toFixed(1) }}L in tank
            </div>
          </div>

          <div class="text-right flex-shrink-0">
            <div class="font-mono text-xs text-r-muted">Laps to go</div>
            <div class="font-mono font-bold text-sm text-r-text">{{ stop.lapsAfter }}</div>
          </div>
        </div>

        <!-- Finish summary -->
        <div class="flex items-center gap-3 px-4 py-3 rounded-xl border border-r-green/30 bg-r-green/5">
          <div class="w-8 h-8 rounded-full bg-r-green/20 flex items-center justify-center flex-shrink-0">
            <svg viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 text-r-green">
              <path d="M13 3H3v2h2v8h2V9h2v4h2V5h2V3z"/>
            </svg>
          </div>
          <div class="flex-1">
            <div class="font-mono text-sm text-r-green font-semibold">Finish — lap {{ totalLaps }}</div>
            <div class="font-mono text-[11px] text-r-muted">~{{ finishFuel.toFixed(1) }}L remaining</div>
          </div>
          <div class="font-mono text-sm font-bold text-r-green">{{ plan.length }} pit{{ plan.length !== 1 ? 's' : '' }}</div>
        </div>
      </div>

      <!-- No data hint -->
      <div v-else-if="totalLaps > 0 && avgFuelPerLap <= 0" class="mt-4 text-center text-r-muted text-xs font-mono py-4">
        Complete at least 1 lap to calculate fuel strategy
      </div>
    </div>

    <!-- ── Fuel per lap chart ──────────────────────────────────────────────── -->
    <div class="r-card p-4">
      <div class="r-label mb-3 flex items-center justify-between">
        <span>Fuel Per Lap</span>
        <span class="text-[10px] font-mono text-r-muted">avg {{ avgFuelPerLap.toFixed(2) }} L/lap</span>
      </div>
      <canvas ref="fuelCanvas" class="w-full rounded-lg bg-r-surface" :width="800" :height="120" />
    </div>

    <!-- ── Lap history ─────────────────────────────────────────────────────── -->
    <TimingLapTable :laps="laps" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry, laps, fuelLaps } = useIRacing()

const totalLaps      = ref(0)
const tankSize       = ref(55)
const targetLapTime  = ref(90)
const safetyMargin   = ref(0.5)
const fuelCanvas     = ref<HTMLCanvasElement | null>(null)

const fuelPct = computed(() =>
  (telemetry.value?.FuelLevel ?? 0) / tankSize.value
)

// Average fuel per lap: prefer lap history, fall back to FuelUsePerHour estimate
const avgFuelPerLap = computed(() => {
  const usedLaps = laps.value.filter(l => l.fuel_used > 0)
  if (usedLaps.length >= 1) {
    // Use last 5 laps for a rolling average (more responsive to fuel map changes)
    const recent = usedLaps.slice(-5)
    return recent.reduce((s, l) => s + l.fuel_used, 0) / recent.length
  }
  const uph = telemetry.value?.FuelUsePerHour ?? 0
  return uph > 0 ? uph * (targetLapTime.value / 3600) : 0
})

const lapsOnCurrentFuel = computed(() => {
  const fuel = telemetry.value?.FuelLevel ?? 0
  if (avgFuelPerLap.value <= 0) return 0
  return fuel / avgFuelPerLap.value
})

// ── Pit strategy plan ────────────────────────────────────────────────────────
interface PitStop {
  pitAfterLap: number
  fuelToAdd:   number
  fuelAfter:   number
  lapsAfter:   number
  urgent:      boolean
}

const plan = computed<PitStop[]>(() => {
  if (totalLaps.value <= 0 || avgFuelPerLap.value <= 0) return []

  const currentLap  = telemetry.value?.Lap ?? 0
  const currentFuel = telemetry.value?.FuelLevel ?? 0
  const lapsLeft    = totalLaps.value - currentLap
  if (lapsLeft <= 0) return []

  const fuelPerLap = avgFuelPerLap.value + safetyMargin.value / Math.max(lapsLeft, 1)
  const stops: PitStop[] = []
  let fuel = currentFuel
  let lap  = currentLap

  while (lap < totalLaps.value) {
    // How many laps can we do on current fuel?
    const lapsCanDo = fuel / fuelPerLap
    const lapAfterBurn = lap + Math.floor(lapsCanDo)

    if (lapAfterBurn >= totalLaps.value) break // can finish!

    // Pit at the last possible lap before running out, minus 1 for safety
    const pitLap    = Math.max(lap + 1, lapAfterBurn)
    const remaining = totalLaps.value - pitLap
    // How much fuel do we need to reach the finish?
    const fuelNeeded = Math.min(tankSize.value, remaining * fuelPerLap + safetyMargin.value)
    const fuelAtPit  = Math.max(0, fuel - (pitLap - lap) * fuelPerLap)
    const fuelToAdd  = Math.max(0, fuelNeeded - fuelAtPit)
    const fuelAfter  = Math.min(tankSize.value, fuelAtPit + fuelToAdd)

    stops.push({
      pitAfterLap: pitLap,
      fuelToAdd,
      fuelAfter,
      lapsAfter:   remaining,
      urgent:      pitLap - currentLap <= 3,
    })

    fuel = fuelAfter
    lap  = pitLap
  }

  return stops
})

const finishFuel = computed(() => {
  if (totalLaps.value <= 0 || avgFuelPerLap.value <= 0) return 0
  const currentFuel = telemetry.value?.FuelLevel ?? 0
  const currentLap  = telemetry.value?.Lap ?? 0
  const lapsLeft    = totalLaps.value - currentLap

  if (plan.value.length === 0) {
    // No pit stops needed
    return Math.max(0, currentFuel - lapsLeft * avgFuelPerLap.value)
  }

  const lastStop = plan.value[plan.value.length - 1]
  const lapsAfterLastPit = totalLaps.value - lastStop.pitAfterLap
  return Math.max(0, lastStop.fuelAfter - lapsAfterLastPit * avgFuelPerLap.value)
})

// ── Fuel per lap bar chart ────────────────────────────────────────────────────
watch([laps, fuelCanvas], () => drawFuelChart(), { deep: true })
onMounted(() => drawFuelChart())

function drawFuelChart() {
  const ctx = fuelCanvas.value?.getContext('2d')
  if (!ctx) return

  const W = 800, H = 120
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
    ctx.font = '12px JetBrains Mono, monospace'
    ctx.textAlign = 'center'
    ctx.fillText('No lap data yet', W / 2, H / 2)
    return
  }

  const avg    = data.reduce((a, b) => a + b, 0) / data.length
  const maxVal = Math.max(...data, avg * 1.3)
  const barW   = Math.max(8, (W - 20) / data.length - 3)

  data.forEach((val, i) => {
    const barH = (val / maxVal) * (H - 30)
    const x    = 10 + i * (barW + 3)
    const y    = H - 20 - barH
    const hot  = val > avg * 1.1

    const g = ctx.createLinearGradient(x, y, x, H - 20)
    g.addColorStop(0, hot ? '#e63946' : '#4cc9f0')
    g.addColorStop(1, hot ? '#e6394633' : '#4cc9f033')
    ctx.fillStyle = g
    ctx.fillRect(x, y, barW, barH)

    // Lap number label
    if (barW >= 12) {
      ctx.fillStyle = '#7070a0'
      ctx.font = '8px JetBrains Mono, monospace'
      ctx.textAlign = 'center'
      ctx.fillText(String(laps.value[i]?.lap_number ?? i + 1), x + barW / 2, H - 6)
    }
  })

  // Average line
  const avgY = H - 20 - (avg / maxVal) * (H - 30)
  ctx.strokeStyle = '#06d6a0'
  ctx.lineWidth = 1.5
  ctx.setLineDash([4, 4])
  ctx.beginPath()
  ctx.moveTo(10, avgY)
  ctx.lineTo(W - 10, avgY)
  ctx.stroke()
  ctx.setLineDash([])

  // Avg label
  ctx.fillStyle = '#06d6a0'
  ctx.font = '9px JetBrains Mono, monospace'
  ctx.textAlign = 'right'
  ctx.fillText(`avg ${avg.toFixed(2)}L`, W - 12, avgY - 4)
}
</script>
