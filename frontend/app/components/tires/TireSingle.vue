<template>
  <div class="flex flex-col items-center gap-1.5">
    <!-- Tire label -->
    <span class="r-label text-[10px]">{{ label }}</span>

    <!-- Three-column tire visualization -->
    <div class="flex gap-0.5 items-end" style="height: 52px">
      <div
        v-for="(temp, i) in [cl, cm, cr]"
        :key="i"
        class="w-4 rounded-sm relative flex items-end justify-center"
        :style="{ height: '100%', backgroundColor: '#15151f' }"
      >
        <div
          class="w-full rounded-sm transition-all duration-500"
          :style="{ height: '100%', backgroundColor: tempToColor(temp) }"
        />
      </div>
    </div>

    <!-- Average temp label -->
    <span class="font-mono text-xs font-semibold" :style="{ color: tempToColor(avgTemp) }">
      {{ Math.round(avgTemp) }}°
    </span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  label: string
  cl: number  // left strip
  cm: number  // mid strip
  cr: number  // right strip
}>()

const avgTemp = computed(() => (props.cl + props.cm + props.cr) / 3)

function lerp(a: number, b: number, t: number): number {
  return a + t * (b - a)
}

function hexToRgb(hex: string) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return [r, g, b]
}

function interpolateColors(c1: string, c2: string, t: number): string {
  const [r1, g1, b1] = hexToRgb(c1)
  const [r2, g2, b2] = hexToRgb(c2)
  const r = Math.round(lerp(r1, r2, t))
  const g = Math.round(lerp(g1, g2, t))
  const b = Math.round(lerp(b1, b2, t))
  return `rgb(${r},${g},${b})`
}

// Temp ranges: 40°C (cold) → 80°C (warm) → 100°C (optimal) → 115°C+ (hot)
function tempToColor(temp: number): string {
  if (temp <= 40) return '#4cc9f0'
  if (temp <= 70) return interpolateColors('#4cc9f0', '#06d6a0', (temp - 40) / 30)
  if (temp <= 95) return interpolateColors('#06d6a0', '#ffd166', (temp - 70) / 25)
  if (temp <= 115) return interpolateColors('#ffd166', '#e63946', (temp - 95) / 20)
  return '#e63946'
}
</script>
