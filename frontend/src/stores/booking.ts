import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface SeatStatus {
  id: string
  showtime_id: string
  row: string
  number: number
  status: 'AVAILABLE' | 'LOCKED' | 'BOOKED'
  screen: string
}

export const useBookingStore = defineStore('booking', () => {
  const selectedSeats = ref<Set<string>>(new Set())
  const currentBookingId = ref<string | null>(null)
  const lockExpiresAt = ref<string | null>(null)

  function toggleSeat(seatId: string) {
    const next = new Set(selectedSeats.value)
    if (next.has(seatId)) {
      next.delete(seatId)
    } else {
      next.add(seatId)
    }
    selectedSeats.value = next
  }

  function clearSelection() {
    selectedSeats.value = new Set()
    currentBookingId.value = null
    lockExpiresAt.value = null
  }

  function setBooking(bookingId: string, expiresAt: string) {
    currentBookingId.value = bookingId
    lockExpiresAt.value = expiresAt
  }

  return {
    selectedSeats,
    currentBookingId,
    lockExpiresAt,
    toggleSeat,
    clearSelection,
    setBooking,
  }
})
