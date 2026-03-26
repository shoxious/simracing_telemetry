import { useIRacingStore } from '~/stores/iracing'

let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let fpsTimer: ReturnType<typeof setInterval> | null = null
let reconnectAttempts = 0

// Shared reactive state accessible from any composable call — exported for direct import
export const wsConnected = ref(false)
export const wsStatusText = ref('Connecting...')

function getWsUrl(): string {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // In Nuxt dev mode, connect to the Go backend on 8080; in production same host
  const host = import.meta.dev ? 'localhost:8080' : window.location.host
  return `${proto}//${host}/ws`
}

export function useWebSocket() {
  const store = useIRacingStore()

  function connect() {
    if (ws && ws.readyState !== WebSocket.CLOSED) return

    wsStatusText.value = 'Connecting to dashboard...'

    try {
      ws = new WebSocket(getWsUrl())
    } catch {
      wsConnected.value = false
      wsStatusText.value = 'Cannot reach backend'
      scheduleReconnect()
      return
    }

    ws.onopen = () => {
      reconnectAttempts = 0
      wsConnected.value = true
      wsStatusText.value = 'Connected – waiting for iRacing...'
      if (!store.demoMode) {
        store.setStatus(true, store.simulate)
      }
    }

    ws.onmessage = (event: MessageEvent) => {
      if (store.demoMode) return
      try {
        const msg = JSON.parse(event.data as string)
        switch (msg.type) {
          case 'telemetry':
            store.setTelemetry(msg.data)
            break
          case 'session':
            store.setSession(msg.yaml)
            break
          case 'status':
            store.setStatus(msg.connected, msg.simulate)
            break
        }
      } catch {
        // ignore malformed messages
      }
    }

    ws.onclose = () => {
      wsConnected.value = false
      wsStatusText.value = reconnectAttempts > 0
        ? `Reconnecting... (attempt ${reconnectAttempts})`
        : 'Connection lost – reconnecting...'
      if (!store.demoMode) {
        store.setStatus(false, store.simulate)
      }
      scheduleReconnect()
    }

    ws.onerror = () => {
      ws?.close()
    }
  }

  function scheduleReconnect() {
    if (reconnectTimer) return
    const delay = Math.min(8000, 1000 * Math.pow(1.5, reconnectAttempts))
    reconnectAttempts++
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      connect()
    }, delay)
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    ws?.close()
    ws = null
  }

  onMounted(() => {
    connect()
    // FPS counter reset every second
    fpsTimer = setInterval(() => store.tickFPS(), 1000)
  })

  onUnmounted(() => {
    disconnect()
    if (fpsTimer) clearInterval(fpsTimer)
  })

  return { connect, disconnect, wsConnected, wsStatusText }
}
