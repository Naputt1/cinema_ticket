<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const {
  data: logs,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['admin-audit-logs'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/audit-logs`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch logs')
    return res.json()
  },
})

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Audit Logs</h1>

    <div v-if="isLoading" class="space-y-2">
      <div v-for="n in 3" :key="n" class="border rounded p-3 text-sm animate-pulse">
        <div class="h-4 bg-muted rounded w-2/5 mb-2" />
        <div class="h-3 bg-muted rounded w-4/5" />
      </div>
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-3">Failed to load audit logs</p>
      <Button variant="outline" @click="refetch()">Retry</Button>
    </div>

    <div v-else-if="!(logs as any[])?.length" class="text-center py-8 text-muted-foreground">No audit logs found</div>

    <div v-else class="space-y-2">
      <div v-for="log in logs as any[]" :key="log.id" class="border rounded p-3 text-sm">
        <div class="flex justify-between">
          <span class="font-medium">{{ log.event_type }}</span>
          <span class="text-muted-foreground">{{ formatDateTime(log.created_at) }}</span>
        </div>
        <p class="text-muted-foreground mt-1">{{ log.detail }}</p>
      </div>
    </div>
  </div>
</template>
