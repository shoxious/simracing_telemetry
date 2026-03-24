<template>
  <div class="flex flex-col items-center gap-2">
    <!-- Position label -->
    <div class="text-[10px] font-mono font-bold tracking-widest" :style="{ color: avgColor }">
      {{ label }}
    </div>

    <!-- Tire body: three vertical segments -->
    <div class="relative flex gap-0.5" style="width: 54px; height: 70px">
      <!-- Outer tire border -->
      <div
        class="absolute inset-0 rounded-lg border-2 pointer-events-none z-10"
        :style="{ borderColor: avgColor + '60' }"
      />
      <div
        v-for="(seg, i) in segments"
        :key="i"
        class="flex-1 rounded-sm flex flex-col items-center justify-end pb-1.5 relative overflow-hidden"
        :style="{ backgroundColor: seg.bgColor }"
      >
        <!-- Glow overlay -->
        <div
          class="absolute inset-0 opacity-40"
          :style="{ background: `linear-gradient(to top, ${seg.color}80, transparent)` }"
        />
        <span
          class="relative z-10 font-mono font-bold leading-none"
          style="font-size: 9px"
          :style="{ color: seg.color }"
        >{{ Math.round(seg.temp) }}</span>
      </div>
    </div>

    <!-- Strip labels -->
    <div class="flex gap-0.5 text-[8px] font-mono text-r-dim" style="width: 54px">
      <div class="flex-1 text-center">I</div>
      <div class="flex-1 text-center">M</div>
      <div class="flex-1 text-center">O</div>
    </div>

    <!-- Average temp badge -->
    <div
      class="px-2 py-0.5 rounded-full font-mono font-bold text-xs border"
      :style="{ color: avgColor, borderColor: avgColor + '50', backgroundColor: avgColor + '15' }"
    >
      {{ Math.round(avgTemp) }}°C
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  label: string
  cl: number  // inner strip
  cm: number  // middle strip
  cr: number  // outer strip
}>()

function lerp(a: number, b: number, t: number) { return a + t * (b - a) }
function hexToRgb(hex: string) {
  return [parseInt(hex.slice(1,3),16), parseInt(hex.slice(3,5),16), parseInt(hex.slice(5,7),16)]
}
function blend(c1: string, c2: string, t: number) {
  const [r1,g1,b1] = hexToRgb(c1); const [r2,g2,b2] = hexToRgb(c2)
  return `rgb(${Math.round(lerp(r1,r2,t))},${Math.round(lerp(g1,g2,t))},${Math.round(lerp(b1,b2,t))})`
}

function tempToColor(t: number) {
  if (t <= 40)  return '#4cc9f0'
  if (t <= 70)  return blend('#4cc9f0','#06d6a0', (t-40)/30)
  if (t <= 95)  return blend('#06d6a0','#ffd166', (t-70)/25)
  if (t <= 115) return blend('#ffd166','#e63946', (t-95)/20)
  return '#e63946'
}

const avgTemp = computed(() => (props.cl + props.cm + props.cr) / 3)
const avgColor = computed(() => tempToColor(avgTemp.value))

const segments = computed(() => [props.cl, props.cm, props.cr].map(temp => ({
  temp,
  color: tempToColor(temp),
  bgColor: tempToColor(temp) + '22',
})))
</script>
