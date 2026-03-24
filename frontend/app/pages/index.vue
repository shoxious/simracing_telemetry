<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">

    <!-- Shift light bar -->
    <GaugesShiftLightBar
      :rpm="telemetry?.RPM ?? 0"
      :max-rpm="8600"
    />

    <!-- Row 1: Speed + RPM gauges -->
    <div class="grid grid-cols-2 gap-3">
      <GaugesSpeedGauge :speed-kmh="speedKmh" :max-speed="350" />
      <GaugesRpmGauge :rpm="telemetry?.RPM ?? 0" :gear="gearLabel" :max-rpm="8600" />
    </div>

    <!-- Row 2: Pedals (vertical) + Session card -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <InputsPedalBars
        :throttle="telemetry?.Throttle ?? 0"
        :brake="telemetry?.Brake ?? 0"
        :clutch="telemetry?.Clutch ?? 0"
        :steering-angle="telemetry?.SteeringAngle ?? 0"
      />
      <SessionCard
        :position="position"
        :delta="deltaToBest"
        :lap="telemetry?.Lap ?? 0"
        :lap-dist-pct="telemetry?.LapDistPct ?? 0"
        :track-temp="telemetry?.TrackTemp ?? 0"
        :air-temp="telemetry?.AirTemp ?? 0"
        :session-time="telemetry?.SessionTime ?? 0"
        :best-lap-time="telemetry?.LapBestLapTime ?? -1"
      />
    </div>

    <!-- Row 3: Lap Timer -->
    <TimingLapTimer
      :lap="telemetry?.Lap ?? 0"
      :current-lap-time="telemetry?.LapCurrentLapTime ?? 0"
      :last-lap-time="telemetry?.LapLastLapTime ?? -1"
      :best-lap-time="telemetry?.LapBestLapTime ?? -1"
      :lap-dist-pct="telemetry?.LapDistPct ?? 0"
    />

    <!-- Row 4: Tire Widget (car footprint) -->
    <TiresTireWidget v-if="telemetry" :tires="telemetry" />
    <div v-else class="r-card p-4">
      <div class="r-label mb-3 flex items-center justify-between">
        <span>Tire Temperatures</span>
        <div class="flex items-center gap-2 text-[10px] font-mono">
          <span style="color:#4cc9f0">COLD</span>
          <div class="w-16 h-1 rounded" style="background: linear-gradient(to right,#4cc9f0,#06d6a0,#ffd166,#e63946)" />
          <span style="color:#e63946">HOT</span>
        </div>
      </div>
      <div class="flex items-center justify-center h-32 text-r-muted text-xs font-mono">
        Waiting for telemetry...
      </div>
    </div>

    <!-- Row 5: Fuel gauge -->
    <GaugesFuelGauge
      :fuel-level="telemetry?.FuelLevel ?? 0"
      :fuel-max="55"
      :fuel-use-per-hour="telemetry?.FuelUsePerHour ?? 0"
      :laps-remaining="fuelLaps"
    />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })
const { telemetry, speedKmh, gearLabel, fuelLaps, position, deltaToBest } = useIRacing()
</script>
