<template>
  <div class="r-card p-4">
    <div class="r-label mb-3">Tire Temperatures</div>

    <!-- 2×2 grid -->
    <div class="grid grid-cols-2 gap-3">
      <TireSingle label="FL" :cl="tires.LFtempCL" :cm="tires.LFtempCM" :cr="tires.LFtempCR" />
      <TireSingle label="FR" :cl="tires.RFtempCL" :cm="tires.RFtempCM" :cr="tires.RFtempCR" />
      <TireSingle label="RL" :cl="tires.LRtempCL" :cm="tires.LRtempCM" :cr="tires.LRtempCR" />
      <TireSingle label="RR" :cl="tires.RRtempCL" :cm="tires.RRtempCM" :cr="tires.RRtempCR" />
    </div>

    <!-- Temperature legend -->
    <div class="mt-3 flex items-center justify-between text-[10px] font-mono text-r-muted">
      <span style="color: #4cc9f0">COLD</span>
      <div class="flex-1 mx-2 h-1 rounded" style="background: linear-gradient(to right, #4cc9f0, #06d6a0, #ffd166, #e63946)" />
      <span style="color: #e63946">HOT</span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface TireData {
  LFtempCL: number; LFtempCM: number; LFtempCR: number
  RFtempCL: number; RFtempCM: number; RFtempCR: number
  LRtempCL: number; LRtempCM: number; LRtempCR: number
  RRtempCL: number; RRtempCM: number; RRtempCR: number
}

withDefaults(defineProps<{ tires: TireData }>(), {
  tires: () => ({
    LFtempCL: 60, LFtempCM: 60, LFtempCR: 60,
    RFtempCL: 60, RFtempCM: 60, RFtempCR: 60,
    LRtempCL: 60, LRtempCM: 60, LRtempCR: 60,
    RRtempCL: 60, RRtempCM: 60, RRtempCR: 60,
  }),
})
</script>
