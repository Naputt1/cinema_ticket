<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref<'bookings' | 'logs'>('bookings')

const { data: bookings, isLoading: bkLoading, isError: bkError, refetch: bkRefetch } = useQuery({
  queryKey: ['admin-bookings'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/bookings`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch bookings')
    return res.json()
  },
})

const { data: logs, isLoading: logsLoading, isError: logsError, refetch: logsRefetch } = useQuery({
  queryKey: ['admin-audit-logs'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/audit-logs`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch logs')
    return res.json()
  },
})

function seatLabel(s: { row: string; number: number }) {
  return `${s.row}${s.number}`
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}

function formatPrice(price: number) {
  return price?.toFixed(2)
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Admin Dashboard</h1>

    <div class="flex gap-2 mb-6">
      <button
        @click="activeTab = 'bookings'"
        class="px-4 py-2 rounded text-sm"
        :class="activeTab === 'bookings' ? 'bg-primary text-primary-foreground' : 'border'"
      >
        Bookings
      </button>
      <button
        @click="activeTab = 'logs'"
        class="px-4 py-2 rounded text-sm"
        :class="activeTab === 'logs' ? 'bg-primary text-primary-foreground' : 'border'"
      >
        Audit Logs
      </button>
    </div>

    <div v-if="activeTab === 'bookings'">
      <div v-if="bkLoading" class="space-y-2">
        <div v-for="n in 4" :key="n" class="border rounded p-3 text-sm animate-pulse">
          <div class="h-4 bg-muted rounded w-3/5 mb-2"></div>
          <div class="h-3 bg-muted rounded w-2/5 mb-1"></div>
          <div class="h-3 bg-muted rounded w-1/3 mb-1"></div>
          <div class="h-3 bg-muted rounded w-1/2"></div>
        </div>
      </div>
      <div v-else-if="bkError" class="text-center py-8">
        <p class="text-destructive mb-3">Failed to load bookings</p>
        <button @click="bkRefetch()" class="px-4 py-2 rounded text-sm border">Retry</button>
      </div>
      <div v-else-if="!(bookings as any[])?.length" class="text-center py-8 text-muted-foreground">
        No bookings found
      </div>
      <div v-else class="space-y-2">
        <div v-for="b in (bookings as any[])" :key="b.id" class="border rounded p-3 text-sm">
          <div class="flex justify-between items-center mb-1">
            <span class="font-medium">{{ b.movie?.title }}</span>
            <span
              class="text-xs font-medium px-2 py-0.5 rounded"
              :class="b.status === 'BOOKED' ? 'text-green-700 bg-green-50' : b.status === 'LOCKED' ? 'text-yellow-700 bg-yellow-50' : 'text-red-700 bg-red-50'"
            >
              {{ b.status }}
            </span>
          </div>
          <p class="text-muted-foreground">
            {{ b.showtime?.screen }} · {{ formatDateTime(b.showtime?.start_time) }}
          </p>
          <p class="text-muted-foreground">
            User: {{ b.user_email || b.user_id?.slice(0, 8) }}
          </p>
          <p class="text-muted-foreground">
            Seats: {{ b.seats?.map(seatLabel).join(', ') || '—' }}
            <span class="ml-2">${{ formatPrice(b.showtime?.price) }} × {{ b.seats?.length || 0 }}</span>
          </p>
          <p class="text-muted-foreground text-xs mt-1">{{ formatDateTime(b.created_at) }}</p>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'logs'">
      <div v-if="logsLoading" class="space-y-2">
        <div v-for="n in 3" :key="n" class="border rounded p-3 text-sm animate-pulse">
          <div class="h-4 bg-muted rounded w-2/5 mb-2"></div>
          <div class="h-3 bg-muted rounded w-4/5"></div>
        </div>
      </div>
      <div v-else-if="logsError" class="text-center py-8">
        <p class="text-destructive mb-3">Failed to load audit logs</p>
        <button @click="logsRefetch()" class="px-4 py-2 rounded text-sm border">Retry</button>
      </div>
      <div v-else-if="!(logs as any[])?.length" class="text-center py-8 text-muted-foreground">
        No audit logs found
      </div>
      <div v-else class="space-y-2">
        <div v-for="log in (logs as any[])" :key="log.id" class="border rounded p-3 text-sm">
          <div class="flex justify-between">
            <span class="font-medium">{{ log.event_type }}</span>
            <span class="text-muted-foreground">{{ formatDateTime(log.created_at) }}</span>
          </div>
          <p class="text-muted-foreground mt-1">{{ log.detail }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
