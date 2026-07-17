<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const { data: movies, isLoading } = useQuery({
  queryKey: ['movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Now Showing</h1>
    <div v-if="isLoading" class="text-center py-8 text-muted-foreground">Loading...</div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
      <div
        v-for="movie in (movies as any[] || [])"
        :key="movie.id"
        class="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
      >
        <img :src="movie.poster_url" :alt="movie.title" class="w-full h-64 object-cover" />
        <div class="p-4">
          <h2 class="font-bold text-lg">{{ movie.title }}</h2>
          <p class="text-sm text-muted-foreground mt-1">{{ movie.genre }} · {{ movie.duration }} min</p>
          <p class="text-sm mt-2 line-clamp-2">{{ movie.description }}</p>
          <router-link
            :to="`/movies/${movie.id}`"
            class="mt-3 inline-block px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90"
          >
            Select Showtime
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
