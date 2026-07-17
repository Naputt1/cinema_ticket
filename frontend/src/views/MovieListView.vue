<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const search = ref('')
const selectedGenre = ref('')

const { data: movies, isLoading, isError, refetch } = useQuery({
  queryKey: ['movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

const genres = computed(() => {
  const all = (movies.value as any[] || []).map(m => m.genre).filter(Boolean)
  return [...new Set(all)] as string[]
})

const filtered = computed(() => {
  const list = (movies.value as any[] || []).filter(m => {
    if (search.value) {
      const q = search.value.toLowerCase()
      if (!m.title?.toLowerCase().includes(q)) return false
    }
    if (selectedGenre.value && m.genre !== selectedGenre.value) return false
    return true
  })
  return list
})

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
  const fallback = el.nextElementSibling as HTMLElement
  if (fallback) fallback.style.display = 'flex'
}
</script>

<template>
  <div>
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <h1 class="text-3xl font-bold">Now Showing</h1>
      <input
        v-model="search"
        type="text"
        placeholder="Search movies..."
        class="w-full sm:w-64 px-3 py-2 border rounded-lg text-sm bg-background focus:outline-none focus:ring-2 focus:ring-primary"
      />
    </div>

    <div class="flex flex-wrap gap-2 mb-6">
      <button
        @click="selectedGenre = ''"
        class="px-3 py-1.5 text-sm rounded-full border transition-colors"
        :class="selectedGenre === '' ? 'bg-primary text-primary-foreground border-primary' : 'hover:bg-accent'"
      >
        All
      </button>
      <button
        v-for="g in genres"
        :key="g"
        @click="selectedGenre = g"
        class="px-3 py-1.5 text-sm rounded-full border transition-colors"
        :class="selectedGenre === g ? 'bg-primary text-primary-foreground border-primary' : 'hover:bg-accent'"
      >
        {{ g }}
      </button>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
      <div v-for="i in 8" :key="i" class="border rounded-lg overflow-hidden animate-pulse">
        <div class="w-full h-64 bg-muted" />
        <div class="p-4 space-y-3">
          <div class="h-5 bg-muted rounded w-3/4" />
          <div class="h-4 bg-muted rounded w-1/2" />
          <div class="h-4 bg-muted rounded w-full" />
          <div class="h-4 bg-muted rounded w-2/3" />
          <div class="h-9 bg-muted rounded w-32" />
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Failed to load movies</p>
      <p class="text-sm text-muted-foreground mb-4">Check that the backend is running.</p>
      <button
        @click="refetch()"
        class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90"
      >
        Retry
      </button>
    </div>

    <div v-else-if="filtered.length === 0" class="text-center py-16">
      <p class="text-lg text-muted-foreground">No movies found</p>
      <p class="text-sm text-muted-foreground mt-1">Try adjusting your search or filter.</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
      <div
        v-for="movie in filtered"
        :key="movie.id"
        class="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow bg-card flex flex-col"
      >
        <div class="relative">
          <img
            :src="movie.poster_url"
            :alt="movie.title"
            class="w-full h-64 object-cover"
            @error="imgError"
          />
          <div
            class="hidden absolute inset-0 items-center justify-center bg-muted text-muted-foreground text-sm"
          >
            {{ movie.title?.[0] }}
          </div>
        </div>
        <div class="p-4 flex flex-col flex-1">
          <div class="flex items-start justify-between gap-2 mb-1">
            <h2 class="font-bold text-lg leading-tight">{{ movie.title }}</h2>
            <span class="shrink-0 text-xs px-2 py-0.5 rounded-full bg-secondary text-secondary-foreground">
              {{ movie.genre }}
            </span>
          </div>
          <p class="text-xs text-muted-foreground mb-2">{{ formatDuration(movie.duration) }}</p>
          <p class="text-sm text-muted-foreground line-clamp-2 mb-4">{{ movie.description }}</p>
          <router-link
            :to="`/movies/${movie.id}`"
            class="mt-auto inline-flex items-center gap-1.5 px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90 self-start"
          >
            Select Showtime
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
