<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent } from '@/components/ui/card'

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

function retryAll() {
  movieRefetch()
  stRefetch()
}

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

    <div v-if="movieLoading || stLoading">
      <div class="flex gap-6 mb-8">
        <Skeleton class="w-48 h-72 rounded-lg shrink-0" />
        <div class="flex-1 space-y-3 pt-2">
          <Skeleton class="h-7 w-1/2" />
          <Skeleton class="h-4 w-1/3" />
          <Skeleton class="h-4 w-2/3" />
          <Skeleton class="h-4 w-1/4" />
        </div>
      </div>
      <div class="grid gap-4">
        <Card v-for="i in 3" :key="i">
          <CardContent class="space-y-2">
            <Skeleton class="h-5 w-1/4" />
            <Skeleton class="h-4 w-1/2" />
            <Skeleton class="h-4 w-1/6" />
          </CardContent>
        </Card>
      </div>
    </div>

    <div v-else-if="movieError || stError" class="text-center py-16">
      <p class="text-lg text-destructive mb-2">Failed to load</p>
      <p class="text-sm text-muted-foreground mb-4">Could not load showtime information. Check that the backend is running.</p>
      <Button @click="retryAll">Retry</Button>
    </div>

    <div v-else-if="(showtimes as any[])?.length === 0" class="text-center py-16">
      <p class="text-lg text-muted-foreground">No showtimes available for this movie</p>
      <p class="text-sm text-muted-foreground mt-1">Check back later for new showtimes.</p>
    </div>

    <template v-else>
      <h1 class="text-2xl font-bold mb-6">{{ (movie as any)?.title }} - Showtimes</h1>
      <div class="grid gap-4">
        <Card v-for="st in (showtimes as any[]) || []" :key="st.id">
          <CardContent class="flex items-center justify-between !py-0">
            <div>
              <p class="font-medium">{{ st.screen }}</p>
              <p class="text-sm text-muted-foreground">{{ formatDate(st.start_time) }} at {{ formatTime(st.start_time) }}</p>
            </div>
            <div class="text-right">
              <p class="font-bold">${{ st.price?.toFixed(2) }}</p>
              <router-link :to="`/showtimes/${st.id}`">
                <Button size="sm"> Select Seats </Button>
              </router-link>
            </div>
          </CardContent>
        </Card>
      </div>
    </template>
  </div>
</template>
