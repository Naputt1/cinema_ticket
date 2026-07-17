<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'
import { useSSE } from './composables/useSSE'

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
        <router-link to="/admin" class="text-sm hover:underline" v-if="auth.isAdmin">Admin</router-link>
        <template v-if="auth.isAuthenticated">
          <span class="text-sm text-muted-foreground">{{ auth.user?.name }}</span>
          <button @click="auth.logout()" class="text-sm text-destructive hover:underline">Logout</button>
        </template>
        <router-link v-else to="/login" class="text-sm hover:underline">Login</router-link>
      </div>
    </nav>
    <main class="container mx-auto p-4">
      <router-view />
    </main>
    <Toaster />
  </div>
</template>
