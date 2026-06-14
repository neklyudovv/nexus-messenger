<template>
  <!-- Rail -->
  <div class="w-[60px] bg-rail flex flex-col items-center py-3 gap-1 shrink-0 border-r border-white/5">

    <!-- Workspace buttons -->
    <div
      v-for="w in workspace.workspaces"
      :key="w.id"
      class="relative flex items-center justify-center w-full h-[52px]"
    >
      <!-- Active indicator pill -->
      <span
        v-if="workspace.activeWorkspace?.id === w.id && channels.activeChannel?.type !== 'dm'"
        class="absolute left-0 top-[14px] h-6 w-[3px] bg-white rounded-r-full pointer-events-none"
      />
      <button
        @click="emit('switch', w)"
        :title="w.name"
        :class="[
          'w-10 h-10 flex items-center justify-center text-sm font-bold transition-all duration-200 cursor-pointer select-none',
          workspace.activeWorkspace?.id === w.id && channels.activeChannel?.type !== 'dm'
            ? 'bg-brand text-white rounded-[12px]'
            : 'bg-surface text-tx-muted rounded-full hover:bg-brand hover:text-white hover:rounded-[12px]',
        ]"
      >{{ w.name[0].toUpperCase() }}</button>
    </div>

    <!-- DM avatars -->
    <template v-if="channels.allDMs.length">
      <div class="w-7 h-px bg-white/10 my-0.5" />
      <div
        v-for="ch in channels.allDMs"
        :key="ch.id"
        class="relative flex items-center justify-center w-full h-[52px]"
      >
        <span
          v-if="channels.activeChannel?.id === ch.id"
          class="absolute left-0 top-[14px] h-6 w-[3px] bg-white rounded-r-full pointer-events-none"
        />
        <button
          @click="selectDM(ch)"
          :title="dmName(ch)"
          :style="{ background: channels.activeChannel?.id === ch.id ? 'var(--color-brand)' : avatarColor(dmName(ch)) }"
          :class="[
            'w-10 h-10 flex items-center justify-center text-white text-xs font-bold transition-all duration-200 cursor-pointer select-none',
            channels.activeChannel?.id === ch.id ? 'rounded-[12px]' : 'rounded-full hover:rounded-[12px]',
          ]"
        >{{ dmName(ch)[0]?.toUpperCase() }}</button>
      </div>
    </template>

    <div class="w-7 h-px bg-white/10 my-0.5" />

    <!-- Add workspace -->
    <div class="relative" ref="addRef">
      <button
        @click="showAdd = !showAdd"
        title="Добавить workspace"
        class="w-10 h-10 rounded-xl bg-surface text-tx-muted hover:bg-online/20 hover:text-online flex items-center justify-center transition-all duration-150 cursor-pointer"
      >
        <svg class="w-4.5 h-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
        </svg>
      </button>

      <div
        v-if="showAdd"
        class="absolute left-full ml-2.5 top-0 z-50 bg-surface border border-white/10 rounded-xl shadow-2xl p-1.5 w-52"
      >
        <button
          @click="showCreate = true; showAdd = false"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2.5"
        >
          <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Создать workspace
        </button>
        <button
          @click="showJoin = true; showAdd = false"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2.5"
        >
          <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          Войти по коду
        </button>
      </div>
    </div>

    <div class="flex-1" />

    <!-- User avatar -->
    <div class="relative" ref="userRef">
      <button
        @click="showUser = !showUser"
        :title="auth.user?.username"
        :style="{ background: avatarColor(auth.user?.username || '') }"
        class="w-10 h-10 rounded-full flex items-center justify-center text-white text-sm font-bold transition-opacity hover:opacity-80 cursor-pointer"
      >{{ auth.user?.username?.[0]?.toUpperCase() }}</button>

      <div
        v-if="showUser"
        class="absolute left-full ml-2.5 bottom-0 z-50 bg-surface border border-white/10 rounded-xl shadow-2xl overflow-hidden w-48"
      >
        <p class="px-3 py-2.5 text-tx-muted text-xs truncate border-b border-white/8">{{ auth.user?.username }}</p>
        <div class="p-1">
          <button
            @click="logout"
            class="w-full text-left px-3 py-2 text-sm text-danger hover:bg-danger/10 rounded-lg transition-colors cursor-pointer flex items-center gap-2"
          >
            <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Выйти
          </button>
        </div>
      </div>
    </div>

  </div>

  <!-- Modal: create workspace -->
  <div
    v-if="showCreate"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    @click.self="showCreate = false"
  >
    <div class="bg-surface border border-white/10 rounded-2xl p-6 w-80 shadow-2xl">
      <h2 class="text-tx-strong font-semibold text-base mb-5">Создать workspace</h2>
      <div class="space-y-4">
        <div>
          <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">Название</label>
          <input v-model="newName" class="w-full px-3.5 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all" />
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">
            Описание
            <span class="normal-case font-normal text-tx-subtle tracking-normal ml-1">необязательно</span>
          </label>
          <input v-model="newDesc" class="w-full px-3.5 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all" />
        </div>
      </div>
      <div class="flex gap-2 mt-5">
        <button @click="createWS" class="flex-1 py-2.5 bg-brand hover:bg-brand-dark text-white rounded-lg text-sm font-medium transition-colors cursor-pointer">Создать</button>
        <button @click="showCreate = false" class="flex-1 py-2.5 bg-overlay hover:bg-white/8 text-tx rounded-lg text-sm transition-colors cursor-pointer">Отмена</button>
      </div>
    </div>
  </div>

  <!-- Modal: join by code -->
  <div
    v-if="showJoin"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    @click.self="showJoin = false"
  >
    <div class="bg-surface border border-white/10 rounded-2xl p-6 w-80 shadow-2xl">
      <h2 class="text-tx-strong font-semibold text-base mb-5">Войти по invite-коду</h2>
      <div>
        <label class="block text-[11px] font-semibold text-tx-muted uppercase tracking-widest mb-2">Код приглашения</label>
        <input v-model="joinCode" class="w-full px-3.5 py-2.5 rounded-lg bg-overlay border border-white/8 text-tx text-sm focus:outline-none focus:ring-1 focus:ring-brand focus:border-brand transition-all" />
      </div>
      <p v-if="joinError" class="text-danger text-sm mt-2">{{ joinError }}</p>
      <div class="flex gap-2 mt-5">
        <button @click="joinWS" class="flex-1 py-2.5 bg-brand hover:bg-brand-dark text-white rounded-lg text-sm font-medium transition-colors cursor-pointer">Войти</button>
        <button @click="showJoin = false" class="flex-1 py-2.5 bg-overlay hover:bg-white/8 text-tx rounded-lg text-sm transition-colors cursor-pointer">Отмена</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter }         from 'vue-router'
import { useWorkspaceStore } from '../stores/workspace'
import { useChannelsStore }  from '../stores/channels'
import { useAuthStore }      from '../stores/auth'
import { useUsersStore }     from '../stores/users'
import { joinChannel as wsJoin } from '../api/ws'

const emit = defineEmits(['switch'])

const workspace = useWorkspaceStore()
const channels  = useChannelsStore()
const auth      = useAuthStore()
const users     = useUsersStore()
const router    = useRouter()

const showAdd    = ref(false)
const showUser   = ref(false)
const showCreate = ref(false)
const showJoin   = ref(false)
const newName    = ref('')
const newDesc    = ref('')
const joinCode   = ref('')
const joinError  = ref('')
const addRef     = ref(null)
const userRef    = ref(null)

const PALETTE = ['#1D6FE8','#0891B2','#0D9488','#16A34A','#D97706','#EA580C','#0E7490','#059669']
function avatarColor(str = '') {
  let h = 0
  for (const c of str) h = (h << 5) - h + c.charCodeAt(0)
  return PALETTE[Math.abs(h) % PALETTE.length]
}

function dmName(ch) {
  const parts = ch.name.split('-')
  const otherId = Number(parts[1]) === auth.user?.id ? Number(parts[2]) : Number(parts[1])
  return users.getById(otherId)?.username || ch.name
}

function selectDM(ch) {
  channels.setActive(ch)
  wsJoin(ch.id)
}

function closeDropdowns(e) {
  if (addRef.value  && !addRef.value.contains(e.target))  showAdd.value  = false
  if (userRef.value && !userRef.value.contains(e.target)) showUser.value = false
}
onMounted(()  => document.addEventListener('mousedown', closeDropdowns))
onUnmounted(() => document.removeEventListener('mousedown', closeDropdowns))

async function createWS() {
  if (!newName.value.trim()) return
  const w = await workspace.create(newName.value.trim(), newDesc.value.trim())
  newName.value = ''
  newDesc.value = ''
  showCreate.value = false
  emit('switch', w)
}

async function joinWS() {
  joinError.value = ''
  try {
    const w = await workspace.join(joinCode.value.trim())
    joinCode.value = ''
    showJoin.value = false
    emit('switch', w)
  } catch {
    joinError.value = 'Неверный код приглашения'
  }
}

async function logout() {
  await auth.logout()
  router.push('/login')
}
</script>
