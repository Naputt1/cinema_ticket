<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation } from '@tanstack/vue-query'
import { useAuthStore } from '../stores/auth'
import { useBookingStore } from '../stores/booking'
import { useSSE } from '../composables/useSSE'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const booking = useBookingStore()
const sse = useSSE()
const API = import.meta.env.VITE_API_URL || '/api/v1'

const countdown = ref<number>(0)
let timer: ReturnType<typeof setInterval> | null = null

const {
  data: seats,
  isLoading,
  isError,
  refetch: refetchSeats,
} = useQuery({
  queryKey: ['seats', route.params.id],
  queryFn: async () => {
    const res = await fetch(`${API}/showtimes/${route.params.id}/seats`, {
      credentials: 'include',
    })
    if (!res.ok) throw new Error('Failed to fetch seats')
    return res.json()
  },
})

const groupedSeats = computed(() => {
  const rows: Record<string, any[]> = {}
  for (const s of (seats.value as any[]) || []) {
    if (!rows[s.row]) rows[s.row] = []
    rows[s.row].push(s)
  }
  for (const row of Object.keys(rows)) {
    rows[row].sort((a, b) => a.number - b.number)
  }
  return rows
})

const selectedSeatLabels = computed(() => {
  const allSeats = (seats.value as any[]) || []
  return Array.from(booking.selectedSeats).map((id) => {
    const s = allSeats.find((s) => s.id === id)
    return s ? `${s.row}${s.number}` : id
  })
})

const lockedMutation = useMutation({
  mutationFn: async (seatIDs: string[]) => {
    const res = await fetch(`${API}/bookings/lock`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ showtime_id: route.params.id, seat_ids: seatIDs }),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'Failed to lock seats')
    }
    return res.json()
  },
  onSuccess: (data: any) => {
    booking.setBooking(data.id, data.expires_at)
    startCountdown(data.expires_at)
    refetchSeats()
  },
})

const payMutation = useMutation({
  mutationFn: async () => {
    const res = await fetch(`${API}/bookings/pay`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ booking_id: booking.currentBookingId }),
    })
    if (!res.ok) throw new Error('Payment failed')
    return res.json()
  },
  onSuccess: () => {
    booking.clearSelection()
    if (timer) clearInterval(timer)
    countdown.value = 0
    refetchSeats()
    router.push('/bookings')
  },
})

const cancelMutation = useMutation({
  mutationFn: async () => {
    const res = await fetch(`${API}/bookings/cancel`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ booking_id: booking.currentBookingId }),
    })
    if (!res.ok) throw new Error('Cancel failed')
  },
  onSuccess: () => {
    booking.clearSelection()
    if (timer) clearInterval(timer)
    countdown.value = 0
    refetchSeats()
  },
})

const extendMutation = useMutation({
  mutationFn: async () => {
    const res = await fetch(`${API}/bookings/${booking.currentBookingId}/extend`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!res.ok) throw new Error('Extend failed')
  },
  onSuccess: () => {
    startCountdown(new Date(Date.now() + 5 * 60 * 1000).toISOString())
  },
})

function getSeatStatus(seat: any): 'AVAILABLE' | 'LOCKED' | 'BOOKED' | 'SELECTED' {
  if (booking.selectedSeats.has(seat.id)) return 'SELECTED'
  return seat.status
}

function toggleSeat(seat: any) {
  if (seat.status !== 'AVAILABLE' && !booking.selectedSeats.has(seat.id)) return
  if (booking.currentBookingId) return
  booking.toggleSeat(seat.id)
}

async function handleLock() {
  const ids = Array.from(booking.selectedSeats)
  if (ids.length === 0) return
  await lockedMutation.mutateAsync(ids)
}

function startCountdown(expiresAt: string) {
  const target = new Date(expiresAt).getTime()
  if (timer) clearInterval(timer)
  timer = setInterval(() => {
    const remaining = Math.max(0, Math.floor((target - Date.now()) / 1000))
    countdown.value = remaining
    if (remaining <= 0) {
      if (timer) clearInterval(timer)
      booking.clearSelection()
      refetchSeats()
    }
  }, 1000)
}

function seatColor(status: string) {
  switch (status) {
    case 'AVAILABLE':
      return 'bg-green-100 border-green-300 hover:bg-green-200 cursor-pointer'
    case 'LOCKED':
      return 'bg-yellow-100 border-yellow-300 cursor-not-allowed opacity-60'
    case 'BOOKED':
      return 'bg-red-100 border-red-300 cursor-not-allowed opacity-60'
    case 'SELECTED':
      return 'bg-blue-500 border-blue-600 text-white cursor-pointer'
    default:
      return 'bg-gray-100 border-gray-300'
  }
}

onMounted(() => {
  sse.connectSeatUpdates(route.params.id as string, () => refetchSeats())
})

onUnmounted(() => {
  sse.disconnectSeatUpdates()
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">Select Seats</h1>

    <div v-if="countdown > 0" class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded text-center">
      <p class="text-sm font-medium">Time left to pay: {{ Math.floor(countdown / 60) }}:{{ String(countdown % 60).padStart(2, '0') }}</p>
    </div>

    <div v-if="isLoading" class="border rounded-lg p-6 mb-6">
      <div class="text-center mb-6">
        <Skeleton class="w-3/4 h-2 mx-auto mb-2" />
        <p class="text-xs text-muted-foreground">SCREEN</p>
      </div>
      <div class="flex flex-col justify-center gap-6">
        <div v-for="r in 8" :key="r" class="flex flex-col items-center">
          <div class="flex flex-row gap-2 items-center">
            <span class="text-xs font-bold text-gray-300">{{ String.fromCharCode(64 + r) }}</span>
            <div class="flex gap-1">
              <div v-for="c in 8" :key="c">
                <Skeleton class="w-8 h-8" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="border rounded-lg p-6 mb-6 text-center">
      <p class="text-destructive mb-4">Failed to load seats</p>
      <Button @click="refetchSeats()">Retry</Button>
    </div>

    <div v-else-if="!seats || seats.length === 0" class="border rounded-lg p-6 mb-6 text-center">
      <p class="text-muted-foreground">No seats available for this showtime</p>
    </div>

    <div v-else class="border rounded-lg p-6 mb-6">
      <div class="text-center mb-6">
        <div class="w-3/4 h-2 bg-gray-300 rounded mx-auto mb-2"></div>
        <p class="text-xs text-muted-foreground">SCREEN</p>
      </div>

      <div class="flex flex-col justify-center gap-6">
        <div v-for="(rowSeats, row) in groupedSeats" :key="row" class="flex flex-col items-center">
          <div class="flex flex-row gap-2 items-center">
            <span class="text-xs font-bold">{{ row }}</span>
            <div class="flex gap-1">
              <button
                v-for="seat in rowSeats"
                :key="seat.id"
                :disabled="seat.status !== 'AVAILABLE' && !booking.selectedSeats.has(seat.id)"
                @click="toggleSeat(seat)"
                class="w-8 h-8 text-xs border rounded transition-colors"
                :class="seatColor(getSeatStatus(seat))"
              >
                {{ seat.number }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-center gap-6 mt-6 text-xs">
        <div class="flex items-center gap-1">
          <div class="w-4 h-4 bg-green-100 border border-green-300 rounded"></div>
          Available
        </div>
        <div class="flex items-center gap-1">
          <div class="w-4 h-4 bg-yellow-100 border border-yellow-300 rounded"></div>
          Locked
        </div>
        <div class="flex items-center gap-1">
          <div class="w-4 h-4 bg-red-100 border border-red-300 rounded"></div>
          Booked
        </div>
        <div class="flex items-center gap-1">
          <div class="w-4 h-4 bg-blue-500 border border-blue-600 rounded"></div>
          Selected
        </div>
      </div>
    </div>

    <div class="flex justify-between items-center">
      <p class="text-sm text-muted-foreground">Selected: {{ selectedSeatLabels.join(', ') || 'None' }}</p>
      <div class="flex gap-2">
        <button
          v-if="booking.currentBookingId"
          @click="extendMutation.mutate()"
          class="px-4 py-2 border rounded text-sm hover:bg-muted"
          :disabled="extendMutation.isPending.value"
        >
          {{ extendMutation.isPending.value ? 'Extending...' : 'Extend 5 min' }}
        </button>
        <button
          v-if="booking.currentBookingId"
          @click="cancelMutation.mutate()"
          class="px-4 py-2 border rounded text-sm hover:bg-muted"
          :disabled="cancelMutation.isPending.value"
        >
          Cancel
        </button>
        <button
          v-if="booking.currentBookingId"
          @click="payMutation.mutate()"
          class="px-4 py-2 bg-green-600 text-white rounded text-sm hover:bg-green-700"
          :disabled="payMutation.isPending.value"
        >
          Confirm Payment (Mock)
        </button>
        <button
          v-else-if="auth.isAuthenticated"
          @click="handleLock()"
          class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90"
          :disabled="booking.selectedSeats.size === 0 || lockedMutation.isPending.value"
        >
          {{ lockedMutation.isPending.value ? 'Locking...' : 'Lock Seats' }}
        </button>
        <router-link
          v-else
          :to="`/login?redirect=/showtimes/${route.params.id}`"
          class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90 inline-block"
        >
          Login to Book
        </router-link>
      </div>
    </div>

    <p v-if="lockedMutation.isError.value" class="mt-2 text-sm text-destructive">
      {{ (lockedMutation.error as any)?.message }}
    </p>
  </div>
</template>
