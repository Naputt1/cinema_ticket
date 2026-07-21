<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>
<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] text-center">
    <h1 class="text-6xl font-bold text-muted-foreground mb-4">404</h1>
    <p class="text-lg text-muted-foreground mb-6">Page not found</p>
    <button @click="router.push('/')" class="px-4 py-2 bg-primary text-primary-foreground rounded text-sm hover:opacity-90">Go Home</button>
  </div>
</template>
