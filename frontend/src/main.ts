import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createRouter, createWebHistory } from 'vue-router'
import { Toaster } from 'vue-sonner'
import App from './App.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'movies',
      component: () => import('./views/MovieListView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('./views/LoginView.vue'),
    },
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: () => import('./views/AuthCallbackView.vue'),
    },
    {
      path: '/movies/:id',
      name: 'showtimes',
      component: () => import('./views/ShowtimeListView.vue'),
    },
    {
      path: '/showtimes/:id',
      name: 'seats',
      component: () => import('./views/SeatMapView.vue'),
    },
    {
      path: '/bookings',
      name: 'bookings',
      component: () => import('./views/MyBookingsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/bookings/:id',
      name: 'booking-detail',
      component: () => import('./views/BookingDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('./views/AdminDashboard.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('./views/NotFoundView.vue'),
    },
  ],
})

router.beforeEach((to, from, next) => {
  const stored = localStorage.getItem('auth')
  const authenticated = stored ? !!JSON.parse(stored).refreshToken : false
  if (to.meta.requiresAuth && !authenticated) {
    next('/login')
  } else {
    next()
  }
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin)
app.component('Toaster', Toaster)
app.mount('#app')
