<template>
  <div class="min-h-screen bg-rail flex items-center justify-center p-6">
    <div class="w-full max-w-sm">

      <!-- Brand mark -->
      <div class="flex items-center justify-center gap-3 mb-8">
        <div class="w-10 h-10 rounded-xl bg-brand flex items-center justify-center shadow-lg">
          <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <span class="text-tx-strong text-xl font-semibold tracking-tight">Nexus</span>
      </div>

      <!-- Card -->
      <div class="bg-surface border border-white/8 rounded-2xl p-8 shadow-2xl">
        <h1 class="text-tx-strong text-xl font-semibold mb-1">Создать аккаунт</h1>
        <p class="text-tx-muted text-sm mb-7">Заполните данные для регистрации</p>

        <form @submit.prevent="submit" class="space-y-5">
          <div>
            <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">Имя пользователя</label>
            <input
              v-model="username"
              type="text"
              autocomplete="username"
              required
              class="w-full px-4 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all duration-150"
            />
          </div>
          <div>
            <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">Email</label>
            <input
              v-model="email"
              type="email"
              autocomplete="email"
              required
              class="w-full px-4 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all duration-150"
            />
          </div>
          <div>
            <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">
              Пароль
              <span class="normal-case font-normal text-tx-subtle tracking-normal ml-1">· мин. 8 символов</span>
            </label>
            <input
              v-model="password"
              type="password"
              autocomplete="new-password"
              required
              class="w-full px-4 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all duration-150"
            />
          </div>

          <p v-if="error" class="text-danger text-sm">{{ error }}</p>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 rounded-lg bg-brand hover:bg-brand-dark text-white text-sm font-semibold transition-colors duration-150 disabled:opacity-50 cursor-pointer"
          >
            {{ loading ? 'Создаём…' : 'Создать аккаунт' }}
          </button>
        </form>
      </div>

      <p class="text-tx-muted text-sm text-center mt-5">
        Уже есть аккаунт?
        <RouterLink to="/login" class="text-brand hover:text-brand-dark font-medium transition-colors">Войти</RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth   = useAuthStore()
const router = useRouter()

const username = ref('')
const email    = ref('')
const password = ref('')
const error    = ref('')
const loading  = ref(false)

async function submit() {
  error.value   = ''
  loading.value = true
  try {
    await auth.register(username.value, email.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>
