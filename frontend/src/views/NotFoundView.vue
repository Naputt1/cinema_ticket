<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
const router = useRouter()
</script>
<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] text-center">
    <h1 class="text-6xl font-bold text-muted-foreground mb-4">404</h1>
    <p class="text-lg text-muted-foreground mb-6">Page not found</p>
    <Button @click="router.push('/')">Go Home</Button>
  </div>
</template>
