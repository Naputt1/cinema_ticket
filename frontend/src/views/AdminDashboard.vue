<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref<'bookings' | 'logs' | 'showtimes'>('bookings')

const { data: bookings } = useQuery({
  queryKey: ['admin-bookings'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/bookings`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch bookings')
    return res.json()
  },
})

const { data: logs } = useQuery({
  queryKey: ['admin-audit-logs'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/audit-logs`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch logs')
    return res.json()
  },
})

const showShowtimeModal = ref(false)
const editingShowtime = ref<any>(null)
const showtimeForm = ref({ screen: '', start_time: '', price: 12 })

const { data: showtimes, isLoading: stLoading, isError: stError, refetch: stRefetch } = useQuery({
  queryKey: ['admin-showtimes'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/showtimes`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch showtimes')
    return res.json()
  },
})

function openCreateShowtime() {
  editingShowtime.value = null
  showtimeForm.value = { screen: '', start_time: '', price: 12 }
  showShowtimeModal.value = true
}

function openEditShowtime(st: any) {
  editingShowtime.value = st
  showtimeForm.value = {
    screen: st.screen,
    start_time: st.start_time?.slice(0, 16),
    price: st.price,
  }
  showShowtimeModal.value = true
}

async function saveShowtime() {
  const body = {
    screen: showtimeForm.value.screen,
    start_time: new Date(showtimeForm.value.start_time).toISOString(),
    price: Number(showtimeForm.value.price),
  }

  if (editingShowtime.value) {
    await fetch(`${API}/admin/showtimes/${editingShowtime.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  } else {
    const movieId = prompt('Movie ID:')
    if (!movieId) return
    await fetch(`${API}/admin/movies/${movieId}/showtimes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  }
  showShowtimeModal.value = false
  stRefetch()
}

async function deleteShowtime(id: string) {
  if (!confirm('Delete this showtime?')) return
  await fetch(`${API}/admin/showtimes/${id}`, { method: 'DELETE', credentials: 'include' })
  stRefetch()
}

function formatSTDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}

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
      <button
        @click="activeTab = 'showtimes'"
        class="px-4 py-2 rounded text-sm"
        :class="activeTab === 'showtimes' ? 'bg-primary text-primary-foreground' : 'border'"
      >
        Showtimes
      </button>
    </div>

    <div v-if="activeTab === 'bookings'">
      <div class="space-y-2">
        <div v-for="b in (bookings as any[] || [])" :key="b.id" class="border rounded p-3 text-sm">
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
      <div class="space-y-2">
        <div v-for="log in (logs as any[] || [])" :key="log.id" class="border rounded p-3 text-sm">
          <div class="flex justify-between">
            <span class="font-medium">{{ log.event_type }}</span>
            <span class="text-muted-foreground">{{ formatDateTime(log.created_at) }}</span>
          </div>
          <p class="text-muted-foreground mt-1">{{ log.detail }}</p>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'showtimes'">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold">Showtimes</h2>
        <button @click="openCreateShowtime()" class="px-3 py-1.5 bg-primary text-primary-foreground rounded text-sm">Add Showtime</button>
      </div>

      <div v-if="stLoading" class="space-y-2">
        <div v-for="i in 5" :key="i" class="h-12 bg-muted rounded animate-pulse" />
      </div>
      <div v-else-if="stError" class="text-center py-8">
        <p class="text-destructive mb-2">Failed to load showtimes</p>
        <button @click="stRefetch()" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Retry</button>
      </div>
      <div v-else-if="!showtimes?.length" class="text-center py-8 text-muted-foreground">No showtimes yet.</div>
      <div v-else class="space-y-2">
        <div v-for="st in showtimes" :key="st.id" class="border rounded p-3 text-sm flex items-center justify-between">
          <div>
            <p class="font-medium">{{ st.screen }}</p>
            <p class="text-muted-foreground">{{ formatSTDateTime(st.start_time) }} · ${{ st.price?.toFixed(2) }}</p>
            <p class="text-xs text-muted-foreground">Movie ID: {{ st.movie_id }}</p>
          </div>
          <div class="flex gap-2">
            <button @click="openEditShowtime(st)" class="text-xs px-2 py-1 border rounded hover:bg-accent">Edit</button>
            <button @click="deleteShowtime(st.id)" class="text-xs px-2 py-1 border border-destructive text-destructive rounded hover:bg-destructive/5">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-if="showShowtimeModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showShowtimeModal = false">
    <div class="bg-background rounded-lg p-6 w-full max-w-md">
      <h2 class="text-lg font-bold mb-4">{{ editingShowtime ? 'Edit Showtime' : 'Add Showtime' }}</h2>
      <div class="space-y-3">
        <input v-model="showtimeForm.screen" placeholder="Screen (e.g. Screen 1)" class="w-full px-3 py-2 border rounded text-sm" />
        <input v-model="showtimeForm.start_time" type="datetime-local" class="w-full px-3 py-2 border rounded text-sm" />
        <input v-model="showtimeForm.price" type="number" step="0.01" placeholder="Price" class="w-full px-3 py-2 border rounded text-sm" />
      </div>
      <div class="flex justify-end gap-2 mt-4">
        <button @click="showShowtimeModal = false" class="px-4 py-2 border rounded text-sm">Cancel</button>
        <button @click="saveShowtime" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">{{ editingShowtime ? 'Update' : 'Create' }}</button>
      </div>
    </div>
  </div>
</template>
