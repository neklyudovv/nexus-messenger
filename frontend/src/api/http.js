import axios from 'axios'

export const http = axios.create({ baseURL: '/api' })
export const api  = axios.create({ baseURL: '/api' })

api.interceptors.request.use(cfg => {
  const token = localStorage.getItem('access_token')
  if (token) cfg.headers.Authorization = `Bearer ${token}`
  return cfg
})

let refreshing = false
let queue = []

api.interceptors.response.use(
  res => res,
  async err => {
    const original = err.config
    if (err.response?.status !== 401 || original._retry) {
      return Promise.reject(err)
    }
    if (refreshing) {
      return new Promise((resolve, reject) => {
        queue.push({ resolve, reject, config: original })
      })
    }
    refreshing = true
    original._retry = true
    try {
      // Refresh token is in an HttpOnly cookie — no body needed
      const { data } = await http.post('/auth/refresh')
      localStorage.setItem('access_token', data.access_token)
      api.defaults.headers.common.Authorization = `Bearer ${data.access_token}`
      queue.forEach(({ resolve, config }) => {
        config.headers.Authorization = `Bearer ${data.access_token}`
        resolve(api(config))
      })
      queue = []
      return api(original)
    } catch {
      queue.forEach(({ reject }) => reject(err))
      queue = []
      localStorage.removeItem('access_token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    } finally {
      refreshing = false
    }
  }
)
