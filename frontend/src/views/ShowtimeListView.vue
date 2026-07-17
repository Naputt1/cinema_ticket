<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'

const route = useRoute()
const API = import.meta.env.VITE_API_URL || '/api/v1'

const { data: movie } = useQuery({
  queryKey: ['movie', route.params.id],
  queryFn: async () => {
    const res = await fetch(`${API}/movies/${route.params.id}`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movie')
    return res.json()
  },
})

const { data: showtimes } = useQuery({
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
    <h1 class="text-2xl font-bold mb-6">{{ (movie as any)?.title }} - Showtimes</h1>
    <div class="grid gap-4">
      <div
        v-for="st in (showtimes as any[] || [])"
        :key="st.id"
        class="border rounded-lg p-4 flex items-center justify-between"
      >
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
  </div>
</template>
