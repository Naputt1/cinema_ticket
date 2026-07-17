<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref<'bookings' | 'logs' | 'movies'>('bookings')

const showMovieModal = ref(false)
const editingMovie = ref<any>(null)
const movieForm = ref({ title: '', description: '', genre: '', duration: 60, poster_url: '' })

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

const { data: movies, isLoading: moviesLoading, isError: moviesError, refetch: moviesRefetch } = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
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

function openCreateMovie() {
  editingMovie.value = null
  movieForm.value = { title: '', description: '', genre: '', duration: 60, poster_url: '' }
  showMovieModal.value = true
}

function openEditMovie(m: any) {
  editingMovie.value = m
  movieForm.value = { title: m.title, description: m.description, genre: m.genre, duration: m.duration, poster_url: m.poster_url }
  showMovieModal.value = true
}

async function saveMovie() {
  const url = editingMovie.value
    ? `${API}/admin/movies/${editingMovie.value.id}`
    : `${API}/admin/movies`
  const method = editingMovie.value ? 'PUT' : 'POST'
  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(movieForm.value),
  })
  if (!res.ok) throw new Error('Failed to save movie')
  showMovieModal.value = false
  moviesRefetch()
}

async function deleteMovie(id: string) {
  if (!confirm('Delete this movie?')) return
  await fetch(`${API}/admin/movies/${id}`, { method: 'DELETE', credentials: 'include' })
  moviesRefetch()
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
        @click="activeTab = 'movies'"
        class="px-4 py-2 rounded text-sm"
        :class="activeTab === 'movies' ? 'bg-primary text-primary-foreground' : 'border'"
      >
        Movies
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

    <div v-if="activeTab === 'movies'">
      <button @click="openCreateMovie" class="mb-4 px-4 py-2 bg-primary text-primary-foreground rounded text-sm">
        Add Movie
      </button>

      <div v-if="moviesLoading" class="space-y-2">
        <div v-for="i in 4" :key="i" class="border rounded p-3 animate-pulse">
          <div class="h-4 bg-muted rounded w-1/3 mb-2"></div>
          <div class="h-3 bg-muted rounded w-1/2"></div>
        </div>
      </div>

      <div v-else-if="moviesError" class="text-red-500">
        Failed to load movies. <button @click="() => moviesRefetch()" class="underline">Retry</button>
      </div>

      <div v-else-if="!movies?.length" class="text-muted-foreground">No movies yet</div>

      <table v-else class="w-full text-sm">
        <thead>
          <tr class="border-b text-left">
            <th class="pb-2">Poster</th>
            <th class="pb-2">Title</th>
            <th class="pb-2">Genre</th>
            <th class="pb-2">Duration</th>
            <th class="pb-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="m in (movies as any[])" :key="m.id" class="border-b">
            <td class="py-2">
              <img v-if="m.poster_url" :src="m.poster_url" alt="" class="w-10 h-14 object-cover rounded" />
              <div v-else class="w-10 h-14 bg-muted rounded"></div>
            </td>
            <td class="py-2 font-medium">{{ m.title }}</td>
            <td class="py-2">{{ m.genre }}</td>
            <td class="py-2">{{ m.duration }} min</td>
            <td class="py-2 flex gap-2">
              <button @click="openEditMovie(m)" class="text-primary underline">Edit</button>
              <button @click="deleteMovie(m.id)" class="text-red-500 underline">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showMovieModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showMovieModal = false">
      <div class="bg-background rounded-lg p-6 w-full max-w-md">
        <h2 class="text-lg font-bold mb-4">{{ editingMovie ? 'Edit Movie' : 'Add Movie' }}</h2>
        <div class="space-y-3">
          <input v-model="movieForm.title" placeholder="Title" class="w-full px-3 py-2 border rounded text-sm" />
          <input v-model="movieForm.genre" placeholder="Genre" class="w-full px-3 py-2 border rounded text-sm" />
          <input v-model="movieForm.duration" type="number" placeholder="Duration (min)" class="w-full px-3 py-2 border rounded text-sm" />
          <textarea v-model="movieForm.description" placeholder="Description" class="w-full px-3 py-2 border rounded text-sm" rows="3" />
          <input v-model="movieForm.poster_url" placeholder="Poster URL" class="w-full px-3 py-2 border rounded text-sm" />
        </div>
        <div class="flex justify-end gap-2 mt-4">
          <button @click="showMovieModal = false" class="px-4 py-2 border rounded text-sm">Cancel</button>
          <button @click="saveMovie" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">{{ editingMovie ? 'Update' : 'Create' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
