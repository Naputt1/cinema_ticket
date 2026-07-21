<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useSSE } from './composables/useSSE'

const route = useRoute()
const auth = useAuthStore()
auth.loadFromStorage()

onMounted(() => {
  auth.checkAuth()
})

useSSE()
</script>

<template>
  <div class="min-h-screen bg-background">
    <nav class="border-b px-4 py-3 flex items-center justify-between">
      <router-link to="/" class="font-bold text-lg">🎬 Cinema Ticket</router-link>
      <div class="flex items-center gap-4">
        <router-link to="/bookings" class="text-sm hover:underline" v-if="auth.isAuthenticated">My Bookings</router-link>
        <template v-if="auth.isAuthenticated">
          <span class="text-sm text-muted-foreground">{{ auth.user?.name }}</span>
          <button @click="auth.logout()" class="text-sm text-destructive hover:underline">Logout</button>
        </template>
        <router-link v-else to="/login" class="text-sm hover:underline">Login</router-link>
      </div>
    </nav>
    <div class="flex">
      <aside v-if="route.path.startsWith('/admin')" class="w-48 border-r min-h-[calc(100vh-57px)] p-4 shrink-0 flex flex-col gap-1">
        <div class="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Admin</div>
        <router-link to="/admin" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Dashboard</router-link
        >
        <router-link to="/admin/bookings" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Bookings</router-link
        >
        <router-link to="/admin/logs" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Audit Logs</router-link
        >
        <router-link to="/admin/movies" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Movies</router-link
        >
        <router-link to="/admin/showtimes" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Showtimes</router-link
        >
        <router-link to="/admin/users" class="px-3 py-2 rounded text-sm hover:bg-accent" active-class="bg-accent font-medium"
          >Users</router-link
        >
      </aside>
      <main class="container mx-auto p-4 flex-1">
        <router-view />
      </main>
    </div>
    <Toaster />
  </div>
</template>
