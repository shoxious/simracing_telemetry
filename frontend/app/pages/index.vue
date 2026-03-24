<template>
  <div class="p-3 md:p-4 space-y-3 max-w-screen-lg mx-auto animate-fade-in">

    <!-- Offline banner -->
    <div
      v-if="!connected"
      class="flex items-center gap-3 px-4 py-3 rounded-xl bg-r-accent/10 border border-r-accent/30 text-r-accent text-sm font-mono"
    >
      <div class="w-2 h-2 rounded-full bg-r-accent animate-pulse flex-shrink-0" />
      Waiting for iRacing connection...
    </div>

    <!-- Row 1: Speed + RPM gauges -->
    <div class="grid grid-cols-2 gap-3">
      <GaugesSpeedGauge :speed-kmh="speedKmh" :max-speed="350" />
      <GaugesRpmGauge :rpm="telemetry?.RPM ?? 0" :gear="gearLabel" :max-rpm="8600" />
    </div>

    <!-- Row 2: Pedals + Lap Timer -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <InputsPedalBars
        :throttle="telemetry?.Throttle ?? 0"
        :brake="telemetry?.Brake ?? 0"
        :clutch="telemetry?.Clutch ?? 0"
        :steering-angle="telemetry?.SteeringAngle ?? 0"
      />
      <TimingLapTimer
        :lap="telemetry?.Lap ?? 0"
        :current-lap-time="telemetry?.LapCurrentLapTime ?? 0"
        :last-lap-time="telemetry?.LapLastLapTime ?? -1"
        :best-lap-time="telemetry?.LapBestLapTime ?? -1"
        :lap-dist-pct="telemetry?.LapDistPct ?? 0"
      />
    </div>

    <!-- Row 3: Fuel + Position info -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <GaugesFuelGauge
        :fuel-level="telemetry?.FuelLevel ?? 0"
        :fuel-max="55"
        :fuel-use-per-hour="telemetry?.FuelUsePerHour ?? 0"
        :laps-remaining="fuelLaps"
      />

      <!-- Position card -->
      <div class="r-card p-4 flex flex-col justify-between gap-4">
        <div class="grid grid-cols-2 gap-4">
          <div class="text-center">
            <div class="r-label text-[10px] mb-1">Position</div>
            <div class="font-mono font-black text-5xl text-r-text tabular-nums leading-none">
              P{{ position > 0 ? position : '--' }}
            </div>
          </div>
          <div class="text-center">
            <div class="r-label text-[10px] mb-1">Lap Dist</div>
            <div class="font-mono font-bold text-2xl text-r-text tabular-nums">
              {{ ((telemetry?.LapDistPct ?? 0) * 100).toFixed(1) }}<span class="text-r-muted text-base">%</span>
            </div>
          </div>
        </div>

        <!-- Track progress bar -->
        <div>
          <div class="flex justify-between text-[10px] font-mono text-r-muted mb-1">
            <span>Track Progress</span>
            <span>{{ fmtTime(telemetry?.SessionTime) }}</span>
          </div>
          <div class="h-2 bg-r-surface rounded-full overflow-hidden border border-r-border">
            <div
              class="h-full bg-r-blue rounded-full transition-all duration-300"
              :style="{ width: `${(telemetry?.LapDistPct ?? 0) * 100}%` }"
            />
          </div>
        </div>

        <!-- Environment -->
        <div class="flex gap-2 text-center">
          <div class="flex-1 bg-r-surface rounded-lg p-2">
            <div class="r-label text-[10px]">Track</div>
            <div class="font-mono text-xs text-r-text">{{ telemetry?.TrackTemp?.toFixed(0) ?? '--' }}°C</div>
          </div>
          <div class="flex-1 bg-r-surface rounded-lg p-2">
            <div class="r-label text-[10px]">Air</div>
            <div class="font-mono text-xs text-r-text">{{ telemetry?.AirTemp?.toFixed(0) ?? '--' }}°C</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Row 4: Tires -->
    <TiresTireWidget v-if="telemetry" :tires="telemetry" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'default' })

const { telemetry, connected, speedKmh, gearLabel, fuelLaps, position, fmtTime } = useIRacing()
</script>
