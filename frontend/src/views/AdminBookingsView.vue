<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const bookingPage = ref(0)
const bookingLimit = ref(20)

const { data: bookings, isLoading, isError, refetch } = useQuery({
  queryKey: ['admin-bookings', bookingPage, bookingLimit],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/bookings?skip=${bookingPage.value * bookingLimit.value}&limit=${bookingLimit.value}`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch bookings')
    return res.json()
  },
})

async function forceCancelBooking(id: string) {
  if (!confirm('Cancel this booking?')) return
  await fetch(`${API}/admin/bookings/${id}/cancel`, { method: 'POST', credentials: 'include' })
  refetch()
}

function seatLabel(s: { row: string; number: number }) {
  return `${s.row}${s.number}`
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Bookings</h1>

    <div v-if="isLoading" class="space-y-2">
      <div v-for="n in 4" :key="n" class="border rounded p-3 text-sm animate-pulse">
        <div class="h-4 bg-muted rounded w-3/5 mb-2" />
        <div class="h-3 bg-muted rounded w-2/5 mb-1" />
        <div class="h-3 bg-muted rounded w-1/3 mb-1" />
        <div class="h-3 bg-muted rounded w-1/2" />
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-3">Failed to load bookings</p>
      <button @click="refetch()" class="px-4 py-2 rounded text-sm border">Retry</button>
    </div>

    <div v-else-if="!(bookings as any[])?.length" class="text-center py-8 text-muted-foreground">
      No bookings found
    </div>

    <div v-else class="space-y-2">
      <div v-for="b in (bookings as any[])" :key="b.id" class="border rounded p-3 text-sm">
        <div class="flex justify-between items-center mb-1">
          <span class="font-medium">{{ b.movie?.title }}</span>
          <span class="flex items-center gap-2">
            <span
              class="text-xs font-medium px-2 py-0.5 rounded"
              :class="b.status === 'BOOKED' ? 'text-green-700 bg-green-50' : b.status === 'LOCKED' ? 'text-yellow-700 bg-yellow-50' : 'text-red-700 bg-red-50'"
            >
              {{ b.status }}
            </span>
            <button
              v-if="b.status === 'LOCKED' || b.status === 'BOOKED'"
              @click="forceCancelBooking(b.id)"
              class="text-red-500 underline text-xs ml-2"
            >
              Cancel
            </button>
          </span>
        </div>
        <p class="text-muted-foreground">{{ b.showtime?.screen }} · {{ formatDateTime(b.showtime?.start_time) }}</p>
        <p class="text-muted-foreground">User: {{ b.user_email || b.user_id?.slice(0, 8) }}</p>
        <p class="text-muted-foreground">
          Seats: {{ b.seats?.map(seatLabel).join(', ') || '—' }}
          <span class="ml-2">${{ b.showtime?.price?.toFixed(2) }} × {{ b.seats?.length || 0 }}</span>
        </p>
        <p class="text-muted-foreground text-xs mt-1">{{ formatDateTime(b.created_at) }}</p>
      </div>

      <div class="flex justify-center gap-4 mt-4">
        <button
          @click="bookingPage = Math.max(0, bookingPage - 1)"
          :disabled="bookingPage === 0"
          class="px-3 py-1 border rounded text-sm disabled:opacity-50"
        >
          Previous
        </button>
        <span class="text-sm text-muted-foreground self-center">Page {{ bookingPage + 1 }}</span>
        <button
          @click="bookingPage++"
          :disabled="(bookings as any[])?.length < bookingLimit"
          class="px-3 py-1 border rounded text-sm disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>
