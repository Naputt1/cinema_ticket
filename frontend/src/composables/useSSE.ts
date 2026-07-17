import { onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '../stores/notification'

const API = import.meta.env.VITE_API_URL || '/api/v1'

export function useSSE() {
  const notification = useNotificationStore()

  let seatEventSource: EventSource | null = null
  let notifEventSource: EventSource | null = null

  function connect() {
    if (notifEventSource) notifEventSource.close()

    try {
      const es = new EventSource(`${API}/events/notifications`)
      es.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          notification.addNotification(data.event, data.data)
        } catch {}
      }
      notifEventSource = es
    } catch {}
  }

  function connectSeatUpdates(showtimeID: string) {
    if (!showtimeID) return

    if (seatEventSource) seatEventSource.close()

    try {
      const es = new EventSource(`${API}/events/seats?showtime_id=${showtimeID}`)
      seatEventSource = es
    } catch {}
  }

  function disconnectSeatUpdates() {
    if (seatEventSource) {
      seatEventSource.close()
      seatEventSource = null
    }
  }

  function disconnect() {
    disconnectSeatUpdates()
    if (notifEventSource) {
      notifEventSource.close()
      notifEventSource = null
    }
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return { connectSeatUpdates, disconnectSeatUpdates }
}
