<script setup lang="ts">
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Admin Dashboard</h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <router-link to="/admin/bookings" class="border rounded-lg p-6 hover:shadow-lg transition-shadow">
        <h2 class="font-bold text-lg">Bookings</h2>
        <p class="text-sm text-muted-foreground mt-1">View and manage all bookings</p>
      </router-link>
      <router-link to="/admin/logs" class="border rounded-lg p-6 hover:shadow-lg transition-shadow">
        <h2 class="font-bold text-lg">Audit Logs</h2>
        <p class="text-sm text-muted-foreground mt-1">View system audit trail</p>
      </router-link>
      <router-link to="/admin/movies" class="border rounded-lg p-6 hover:shadow-lg transition-shadow">
        <h2 class="font-bold text-lg">Movies</h2>
        <p class="text-sm text-muted-foreground mt-1">Add, edit, or remove movies</p>
      </router-link>
      <router-link to="/admin/showtimes" class="border rounded-lg p-6 hover:shadow-lg transition-shadow">
        <h2 class="font-bold text-lg">Showtimes</h2>
        <p class="text-sm text-muted-foreground mt-1">Manage showtimes and schedules</p>
      </router-link>
      <router-link to="/admin/users" class="border rounded-lg p-6 hover:shadow-lg transition-shadow">
        <h2 class="font-bold text-lg">Users</h2>
        <p class="text-sm text-muted-foreground mt-1">Manage user roles and permissions</p>
      </router-link>
    </div>
  </div>
</template>
