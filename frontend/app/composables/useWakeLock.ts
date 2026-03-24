/**
 * useWakeLock – prevents the screen from turning off while the app is open.
 * Uses the Screen Wake Lock API (Safari 16.4+, Chrome 84+, Firefox 126+).
 * Automatically re-acquires the lock when the tab becomes visible again
 * (e.g. after the user switches apps and returns).
 */

let wakeLock: WakeLockSentinel | null = null

export function useWakeLock() {
  const supported = ref(false)
  const active    = ref(false)

  async function acquire() {
    if (!('wakeLock' in navigator)) return
    try {
      wakeLock        = await navigator.wakeLock.request('screen')
      active.value    = true
      wakeLock.addEventListener('release', () => {
        active.value = false
      })
    } catch {
      // Permission denied or not supported – fail silently
      active.value = false
    }
  }

  async function release() {
    if (!wakeLock) return
    await wakeLock.release()
    wakeLock     = null
    active.value = false
  }

  // Re-acquire when the page becomes visible again (tab switch / app switch)
  function onVisibilityChange() {
    if (document.visibilityState === 'visible') acquire()
  }

  onMounted(async () => {
    supported.value = 'wakeLock' in navigator
    await acquire()
    document.addEventListener('visibilitychange', onVisibilityChange)
  })

  onUnmounted(() => {
    release()
    document.removeEventListener('visibilitychange', onVisibilityChange)
  })

  return { supported, active }
}
