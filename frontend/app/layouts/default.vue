<template>
  <div class="flex flex-col min-h-dvh bg-r-bg">
    <!-- Top status bar -->
    <AppStatusBar />

    <!-- Main content area with bottom padding for nav -->
    <main class="flex-1 overflow-y-auto pb-20 md:pb-4">
      <slot />
    </main>

    <!-- Bottom navigation (mobile/tablet) -->
    <AppNavBar class="md:hidden" />
  </div>
</template>

<script setup lang="ts">
// Initialise the WebSocket connection for the whole app lifetime
useWebSocket()
</script>
