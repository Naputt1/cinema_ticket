import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface User {
  id: string
  email: string
  name: string
  avatar_url: string
  role: 'USER' | 'ADMIN'
}

const API = import.meta.env.VITE_API_URL || '/api/v1'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'ADMIN')

  function loadFromStorage() {
    const stored = localStorage.getItem('auth')
    if (stored) {
      try {
        const data = JSON.parse(stored)
        refreshToken.value = data.refreshToken
        user.value = data.user
      } catch { }
    }
  }

  function saveToStorage() {
    localStorage.setItem('auth', JSON.stringify({
      refreshToken: refreshToken.value,
      user: user.value,
    }))
  }

  function clearStorage() {
    localStorage.removeItem('auth')
  }

  async function checkAuth() {
    try {
      const res = await fetch(`${API}/auth/me`, { credentials: 'include' })
      if (res.ok) {
        user.value = await res.json()
        accessToken.value = 'authenticated'
        saveToStorage()
      }
    } catch { }
  }

  async function fetchUser() {
    try {
      const res = await fetch(`${API}/auth/me`, { credentials: 'include' })
      if (res.ok) {
        user.value = await res.json()
        saveToStorage()
      }
    } catch { }
  }

  async function handleGoogleCallback(code: string) {
    const res = await fetch(`${API}/auth/google/callback?code=${code}`, { credentials: 'include' })
    if (!res.ok) throw new Error('auth failed')
    const data = await res.json()
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    await fetchUser()
    saveToStorage()
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) return false
    try {
      const res = await fetch(`${API}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken.value }),
        credentials: 'include',
      })
      if (!res.ok) throw new Error('refresh failed')
      const data = await res.json()
      accessToken.value = data.access_token
      saveToStorage()
      return true
    } catch {
      logout()
      return false
    }
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    saveToStorage()
  }

  function setUser(u: User) {
    user.value = u
    saveToStorage()
  }

  async function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    clearStorage()
    try {
      await fetch(`${API}/auth/logout`, { credentials: 'include' })
    } catch { }
    window.location.href = '/'
  }

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    isAdmin,
    loadFromStorage,
    checkAuth,
    fetchUser,
    handleGoogleCallback,
    refreshAccessToken,
    setTokens,
    setUser,
    logout,
  }
})
