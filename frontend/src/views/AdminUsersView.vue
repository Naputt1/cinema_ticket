<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const {
  data: users,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['admin-users'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/users`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch users')
    return res.json()
  },
})

async function toggleUserRole(u: any) {
  const newRole = u.role === 'ADMIN' ? 'USER' : 'ADMIN'
  await fetch(`${API}/admin/users/${u.id}/role`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ role: newRole }),
  })
  refetch()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Users</h1>

    <div v-if="isLoading" class="space-y-2">
      <div v-for="i in 5" :key="i" class="h-12 bg-muted rounded animate-pulse" />
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-3">Failed to load users</p>
      <Button variant="outline" @click="refetch()">Retry</Button>
    </div>

    <div v-else-if="!users?.length" class="text-center py-8 text-muted-foreground">No users found.</div>

    <table v-else class="w-full text-sm">
      <thead>
        <tr class="border-b text-left">
          <th class="pb-2">Email</th>
          <th class="pb-2">Name</th>
          <th class="pb-2">Role</th>
          <th class="pb-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="u in users as any[]" :key="u.id" class="border-b">
          <td class="py-2">{{ u.email }}</td>
          <td class="py-2">{{ u.name }}</td>
          <td class="py-2">
            <span
              class="text-xs font-medium px-2 py-0.5 rounded"
              :class="u.role === 'ADMIN' ? 'text-purple-700 bg-purple-50' : 'text-blue-700 bg-blue-50'"
            >
              {{ u.role }}
            </span>
          </td>
          <td class="py-2">
            <Button variant="outline" size="xs" @click="toggleUserRole(u)">
              {{ u.role === 'ADMIN' ? 'Remove Admin' : 'Make Admin' }}
            </Button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
