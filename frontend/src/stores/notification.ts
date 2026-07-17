import { defineStore } from 'pinia'
import { ref } from 'vue'
import { toast } from 'vue-sonner'

export interface Notification {
  id: string
  type: string
  message: string
  timestamp: Date
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])

  function addNotification(type: string, data: any) {
    const notif: Notification = {
      id: crypto.randomUUID(),
      type,
      message: formatMessage(type, data),
      timestamp: new Date(),
    }
    notifications.value.unshift(notif)
    toast(notif.message)
  }

  function formatMessage(type: string, data: any): string {
    switch (type) {
      case 'booking.success':
        return `Booking confirmed! Seats: ${data.seat_ids?.join(', ') || ''}`
      case 'booking.timeout':
        return 'Booking expired - seats have been released'
      case 'seat.locked':
        return `Seats ${data.seat_ids?.join(', ') || ''} are now locked`
      case 'seat.booked':
        return `Seats ${data.seat_ids?.join(', ') || ''} are now booked`
      case 'seat.released':
        return `Seats ${data.seat_ids?.join(', ') || ''} have been released`
      default:
        return `Event: ${type}`
    }
  }

  return { notifications, addNotification }
})
