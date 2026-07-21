<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'

const route = useRoute()
const router = useRouter()
const API = import.meta.env.VITE_API_URL || '/api/v1'

const {
  data: booking,
  isLoading,
  isError,
} = useQuery({
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

const queryClient = useQueryClient()

const cancelMutation = useMutation({
  mutationFn: async () => {
    const res = await fetch(`${API}/bookings/${route.params.id}/refund`, {
      method: 'POST',
      credentials: 'include',
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error || 'Cancel failed')
    }
  },
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['booking', route.params.id] })
    toast.success('Booking cancelled')
  },
  onError: (err: Error) => {
    toast.error(err.message)
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

function imgError(e: Event) {
  const el = e.target as HTMLImageElement
  el.style.display = 'none'
}
</script>

<template>
  <div>
    <Button variant="link" class="p-0 text-muted-foreground mb-4" @click="router.back()">&larr; Back</Button>

    <div v-if="isLoading" class="space-y-4">
      <div class="flex gap-6">
        <Skeleton class="w-48 h-72 shrink-0 rounded-lg" />
        <div class="flex-1 space-y-3">
          <Skeleton class="h-8 w-1/2" />
          <Skeleton class="h-5 w-1/3" />
          <Skeleton class="h-5 w-1/4" />
          <Skeleton class="h-5 w-1/4" />
          <Skeleton class="h-5 w-1/4" />
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Booking not found</p>
      <p class="text-sm text-muted-foreground mb-4">This booking may have been removed or the link is invalid.</p>
      <router-link to="/bookings">
        <Button>My Bookings</Button>
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

      <Card>
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Status</span>
          <Badge :variant="booking.status === 'BOOKED' ? 'default' : booking.status === 'LOCKED' ? 'secondary' : 'destructive'">{{
            booking.status
          }}</Badge>
          <Button
            v-if="booking.status === 'LOCKED' || booking.status === 'BOOKED'"
            variant="destructive"
            size="sm"
            @click="cancelMutation.mutate()"
            :disabled="cancelMutation.isPending.value"
          >
            {{ cancelMutation.isPending.value ? 'Cancelling...' : 'Cancel Booking' }}
          </Button>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Screen</span>
          <span class="text-sm font-medium">{{ booking.showtime?.screen }}</span>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Date</span>
          <span class="text-sm font-medium">{{ formatDate(booking.showtime?.start_time) }}</span>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Time</span>
          <span class="text-sm font-medium">{{ formatTime(booking.showtime?.start_time) }}</span>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Seats ({{ booking.seats?.length }} × ${{ booking.price?.toFixed(2) }})</span>
          <div class="text-right">
            <p v-for="s in booking.seats" :key="s.id" class="text-sm font-medium">{{ seatLabel(s) }}</p>
          </div>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between bg-muted/50">
          <span class="text-sm font-semibold">Total</span>
          <span class="text-lg font-bold">${{ booking.total_amount?.toFixed(2) }}</span>
        </div>
        <Separator />
        <div class="p-4 flex items-center justify-between">
          <span class="text-sm text-muted-foreground">Booked on</span>
          <span class="text-sm text-muted-foreground">{{ new Date(booking.created_at).toLocaleString() }}</span>
        </div>
      </Card>
    </div>
  </div>
</template>
