import { useIRacingStore } from '~/stores/iracing'

let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let fpsTimer: ReturnType<typeof setInterval> | null = null
let reconnectAttempts = 0

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

    try {
      ws = new WebSocket(getWsUrl())
    } catch {
      scheduleReconnect()
      return
    }

    ws.onopen = () => {
      reconnectAttempts = 0
      store.setStatus(true, store.simulate)
    }

    ws.onmessage = (event: MessageEvent) => {
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
      store.setStatus(false, store.simulate)
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

  return { connect, disconnect }
}
