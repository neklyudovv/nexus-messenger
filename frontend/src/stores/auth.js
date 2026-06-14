import { defineStore } from 'pinia'
import { ref } from 'vue'
import { http } from '../api/http'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(localStorage.getItem('access_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  function setAccessToken(token) {
    accessToken.value = token
    localStorage.setItem('access_token', token)
  }

  function setUser(u) {
    user.value = u
    localStorage.setItem('user', JSON.stringify(u))
  }

  async function login(email, password) {
    const { data } = await http.post('/auth/login', { email, password })
    // Refresh token is set as HttpOnly cookie by the server — never touches JS
    setAccessToken(data.access_token)
    setUser(data.user)
  }

  async function register(username, email, password) {
    await http.post('/auth/register', { username, email, password })
    await login(email, password)
  }

  async function logout() {
    // Server reads refresh token from cookie and revokes it
    await http.post('/auth/logout').catch(() => {})
    accessToken.value = ''
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
  }

  return { accessToken, user, login, register, logout }
})
