<script setup lang="ts">
import { Card } from '@/components/ui/card'
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Admin Dashboard</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <router-link to="/admin/bookings">
        <Card class="p-6 hover:shadow-lg transition-shadow">
          <h2 class="font-bold text-lg">Bookings</h2>
          <p class="text-sm text-muted-foreground mt-1">View and manage all bookings</p>
        </Card>
      </router-link>
      <router-link to="/admin/logs">
        <Card class="p-6 hover:shadow-lg transition-shadow">
          <h2 class="font-bold text-lg">Audit Logs</h2>
          <p class="text-sm text-muted-foreground mt-1">View system audit trail</p>
        </Card>
      </router-link>
      <router-link to="/admin/movies">
        <Card class="p-6 hover:shadow-lg transition-shadow">
          <h2 class="font-bold text-lg">Movies</h2>
          <p class="text-sm text-muted-foreground mt-1">Add, edit, or remove movies</p>
        </Card>
      </router-link>
      <router-link to="/admin/showtimes">
        <Card class="p-6 hover:shadow-lg transition-shadow">
          <h2 class="font-bold text-lg">Showtimes</h2>
          <p class="text-sm text-muted-foreground mt-1">Manage showtimes and schedules</p>
        </Card>
      </router-link>
      <router-link to="/admin/users">
        <Card class="p-6 hover:shadow-lg transition-shadow">
          <h2 class="font-bold text-lg">Users</h2>
          <p class="text-sm text-muted-foreground mt-1">Manage user roles and permissions</p>
        </Card>
      </router-link>
    </div>
  </div>
</template>
