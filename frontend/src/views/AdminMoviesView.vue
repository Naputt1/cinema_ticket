<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const showMovieModal = ref(false)
const editingMovie = ref<any>(null)
const movieForm = ref({ title: '', description: '', genre: '', duration: 60, poster_url: '' })

const { data: movies, isLoading, isError, refetch } = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

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
  const url = editingMovie.value ? `${API}/admin/movies/${editingMovie.value.id}` : `${API}/admin/movies`
  const method = editingMovie.value ? 'PUT' : 'POST'
  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(movieForm.value),
  })
  if (!res.ok) throw new Error('Failed to save movie')
  showMovieModal.value = false
  refetch()
}

async function deleteMovie(id: string) {
  if (!confirm('Delete this movie?')) return
  await fetch(`${API}/admin/movies/${id}`, { method: 'DELETE', credentials: 'include' })
  refetch()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Movies</h1>

    <button @click="openCreateMovie" class="mb-4 px-4 py-2 bg-primary text-primary-foreground rounded text-sm">Add Movie</button>

    <div v-if="isLoading" class="space-y-2">
      <div v-for="i in 4" :key="i" class="border rounded p-3 animate-pulse">
        <div class="h-4 bg-muted rounded w-1/3 mb-2" />
        <div class="h-3 bg-muted rounded w-1/2" />
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-3">Failed to load movies</p>
      <button @click="refetch()" class="px-4 py-2 rounded text-sm border">Retry</button>
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
            <div v-else class="w-10 h-14 bg-muted rounded" />
          </td>
          <td class="py-2 font-medium">{{ m.title }}</td>
          <td class="py-2">{{ m.genre }}</td>
          <td class="py-2">{{ m.duration }} min</td>
          <td class="py-2">
            <div class="flex gap-2 items-center">
              <button @click="openEditMovie(m)" class="text-primary underline">Edit</button>
              <button @click="deleteMovie(m.id)" class="text-red-500 underline">Delete</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

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
