<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">
    <!-- Telemetry trace chart -->
    <ChartsTelemetryChart />

    <!-- Detailed gauges row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <!-- Speed -->
      <div class="r-card p-3 text-center">
        <div class="r-label text-[10px] mb-1">Speed</div>
        <div class="font-mono font-bold text-2xl tabular-nums" :class="speedColor">
          {{ speedKmh.toFixed(1) }}
        </div>
        <div class="r-label text-[10px]">KM/H</div>
      </div>

      <!-- RPM -->
      <div class="r-card p-3 text-center">
        <div class="r-label text-[10px] mb-1">Engine RPM</div>
        <div class="font-mono font-bold text-2xl tabular-nums" :class="rpmColor">
          {{ Math.round(telemetry?.RPM ?? 0) }}
        </div>
        <div class="r-label text-[10px]">RPM</div>
      </div>

      <!-- Gear -->
      <div class="r-card p-3 text-center">
        <div class="r-label text-[10px] mb-1">Gear</div>
        <div class="font-mono font-black text-4xl tabular-nums text-r-text">
          {{ gearLabel }}
        </div>
      </div>

      <!-- Position -->
      <div class="r-card p-3 text-center">
        <div class="r-label text-[10px] mb-1">Position</div>
        <div class="font-mono font-black text-4xl tabular-nums text-r-accent">
          P{{ position > 0 ? position : '--' }}
        </div>
      </div>
    </div>

    <!-- Inputs detail -->
    <InputsPedalBars
      :throttle="telemetry?.Throttle ?? 0"
      :brake="telemetry?.Brake ?? 0"
      :clutch="telemetry?.Clutch ?? 0"
      :steering-angle="telemetry?.SteeringAngle ?? 0"
    />

    <!-- G-Forces (derived) -->
    <div class="r-card p-4">
      <div class="r-label mb-3">Tire Temperatures Detail</div>
      <TiresTireWidget v-if="telemetry" :tires="telemetry" />
    </div>

    <!-- Live data table -->
    <div class="r-card overflow-hidden">
      <div class="px-4 py-3 border-b border-r-border">
        <span class="r-label">Live Data</span>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 divide-r-border">
        <div
          v-for="row in dataRows"
          :key="row.label"
          class="px-4 py-3 border-b border-r-border/60 hover:bg-r-surface/40"
        >
          <div class="r-label text-[10px]">{{ row.label }}</div>
          <div class="font-mono text-sm text-r-text tabular-nums mt-0.5">{{ row.value }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry, speedKmh, gearLabel, position, fmtTime, rpmPct } = useIRacing()

const speedColor = computed(() => {
  const kmh = speedKmh.value
  if (kmh > 280) return 'text-r-accent'
  if (kmh > 200) return 'text-r-yellow'
  return 'text-r-blue'
})

const rpmColor = computed(() => {
  const pct = rpmPct.value
  if (pct > 0.9) return 'text-r-accent'
  if (pct > 0.8) return 'text-r-yellow'
  return 'text-r-green'
})

const dataRows = computed(() => {
  const t = telemetry.value
  if (!t) return []
  return [
    { label: 'Speed (m/s)',     value: t.Speed?.toFixed(2) ?? '--' },
    { label: 'Speed (km/h)',    value: speedKmh.value.toFixed(1) },
    { label: 'RPM',             value: Math.round(t.RPM ?? 0).toString() },
    { label: 'Gear',            value: gearLabel.value },
    { label: 'Throttle',        value: `${((t.Throttle ?? 0) * 100).toFixed(0)}%` },
    { label: 'Brake',           value: `${((t.Brake ?? 0) * 100).toFixed(0)}%` },
    { label: 'Fuel Level',      value: `${t.FuelLevel?.toFixed(2) ?? '--'} L` },
    { label: 'Fuel Use/hr',     value: `${t.FuelUsePerHour?.toFixed(2) ?? '--'} L` },
    { label: 'Lap',             value: t.Lap?.toString() ?? '--' },
    { label: 'Lap Dist %',      value: `${((t.LapDistPct ?? 0) * 100).toFixed(1)}%` },
    { label: 'Current Lap',     value: fmtTime(t.LapCurrentLapTime) },
    { label: 'Last Lap',        value: fmtTime(t.LapLastLapTime) },
    { label: 'Best Lap',        value: fmtTime(t.LapBestLapTime) },
    { label: 'Track Temp',      value: `${t.TrackTemp?.toFixed(1) ?? '--'}°C` },
    { label: 'Air Temp',        value: `${t.AirTemp?.toFixed(1) ?? '--'}°C` },
    { label: 'Session Time',    value: fmtTime(t.SessionTime) },
  ]
})
</script>
