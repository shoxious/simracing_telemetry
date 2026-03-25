<template>
  <div class="r-card p-4">
    <div class="r-label mb-3 flex items-center justify-between">
      <span>Tire Temperatures</span>
      <div class="hidden sm:flex items-center gap-2 text-[10px] font-mono">
        <span style="color:#4cc9f0">COLD</span>
        <div class="w-16 h-1 rounded" style="background: linear-gradient(to right,#4cc9f0,#06d6a0,#ffd166,#e63946)" />
        <span style="color:#e63946">HOT</span>
      </div>
    </div>

    <!-- Car footprint grid: [FL] [car] [FR] / [RL] [car] [RR] -->
    <div class="grid gap-2" style="grid-template-columns: 1fr 80px 1fr; grid-template-rows: auto auto">

      <!-- FL -->
      <div class="flex justify-end">
        <TireSingle label="FL" :cl="tires.LFtempCL" :cm="tires.LFtempCM" :cr="tires.LFtempCR" />
      </div>

      <!-- Car top half -->
      <div class="flex justify-center">
        <svg viewBox="0 0 80 160" width="80" height="80" class="overflow-visible">
          <!-- Car body front -->
          <rect x="18" y="5" width="44" height="62" rx="12" fill="#15151f" stroke="#3a3a55" stroke-width="1.5"/>
          <!-- Windscreen -->
          <rect x="24" y="12" width="32" height="22" rx="4" fill="#1e1e35" stroke="#4cc9f0" stroke-width="0.8" opacity="0.6"/>
          <!-- Front wing -->
          <rect x="8" y="5" width="10" height="20" rx="3" fill="#1a1a2e" stroke="#3a3a55" stroke-width="1"/>
          <rect x="62" y="5" width="10" height="20" rx="3" fill="#1a1a2e" stroke="#3a3a55" stroke-width="1"/>
        </svg>
      </div>

      <!-- FR -->
      <div class="flex justify-start">
        <TireSingle label="FR" :cl="tires.RFtempCL" :cm="tires.RFtempCM" :cr="tires.RFtempCR" />
      </div>

      <!-- RL -->
      <div class="flex justify-end">
        <TireSingle label="RL" :cl="tires.LRtempCL" :cm="tires.LRtempCM" :cr="tires.LRtempCR" />
      </div>

      <!-- Car bottom half -->
      <div class="flex justify-center">
        <svg viewBox="0 0 80 160" width="80" height="80" class="overflow-visible">
          <!-- Car body rear -->
          <rect x="18" y="78" width="44" height="62" rx="8" fill="#15151f" stroke="#3a3a55" stroke-width="1.5"/>
          <!-- Rear diffuser detail -->
          <rect x="22" y="120" width="36" height="12" rx="2" fill="#1a1a2e" stroke="#3a3a55" stroke-width="1"/>
          <!-- Exhaust -->
          <circle cx="30" cy="136" r="3" fill="#e63946" opacity="0.3"/>
          <circle cx="50" cy="136" r="3" fill="#e63946" opacity="0.3"/>
          <!-- Rear wing -->
          <rect x="10" y="130" width="60" height="8" rx="2" fill="#1a1a2e" stroke="#3a3a55" stroke-width="1.5"/>
        </svg>
      </div>

      <!-- RR -->
      <div class="flex justify-start">
        <TireSingle label="RR" :cl="tires.RRtempCL" :cm="tires.RRtempCM" :cr="tires.RRtempCR" />
      </div>
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
    LFtempCL:80,LFtempCM:85,LFtempCR:82,
    RFtempCL:81,RFtempCM:86,RFtempCR:83,
    LRtempCL:88,LRtempCM:92,LRtempCR:90,
    RRtempCL:87,RRtempCM:91,RRtempCR:89,
  }),
})
</script>
