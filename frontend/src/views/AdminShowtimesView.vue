<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogClose } from '@/components/ui/dialog'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const showShowtimeModal = ref(false)
const editingShowtime = ref<any>(null)
const selectedShowtimeMovieId = ref('')
const showtimeForm = ref({ screen: '', start_time: '', price: 12 })

const { data: movies } = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

const {
  data: showtimes,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['admin-showtimes'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/showtimes`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch showtimes')
    return res.json()
  },
})

function openCreateShowtime() {
  editingShowtime.value = null
  selectedShowtimeMovieId.value = ''
  showtimeForm.value = { screen: '', start_time: '', price: 12 }
  showShowtimeModal.value = true
}

function openEditShowtime(st: any) {
  editingShowtime.value = st
  showtimeForm.value = { screen: st.screen, start_time: st.start_time?.slice(0, 16), price: st.price }
  showShowtimeModal.value = true
}

async function saveShowtime() {
  const body = {
    screen: showtimeForm.value.screen,
    start_time: new Date(showtimeForm.value.start_time).toISOString(),
    price: Number(showtimeForm.value.price),
  }

  if (editingShowtime.value) {
    await fetch(`${API}/admin/showtimes/${editingShowtime.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  } else {
    if (!selectedShowtimeMovieId.value) {
      alert('Please select a movie')
      return
    }
    await fetch(`${API}/admin/movies/${selectedShowtimeMovieId.value}/showtimes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  }
  showShowtimeModal.value = false
  selectedShowtimeMovieId.value = ''
  refetch()
}

async function deleteShowtime(id: string) {
  if (!confirm('Delete this showtime?')) return
  await fetch(`${API}/admin/showtimes/${id}`, { method: 'DELETE', credentials: 'include' })
  refetch()
}

function movieTitleForShowtime(st: any) {
  return `Movie ID: ${st.movie_id}`
}

function formatSTDateTime(iso: string) {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Showtimes</h1>

    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-semibold">All Showtimes</h2>
      <Button @click="openCreateShowtime()" class="shrink-0">Add Showtime</Button>
    </div>

    <div v-if="isLoading" class="space-y-2">
      <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-2">Failed to load showtimes</p>
      <Button variant="outline" @click="refetch()">Retry</Button>
    </div>

    <div v-else-if="!showtimes?.length" class="text-center py-8 text-muted-foreground">No showtimes yet.</div>

    <Table v-else>
      <TableHeader>
        <TableRow>
          <TableHead>Movie</TableHead>
          <TableHead>Screen</TableHead>
          <TableHead>Start Time</TableHead>
          <TableHead>Price</TableHead>
          <TableHead class="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="st in showtimes" :key="st.id">
          <TableCell class="font-medium">{{ st.movie_title || movieTitleForShowtime(st) }}</TableCell>
          <TableCell>{{ st.screen }}</TableCell>
          <TableCell>{{ formatSTDateTime(st.start_time) }}</TableCell>
          <TableCell>${{ st.price?.toFixed(2) }}</TableCell>
          <TableCell class="text-right">
            <div class="flex justify-end gap-2">
              <Button variant="outline" size="xs" @click="openEditShowtime(st)">Edit</Button>
              <Button variant="destructive" size="xs" @click="deleteShowtime(st.id)">Delete</Button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <Dialog v-model:open="showShowtimeModal">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ editingShowtime ? 'Edit Showtime' : 'Add Showtime' }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-3">
          <Select v-model="selectedShowtimeMovieId">
            <SelectTrigger class="w-full">
              <SelectValue placeholder="Select movie..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="m in (movies as any[]) || []" :key="m.id" :value="m.id">{{ m.title }}</SelectItem>
            </SelectContent>
          </Select>
          <Input v-model="showtimeForm.screen" placeholder="Screen (e.g. Screen 1)" />
          <Input v-model="showtimeForm.start_time" type="datetime-local" />
          <Input v-model="showtimeForm.price" type="number" step="0.01" placeholder="Price" />
        </div>
        <DialogFooter>
          <DialogClose as-child>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button @click="saveShowtime">{{ editingShowtime ? 'Update' : 'Create' }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
