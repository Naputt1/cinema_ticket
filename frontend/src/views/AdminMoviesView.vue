<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent } from '@/components/ui/card'
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog'

const API = import.meta.env.VITE_API_URL || '/api/v1'

const showMovieModal = ref(false)
const editingMovie = ref<any>(null)
const movieForm = ref({ title: '', description: '', genre: '', duration: 60, poster_url: '' })

const {
  data: movies,
  isLoading,
  isError,
  refetch,
} = useQuery({
  queryKey: ['admin-movies'],
  queryFn: async () => {
    const res = await fetch(`${API}/admin/movies`, { credentials: 'include' })
    if (!res.ok) throw new Error('Failed to fetch movies')
    return res.json()
  },
})

function openCreateMovie() {
  editingMovie.value = null
  movieForm.value = { title: '', description: '', genre: '', duration: 60, poster_url: '' }
  showMovieModal.value = true
}

function openEditMovie(m: any) {
  editingMovie.value = m
  movieForm.value = { title: m.title, description: m.description, genre: m.genre, duration: m.duration, poster_url: m.poster_url }
  showMovieModal.value = true
}

async function saveMovie() {
  const url = editingMovie.value ? `${API}/admin/movies/${editingMovie.value.id}` : `${API}/admin/movies`
  const method = editingMovie.value ? 'PUT' : 'POST'
  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(movieForm.value),
  })
  if (!res.ok) throw new Error('Failed to save movie')
  showMovieModal.value = false
  refetch()
}

async function deleteMovie(id: string) {
  if (!confirm('Delete this movie?')) return
  await fetch(`${API}/admin/movies/${id}`, { method: 'DELETE', credentials: 'include' })
  refetch()
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">Movies</h1>

    <Button @click="openCreateMovie" class="mb-4">Add Movie</Button>

    <div v-if="isLoading" class="space-y-2">
      <Card v-for="i in 4" :key="i">
        <CardContent class="space-y-2">
          <Skeleton class="h-4 bg-muted rounded w-1/3" />
          <Skeleton class="h-3 bg-muted rounded w-1/2" />
        </CardContent>
      </Card>
    </div>

    <div v-else-if="isError" class="text-center py-8">
      <p class="text-destructive mb-3">Failed to load movies</p>
      <Button variant="outline" @click="refetch()">Retry</Button>
    </div>

    <div v-else-if="!movies?.length" class="text-muted-foreground">No movies yet</div>

    <Table v-else>
      <TableHeader>
        <TableRow>
          <TableHead>Poster</TableHead>
          <TableHead>Title</TableHead>
          <TableHead>Genre</TableHead>
          <TableHead>Duration</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="m in movies as any[]" :key="m.id">
          <TableCell>
            <img v-if="m.poster_url" :src="m.poster_url" alt="" class="w-10 h-14 object-cover rounded" />
            <div v-else class="w-10 h-14 bg-muted rounded" />
          </TableCell>
          <TableCell class="font-medium">{{ m.title }}</TableCell>
          <TableCell>{{ m.genre }}</TableCell>
          <TableCell>{{ m.duration }} min</TableCell>
          <TableCell>
            <div class="flex gap-2 items-center">
              <Button variant="outline" size="xs" @click="openEditMovie(m)">Edit</Button>
              <Button variant="destructive" size="xs" @click="deleteMovie(m.id)">Delete</Button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <Dialog v-model:open="showMovieModal">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ editingMovie ? 'Edit Movie' : 'Add Movie' }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-3">
          <Input v-model="movieForm.title" placeholder="Title" />
          <Input v-model="movieForm.genre" placeholder="Genre" />
          <Input v-model="movieForm.duration" type="number" placeholder="Duration (min)" />
          <Textarea v-model="movieForm.description" placeholder="Description" />
          <Input v-model="movieForm.poster_url" placeholder="Poster URL" />
        </div>
        <DialogFooter>
          <DialogClose as-child>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button @click="saveMovie">{{ editingMovie ? 'Update' : 'Create' }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
