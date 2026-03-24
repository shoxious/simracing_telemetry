<template>
  <nav class="safe-bottom fixed bottom-0 left-0 right-0 z-50 bg-r-surface/95 backdrop-blur border-t border-r-border">
    <div class="flex items-stretch h-16">
      <NuxtLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="flex-1 flex flex-col items-center justify-center gap-0.5 text-r-muted transition-colors duration-200 active:opacity-60"
        active-class="!text-r-accent"
      >
        <component :is="item.icon" class="w-5 h-5" />
        <span class="text-[10px] font-mono tracking-wide">{{ item.label }}</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { defineComponent, h } from 'vue'

function svgIcon(paths: string) {
  return defineComponent({
    render: () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' },
      [h('path', { d: paths })]
    ),
  })
}

const navItems = [
  {
    to: '/',
    label: 'DASH',
    icon: svgIcon('M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6'),
  },
  {
    to: '/telemetry',
    label: 'TELEM',
    icon: svgIcon('M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z'),
  },
  {
    to: '/timing',
    label: 'TIMING',
    icon: svgIcon('M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z'),
  },
  {
    to: '/strategy',
    label: 'STRAT',
    icon: svgIcon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01'),
  },
]
</script>
