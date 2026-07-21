<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'

const route = useRoute()
const API = import.meta.env.VITE_API_URL || '/api/v1'

const {
  data: movie,
  isLoading: movieLoading,
  isError: movieError,
  refetch: movieRefetch,
} = useQuery({
  queryKey: ['movie', route.params.id],
  queryFn: async () => {
    const res = await fetch(`${API}/movies/${route.params.id}`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movie')
    return res.json()
  },
})

const {
  data: showtimes,
  isLoading: stLoading,
  isError: stError,
  refetch: stRefetch,
} = useQuery({
  queryKey: ['showtimes', route.params.id],
  queryFn: async () => {
    const res = await fetch(`${API}/movies/${route.params.id}/showtimes`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch showtimes')
    return res.json()
  },
})

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString([], { month: 'short', day: 'numeric' })
}
</script>

<template>
  <div>
    <router-link to="/" class="text-sm text-muted-foreground hover:underline mb-4 inline-block">&larr; Back to Movies</router-link>

    <div v-if="movieLoading || stLoading" class="animate-pulse">
      <div class="flex gap-6 mb-8">
        <div class="w-48 h-72 bg-muted rounded-lg shrink-0" />
        <div class="flex-1 space-y-3 pt-2">
          <div class="h-7 bg-muted rounded w-1/2" />
          <div class="h-4 bg-muted rounded w-1/3" />
          <div class="h-4 bg-muted rounded w-2/3" />
          <div class="h-4 bg-muted rounded w-1/4" />
        </div>
      </div>
      <div class="grid gap-4">
        <div v-for="i in 3" :key="i" class="border rounded-lg p-4">
          <div class="h-5 bg-muted rounded w-1/4 mb-2" />
          <div class="h-4 bg-muted rounded w-1/2 mb-2" />
          <div class="h-4 bg-muted rounded w-1/6" />
        </div>
      </div>
    </div>

    <div v-else-if="movieError || stError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Failed to load</p>
      <p class="text-sm text-muted-foreground mb-4">Could not load showtime information. Check that the backend is running.</p>
      <button
        @click="
          movieRefetch()
          stRefetch()
        "
        class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90"
      >
        Retry
      </button>
    </div>

    <div v-else-if="(showtimes as any[])?.length === 0" class="text-center py-16">
      <p class="text-lg text-muted-foreground">No showtimes available for this movie</p>
      <p class="text-sm text-muted-foreground mt-1">Check back later for new showtimes.</p>
    </div>

    <template v-else>
      <h1 class="text-2xl font-bold mb-6">{{ (movie as any)?.title }} - Showtimes</h1>
      <div class="grid gap-4">
        <div v-for="st in (showtimes as any[]) || []" :key="st.id" class="border rounded-lg p-4 flex items-center justify-between">
          <div>
            <p class="font-medium">{{ st.screen }}</p>
            <p class="text-sm text-muted-foreground">{{ formatDate(st.start_time) }} at {{ formatTime(st.start_time) }}</p>
          </div>
          <div class="text-right">
            <p class="font-bold">${{ st.price?.toFixed(2) }}</p>
            <router-link
              :to="`/showtimes/${st.id}`"
              class="mt-1 inline-block px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90"
            >
              Select Seats
            </router-link>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
