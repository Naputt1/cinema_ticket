<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'

const route = useRoute()
const router = useRouter()
const API = import.meta.env.VITE_API_URL || '/api/v1'

const { data: booking, isLoading, isError } = useQuery({
  queryKey: ['booking', route.params.id],
  queryFn: async () => {
    const res = await fetch(`${API}/bookings/${route.params.id}`, { credentials: 'include' })
    if (!res.ok) {
      if (res.status === 404) throw new Error('not-found')
      throw new Error('Failed to fetch booking')
    }
    return res.json()
  },
})

function seatLabel(s: { row: string; number: number }) {
  return `Row ${s.row}, Seat ${s.number}`
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString([], { weekday: 'long', month: 'long', day: 'numeric', year: 'numeric' })
}

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function formatDuration(min: number) {
  const h = Math.floor(min / 60)
  const m = min % 60
  if (h === 0) return `${m}m`
  if (m === 0) return `${h}h`
  return `${h}h ${m}m`
}

function statusColor(status: string) {
  switch (status) {
    case 'BOOKED': return 'text-green-600 border-green-200 bg-green-50'
    case 'LOCKED': return 'text-yellow-600 border-yellow-200 bg-yellow-50'
    case 'CANCELLED': case 'EXPIRED': return 'text-red-600 border-red-200 bg-red-50'
    default: return 'text-muted-foreground border-border'
  }
}

function imgError(e: Event) {
  const el = e.target as HTMLImageElement
  el.style.display = 'none'
}
</script>

<template>
  <div>
    <button @click="router.back()" class="text-sm text-muted-foreground hover:underline mb-4 inline-block">&larr; Back</button>

    <div v-if="isLoading" class="animate-pulse space-y-4">
      <div class="flex gap-6">
        <div class="w-48 h-72 bg-muted rounded-lg shrink-0" />
        <div class="flex-1 space-y-3">
          <div class="h-8 bg-muted rounded w-1/2" />
          <div class="h-5 bg-muted rounded w-1/3" />
          <div class="h-5 bg-muted rounded w-1/4" />
          <div class="h-5 bg-muted rounded w-1/4" />
          <div class="h-5 bg-muted rounded w-1/4" />
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Booking not found</p>
      <p class="text-sm text-muted-foreground mb-4">This booking may have been removed or the link is invalid.</p>
      <router-link to="/bookings" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90 inline-block">
        My Bookings
      </router-link>
    </div>

    <div v-else-if="booking" class="max-w-2xl">
      <div class="flex flex-col sm:flex-row gap-6 mb-8">
        <img
          :src="booking.movie?.poster_url"
          :alt="booking.movie?.title"
          class="w-48 h-72 object-cover rounded-lg shrink-0"
          @error="imgError"
        />
        <div class="space-y-3">
          <h1 class="text-2xl font-bold">{{ booking.movie?.title }}</h1>
          <p class="text-sm text-muted-foreground">{{ booking.movie?.genre }} · {{ formatDuration(booking.movie?.duration) }}</p>
          <p class="text-sm text-muted-foreground">{{ booking.movie?.description }}</p>
        </div>
      </div>

      <div class="border rounded-lg divide-y">
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Status</span>
          <span
            class="text-sm font-semibold px-3 py-1 rounded-full border"
            :class="statusColor(booking.status)"
          >
            {{ booking.status }}
          </span>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Screen</span>
          <span class="text-sm font-medium">{{ booking.showtime?.screen }}</span>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Date</span>
          <span class="text-sm font-medium">{{ formatDate(booking.showtime?.start_time) }}</span>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Time</span>
          <span class="text-sm font-medium">{{ formatTime(booking.showtime?.start_time) }}</span>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Seats</span>
          <div class="text-right">
            <p
              v-for="s in booking.seats"
              :key="s.id"
              class="text-sm font-medium"
            >
              {{ seatLabel(s) }}
            </p>
          </div>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Price per seat</span>
          <span class="text-sm font-medium">${{ booking.showtime?.price?.toFixed(2) }}</span>
        </div>
        <div class="p-4 flex items-center justify-between bg-muted/50">
          <span class="text-sm font-semibold">Total</span>
          <span class="text-lg font-bold">
            ${{ (booking.showtime?.price * (booking.seats?.length || 1))?.toFixed(2) }}
          </span>
        </div>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Booked on</span>
          <span class="text-sm text-muted-foreground">{{ new Date(booking.created_at).toLocaleString() }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
