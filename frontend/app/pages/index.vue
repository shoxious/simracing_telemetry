<template>
  <div class="grid grid-cols-2 md:grid-cols-3 gap-2 p-2 md:p-3 animate-fade-in">

    <!-- Shift light bar — full width on both layouts -->
    <GaugesShiftLightBar
      class="col-span-2 md:col-span-3"
      :rpm="telemetry?.RPM ?? 0"
      :max-rpm="8600"
    />

    <!-- Speed gauge -->
    <GaugesSpeedGauge :speed-kmh="speedKmh" :max-speed="350" />

    <!-- RPM gauge -->
    <GaugesRpmGauge :rpm="telemetry?.RPM ?? 0" :gear="gearLabel" :max-rpm="8600" />

    <!-- Lap timer — full width on mobile, 1 col on tablet (fills the 3rd slot in row 2) -->
    <TimingLapTimer
      class="col-span-2 md:col-span-1"
      :lap="telemetry?.Lap ?? 0"
      :current-lap-time="telemetry?.LapCurrentLapTime ?? 0"
      :last-lap-time="telemetry?.LapLastLapTime ?? -1"
      :best-lap-time="telemetry?.LapBestLapTime ?? -1"
      :lap-dist-pct="telemetry?.LapDistPct ?? 0"
    />

    <!-- Pedals — full width on mobile, 1 col on tablet -->
    <InputsPedalBars
      class="col-span-2 md:col-span-1"
      :throttle="telemetry?.Throttle ?? 0"
      :brake="telemetry?.Brake ?? 0"
      :clutch="telemetry?.Clutch ?? 0"
      :steering-angle="telemetry?.SteeringAngle ?? 0"
    />

    <!-- Tire widget — full width on mobile, 1 col on tablet -->
    <TiresTireWidget
      v-if="telemetry"
      class="col-span-2 md:col-span-1"
      :tires="telemetry"
    />
    <div v-else class="r-card p-4 col-span-2 md:col-span-1">
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

    <!-- Fuel gauge — full width on mobile, 1 col on tablet -->
    <GaugesFuelGauge
      class="col-span-2 md:col-span-1"
      :fuel-level="telemetry?.FuelLevel ?? 0"
      :fuel-max="55"
      :fuel-use-per-hour="telemetry?.FuelUsePerHour ?? 0"
      :laps-remaining="fuelLaps"
    />

    <!-- Session card — only on mobile (tablet has AppStatusBar) -->
    <SessionCard
      class="col-span-2 md:hidden"
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
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })
const { telemetry, speedKmh, gearLabel, fuelLaps, position, deltaToBest } = useIRacing()
</script>
