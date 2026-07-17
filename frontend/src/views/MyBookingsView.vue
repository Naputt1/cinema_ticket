<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const { data: bookings, isLoading } = useQuery({
  queryKey: ['my-bookings'],
  queryFn: async () => {
    const res = await fetch(`${API}/bookings/my`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch bookings')
    return res.json()
  },
})

function statusColor(status: string) {
  switch (status) {
    case 'BOOKED': return 'text-green-600'
    case 'LOCKED': return 'text-yellow-600'
    case 'CANCELLED': case 'EXPIRED': return 'text-red-600'
    default: return ''
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">My Bookings</h1>
    <div v-if="isLoading" class="text-center py-8 text-muted-foreground">Loading...</div>
    <div v-else-if="!bookings?.length" class="text-center py-8 text-muted-foreground">No bookings yet.</div>
    <div v-else class="space-y-4">
      <div v-for="b in (bookings as any[])" :key="b.id" class="border rounded-lg p-4">
        <div class="flex justify-between items-start">
          <div>
            <p class="font-medium">Booking #{{ b.id?.slice(0, 8) }}</p>
            <p class="text-sm text-muted-foreground">Seats: {{ b.seat_ids?.join(', ') }}</p>
            <p class="text-sm text-muted-foreground">{{ formatDate(b.created_at) }}</p>
          </div>
          <span class="font-bold text-sm" :class="statusColor(b.status)">{{ b.status }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
