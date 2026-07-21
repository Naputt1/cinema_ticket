<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const showShowtimeModal = ref(false)
const editingShowtime = ref<any>(null)
const selectedShowtimeMovieId = ref('')
const showtimeForm = ref({ screen: '', start_time: '', price: 12 })

const { data: movies } = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

const {
  data: showtimes,
  isLoading,
  isError,
  refetch,
} = useQuery({
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
  showtimeForm.value = { screen: st.screen, start_time: st.start_time?.slice(0, 16), price: st.price }
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
  refetch()
}

async function deleteShowtime(id: string) {
  if (!confirm('Delete this showtime?')) return
  await fetch(`${API}/admin/showtimes/${id}`, { method: 'DELETE', credentials: 'include' })
  refetch()
}

function movieTitleForShowtime(st: any) {
  return `Movie ID: ${st.movie_id}`
}

function formatSTDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Showtimes</h1>

    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">All Showtimes</h2>
      <button @click="openCreateShowtime()" class="px-3 py-1.5 bg-primary text-primary-foreground rounded text-sm">Add Showtime</button>
    </div>

    <div v-if="isLoading" class="space-y-2">
      <div v-for="i in 5" :key="i" class="h-12 bg-muted rounded animate-pulse" />
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-2">Failed to load showtimes</p>
      <button @click="refetch()" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Retry</button>
    </div>

    <div v-else-if="!showtimes?.length" class="text-center py-8 text-muted-foreground">No showtimes yet.</div>

    <div v-else class="space-y-2">
      <div v-for="st in showtimes" :key="st.id" class="border rounded p-3 text-sm flex items-center justify-between">
        <div>
          <p class="font-medium">{{ st.movie_title || movieTitleForShowtime(st) }}</p>
          <p class="text-muted-foreground text-xs">{{ st.screen }} · {{ formatSTDateTime(st.start_time) }} · ${{ st.price?.toFixed(2) }}</p>
        </div>
        <div class="flex gap-2">
          <button @click="openEditShowtime(st)" class="text-xs px-2 py-1 border rounded hover:bg-accent">Edit</button>
          <button
            @click="deleteShowtime(st.id)"
            class="text-xs px-2 py-1 border border-destructive text-destructive rounded hover:bg-destructive/5"
          >
            Delete
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showShowtimeModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @click.self="showShowtimeModal = false"
    >
      <div class="bg-background rounded-lg p-6 w-full max-w-md">
        <h2 class="text-lg font-bold mb-4">{{ editingShowtime ? 'Edit Showtime' : 'Add Showtime' }}</h2>
        <div class="space-y-3">
          <select v-model="selectedShowtimeMovieId" class="w-full px-3 py-2 border rounded text-sm">
            <option value="">Select movie...</option>
            <option v-for="m in (movies as any[]) || []" :key="m.id" :value="m.id">{{ m.title }}</option>
          </select>
          <input v-model="showtimeForm.screen" placeholder="Screen (e.g. Screen 1)" class="w-full px-3 py-2 border rounded text-sm" />
          <input v-model="showtimeForm.start_time" type="datetime-local" class="w-full px-3 py-2 border rounded text-sm" />
          <input
            v-model="showtimeForm.price"
            type="number"
            step="0.01"
            placeholder="Price"
            class="w-full px-3 py-2 border rounded text-sm"
          />
        </div>
        <div class="flex justify-end gap-2 mt-4">
          <button @click="showShowtimeModal = false" class="px-4 py-2 border rounded text-sm">Cancel</button>
          <button @click="saveShowtime" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm">
            {{ editingShowtime ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
