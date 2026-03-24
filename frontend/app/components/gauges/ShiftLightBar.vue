<template>
  <div class="r-card px-3 py-2">
    <div class="flex items-center gap-1">
      <div
        v-for="(led, i) in leds"
        :key="i"
        class="flex-1 h-3.5 rounded-full transition-colors duration-50"
        :style="{
          backgroundColor: led.lit ? led.color : '#0d0d1a',
          boxShadow: led.lit ? `0 0 8px ${led.color}, 0 0 2px ${led.color}` : 'none',
        }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  rpm?: number
  maxRpm?: number
}>(), { rpm: 0, maxRpm: 8600 })

const leds = computed(() => {
  const pct = props.rpm / props.maxRpm
  // LED activates starting at 65% RPM, 15 LEDs total
  return Array.from({ length: 15 }, (_, i) => {
    const threshold = 0.65 + (i / 14) * 0.35   // 65% → 100%
    let color: string
    if (i < 5)       color = '#06d6a0'  // green zone
    else if (i < 10) color = '#ffd166'  // yellow zone
    else             color = '#e63946'  // red zone
    return { lit: pct >= threshold, color }
  })
})
</script>
