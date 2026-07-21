<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const {
  data: bookings,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['my-bookings'],
  queryFn: async () => {
    const res = await fetch(`${API}/bookings/my`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch bookings')
    return res.json()
  },
})

function seatLabel(s: { row: string; number: number }) {
  return `${s.row}${s.number}`
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function statusColor(status: string) {
  switch (status) {
    case 'BOOKED':
      return 'text-green-600 border-green-200 bg-green-50'
    case 'LOCKED':
      return 'text-yellow-600 border-yellow-200 bg-yellow-50'
    case 'CANCELLED':
    case 'EXPIRED':
      return 'text-red-600 border-red-200 bg-red-50'
    default:
      return 'text-muted-foreground border-border'
  }
}

function imgError(e: Event) {
  const el = e.target as HTMLImageElement
  el.style.display = 'none'
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">My Bookings</h1>

    <div v-if="isLoading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="border rounded-lg p-4 animate-pulse space-y-3">
        <div class="flex gap-4">
          <div class="w-20 h-28 bg-muted rounded shrink-0" />
          <div class="flex-1 space-y-2">
            <div class="h-5 bg-muted rounded w-1/2" />
            <div class="h-4 bg-muted rounded w-1/3" />
            <div class="h-4 bg-muted rounded w-1/4" />
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Failed to load bookings</p>
      <p class="text-sm text-muted-foreground mb-4">Check that the backend is running.</p>
      <button @click="refetch()" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90">Retry</button>
    </div>

    <div v-else-if="!bookings?.length" class="text-center py-16">
      <p class="text-lg text-muted-foreground">No bookings yet</p>
      <p class="text-sm text-muted-foreground mt-1">Browse movies and book your first ticket.</p>
    </div>

    <div v-else class="space-y-4">
      <router-link
        v-for="b in bookings as any[]"
        :key="b.id"
        :to="`/bookings/${b.id}`"
        class="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow bg-card flex flex-col sm:flex-row"
      >
        <img :src="b.movie?.poster_url" :alt="b.movie?.title" class="w-full sm:w-24 h-32 sm:h-28 object-cover shrink-0" @error="imgError" />
        <div class="p-4 flex-1 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
          <div class="space-y-1">
            <h2 class="font-bold">{{ b.movie?.title }}</h2>
            <p class="text-sm text-muted-foreground">{{ b.showtime?.screen }} · {{ formatDateTime(b.showtime?.start_time) }}</p>
            <p class="text-sm text-muted-foreground">Seats: {{ b.seats?.map(seatLabel).join(', ') }}</p>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-sm font-semibold">${{ b.total_amount?.toFixed(2) }}</span>
            <span class="text-xs font-medium px-2.5 py-1 rounded-full border" :class="statusColor(b.status)">
              {{ b.status }}
            </span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>
