<template>
  <div class="flex flex-col min-h-dvh bg-r-bg">
    <!-- Top status bar -->
    <AppStatusBar />

    <!-- Main content area with bottom padding for nav -->
    <main class="flex-1 overflow-y-auto pb-20 md:pb-4 relative">
      <!-- Offline overlay – shown when not connected to the backend -->
      <Transition name="fade">
        <div
          v-if="!connected"
          class="absolute inset-0 z-40 flex flex-col items-center justify-center gap-6 bg-r-bg px-6 text-center"
        >
          <!-- Animated logo -->
          <div class="relative w-20 h-20">
            <div class="absolute inset-0 rounded-full border-2 border-r-accent/20 animate-ping" />
            <div class="absolute inset-2 rounded-full border-2 border-r-accent/40 animate-ping [animation-delay:300ms]" />
            <div class="relative w-20 h-20 rounded-full bg-r-surface border border-r-border flex items-center justify-center">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="w-9 h-9 text-r-accent">
                <circle cx="12" cy="12" r="9" />
                <circle cx="12" cy="12" r="3" />
                <path d="M12 3v2M12 19v2M3 12h2M19 12h2" stroke-linecap="round" />
              </svg>
            </div>
          </div>

          <!-- Status text -->
          <div class="space-y-1">
            <div class="font-mono font-bold text-lg text-r-text tracking-tight">
              Waiting for iRacing
            </div>
            <div class="text-sm text-r-muted font-mono">
              {{ wsStatusText }}
            </div>
          </div>

          <!-- Hint cards -->
          <div class="w-full max-w-sm space-y-2 text-left">
            <div class="r-card p-3 flex items-start gap-3">
              <div class="w-6 h-6 rounded bg-r-blue/20 flex-shrink-0 flex items-center justify-center mt-0.5">
                <span class="text-r-blue text-xs font-bold">1</span>
              </div>
              <div>
                <div class="text-xs font-mono text-r-text font-medium">Start iRacing</div>
                <div class="text-[11px] text-r-muted font-mono mt-0.5">Launch iRacing and load into a session</div>
              </div>
            </div>
            <div class="r-card p-3 flex items-start gap-3">
              <div class="w-6 h-6 rounded bg-r-blue/20 flex-shrink-0 flex items-center justify-center mt-0.5">
                <span class="text-r-blue text-xs font-bold">2</span>
              </div>
              <div>
                <div class="text-xs font-mono text-r-text font-medium">Drive onto the track</div>
                <div class="text-[11px] text-r-muted font-mono mt-0.5">Dashboard activates automatically</div>
              </div>
            </div>
            <div v-if="!wsConnected" class="r-card p-3 flex items-start gap-3 border border-r-accent/30">
              <div class="w-6 h-6 rounded bg-r-accent/20 flex-shrink-0 flex items-center justify-center mt-0.5">
                <svg viewBox="0 0 16 16" fill="currentColor" class="w-3 h-3 text-r-accent">
                  <path d="M8 1a7 7 0 100 14A7 7 0 008 1zm.75 4a.75.75 0 00-1.5 0v4a.75.75 0 001.5 0V5zm-.75 7a1 1 0 110-2 1 1 0 010 2z"/>
                </svg>
              </div>
              <div>
                <div class="text-xs font-mono text-r-accent font-medium">No connection to backend</div>
                <div class="text-[11px] text-r-muted font-mono mt-0.5">Make sure simracing-dashboard.exe is running</div>
              </div>
            </div>
          </div>

          <!-- Reconnect dots -->
          <div class="flex gap-1.5 items-center">
            <div class="w-1.5 h-1.5 rounded-full bg-r-accent animate-bounce [animation-delay:0ms]" />
            <div class="w-1.5 h-1.5 rounded-full bg-r-accent animate-bounce [animation-delay:150ms]" />
            <div class="w-1.5 h-1.5 rounded-full bg-r-accent animate-bounce [animation-delay:300ms]" />
          </div>
        </div>
      </Transition>

      <slot />
    </main>

    <!-- Bottom navigation (mobile/tablet) -->
    <AppNavBar class="md:hidden" />
  </div>
</template>

<script setup lang="ts">
const { connected } = useIRacing()
const { wsConnected, wsStatusText } = useWebSocket()
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
