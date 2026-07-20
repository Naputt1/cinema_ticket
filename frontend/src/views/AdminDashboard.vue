<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'
const activeTab = ref<'bookings' | 'logs' | 'movies' | 'showtimes'>('bookings')

const showMovieModal = ref(false)
const editingMovie = ref<any>(null)
const movieForm = ref({ title: '', description: '', genre: '', duration: 60, poster_url: '' })

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

const { data: movies, isLoading: moviesLoading, isError: moviesError, refetch: moviesRefetch } = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

const showShowtimeModal = ref(false)
const editingShowtime = ref<any>(null)
const selectedShowtimeMovieId = ref('')
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
  selectedShowtimeMovieId.value = ''
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
    if (!selectedShowtimeMovieId.value) {
      alert('Please select a movie')
      return
    }
    await fetch(`${API}/admin/movies/${selectedShowtimeMovieId.value}/showtimes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  }
  showShowtimeModal.value = false
  selectedShowtimeMovieId.value = ''
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

function movieTitleForShowtime(st: any) {
  const m = (movies.value as any[] || []).find(m => m.id === st.movie_id)
  return m?.title || st.movie_id?.slice(0, 8)
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
      <button
        @click="activeTab = 'showtimes'"
        class="px-4 py-2 rounded text-sm"
        :class="activeTab === 'showtimes' ? 'bg-primary text-primary-foreground' : 'border'"
      >
        Showtimes
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

      <div v-else-if="moviesError" class="text-center py-8">
        <p class="text-destructive mb-3">Failed to load movies</p>
        <button @click="moviesRefetch()" class="px-4 py-2 rounded text-sm border">Retry</button>
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
            <p class="font-medium">{{ movieTitleForShowtime(st) }}</p>
            <p class="text-muted-foreground text-xs">{{ st.screen }} · {{ formatSTDateTime(st.start_time) }} · ${{ st.price?.toFixed(2) }}</p>
          </div>
          <div class="flex gap-2">
            <button @click="openEditShowtime(st)" class="text-xs px-2 py-1 border rounded hover:bg-accent">Edit</button>
            <button @click="deleteShowtime(st.id)" class="text-xs px-2 py-1 border border-destructive text-destructive rounded hover:bg-destructive/5">Delete</button>
          </div>
        </div>
      </div>
    </div>
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

  <div v-if="showShowtimeModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showShowtimeModal = false">
    <div class="bg-background rounded-lg p-6 w-full max-w-md">
      <h2 class="text-lg font-bold mb-4">{{ editingShowtime ? 'Edit Showtime' : 'Add Showtime' }}</h2>
      <div class="space-y-3">
        <select v-model="selectedShowtimeMovieId" class="w-full px-3 py-2 border rounded text-sm">
          <option value="">Select movie...</option>
          <option v-for="m in (movies as any[] || [])" :key="m.id" :value="m.id">{{ m.title }}</option>
        </select>
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
