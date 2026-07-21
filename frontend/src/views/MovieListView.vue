<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const search = ref('')
const selectedGenre = ref('')

const {
  data: movies,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

const genres = computed(() => {
  const all = ((movies.value as any[]) || []).map((m) => m.genre).filter(Boolean)
  return [...new Set(all)] as string[]
})

const filtered = computed(() => {
  const list = ((movies.value as any[]) || []).filter((m) => {
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
      <Input v-model="search" type="text" placeholder="Search movies..." class="w-full sm:w-64" />
    </div>

    <div class="flex flex-wrap gap-2 mb-6">
      <Badge :variant="selectedGenre === '' ? 'default' : 'outline'" class="cursor-pointer" @click="selectedGenre = ''">All</Badge>
      <Badge
        v-for="g in genres"
        :key="g"
        :variant="selectedGenre === g ? 'default' : 'outline'"
        class="cursor-pointer"
        @click="selectedGenre = g"
      >
        {{ g }}
      </Badge>
    </div>

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
      <div v-for="i in 8" :key="i" class="border rounded-lg overflow-hidden">
        <Skeleton class="w-full h-64 rounded-none" />
        <div class="p-4 space-y-3">
          <Skeleton class="h-5 w-3/4" />
          <Skeleton class="h-4 w-1/2" />
          <Skeleton class="h-4 w-full" />
          <Skeleton class="h-4 w-2/3" />
          <Skeleton class="h-9 w-32" />
        </div>
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Failed to load movies</p>
      <p class="text-sm text-muted-foreground mb-4">Check that the backend is running.</p>
      <Button @click="refetch()">Retry</Button>
    </div>

    <div v-else-if="filtered.length === 0" class="text-center py-16">
      <p class="text-lg text-muted-foreground">No movies found</p>
      <p class="text-sm text-muted-foreground mt-1">Try adjusting your search or filter.</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
      <Card
        v-for="movie in filtered"
        :key="movie.id"
        class="flex flex-col overflow-hidden"
      >
        <div class="relative">
          <img :src="movie.poster_url" :alt="movie.title" class="w-full h-64 object-cover" @error="imgError" />
          <div class="hidden absolute inset-0 items-center justify-center bg-muted text-muted-foreground text-sm">
            {{ movie.title?.[0] }}
          </div>
        </div>
        <div class="p-4 flex flex-col flex-1">
          <div class="flex items-start justify-between gap-2 mb-1">
            <h2 class="font-bold text-lg leading-tight">{{ movie.title }}</h2>
            <Badge variant="secondary">
              {{ movie.genre }}
            </Badge>
          </div>
          <p class="text-xs text-muted-foreground mb-2">{{ formatDuration(movie.duration) }}</p>
          <p class="text-sm text-muted-foreground line-clamp-2 mb-4">{{ movie.description }}</p>
          <router-link :to="`/movies/${movie.id}`">
            <Button size="default" class="mt-auto self-start">
              Select Showtime
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </Button>
          </router-link>
        </div>
      </Card>
    </div>
  </div>
</template>
