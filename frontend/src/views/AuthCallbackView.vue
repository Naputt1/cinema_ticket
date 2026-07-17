<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { toast } from 'vue-sonner'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

onMounted(async () => {
  const code = route.query.code as string
  if (!code) {
    router.replace('/login')
    return
  }

  try {
    await auth.handleGoogleCallback(code)
    const redirect = sessionStorage.getItem('redirect') || '/'
    sessionStorage.removeItem('redirect')
    router.replace(redirect)
  } catch {
    toast.error('Authentication failed. Please try again.')
    router.replace('/login')
  }
})
</script>

<template>
  <div class="flex items-center justify-center min-h-[60vh]">
    <div class="text-center space-y-4">
      <div class="animate-spin w-8 h-8 border-4 border-primary border-t-transparent rounded-full mx-auto" />
      <p class="text-muted-foreground">Signing you in...</p>
    </div>
  </div>
</template>
