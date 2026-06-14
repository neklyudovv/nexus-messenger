<template>
  <aside class="w-56 bg-gray-800 flex flex-col h-full shrink-0 border-r border-gray-700/40">
    <!-- workspace name — click opens dropdown -->
    <div class="relative shrink-0" ref="headerRef">
      <button
        @click="!isDM && (showMenu = !showMenu)"
        :class="['w-full px-4 py-3 border-b border-gray-700 flex items-center justify-between gap-2 transition text-left', isDM ? 'cursor-default' : 'hover:bg-gray-700/50']"
      >
        <span :class="['font-semibold truncate', isDM ? 'invisible' : 'text-white']">
          {{ workspace.activeWorkspace?.name }}
        </span>
        <svg v-if="!isDM" class="w-3.5 h-3.5 text-gray-400 shrink-0 transition-transform" :class="showMenu ? 'rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" /></svg>
      </button>

      <div
        v-if="showMenu && !isDM"
        class="absolute left-2 right-2 top-full mt-1 z-50 bg-gray-800 border border-gray-700 rounded-xl shadow-xl p-1.5"
      >
        <button
          v-if="isAdmin"
          @click="showInvite = true; showMenu = false"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition flex items-center gap-2"
        >
          <span>🔗</span> Invite-код
        </button>
        <p v-if="!isAdmin" class="px-3 py-2 text-xs text-gray-500">Нет доступных действий</p>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto px-2 py-3">
      <template v-if="!isDM">
        <div class="flex items-center justify-between px-2 mb-1">
          <span class="text-gray-400 text-xs uppercase font-semibold tracking-wide">Каналы</span>
          <button @click="showCreate = true" class="text-gray-400 hover:text-white text-lg leading-none">+</button>
        </div>

        <button
          v-for="ch in channels.myChannels"
          :key="ch.id"
          @click="select(ch)"
          :class="channelClass(ch)"
        ># {{ ch.name }}</button>
      </template>
    </div>
  </aside>

  <!-- Modal: create channel -->
  <div v-if="showCreate" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showCreate = false">
    <div class="bg-gray-800 rounded-xl p-6 w-80 space-y-4">
      <h2 class="text-white font-semibold text-lg">Новый канал</h2>
      <input v-model="newName" placeholder="Название" class="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" />
      <input v-model="newDesc" placeholder="Описание (необязательно)" class="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" />
      <div class="flex gap-2">
        <button @click="createChannel" class="flex-1 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition">Создать</button>
        <button @click="showCreate = false" class="flex-1 py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg transition">Отмена</button>
      </div>
    </div>
  </div>

  <!-- Modal: invite code -->
  <div v-if="showInvite" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showInvite = false">
    <div class="bg-gray-800 rounded-xl p-6 w-80 space-y-4">
      <h2 class="text-white font-semibold text-lg">Invite-код</h2>
      <div class="flex items-center gap-2 bg-gray-700 rounded-lg px-3 py-2">
        <span class="text-white font-mono flex-1 text-sm">{{ workspace.activeWorkspace?.invite_code }}</span>
        <button @click="copyCode" class="text-indigo-400 hover:text-indigo-300 text-sm">{{ copied ? '✓' : 'Копировать' }}</button>
      </div>
      <button @click="regen" class="w-full py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg text-sm transition">Обновить код</button>
      <button @click="showInvite = false" class="w-full py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition">Закрыть</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useChannelsStore }  from '../stores/channels'
import { useWorkspaceStore } from '../stores/workspace'
import { useUsersStore }     from '../stores/users'
import { useAuthStore }      from '../stores/auth'
import { joinChannel as wsJoin } from '../api/ws'

const channels  = useChannelsStore()
const workspace = useWorkspaceStore()
const users     = useUsersStore()
const auth      = useAuthStore()

const isDM    = computed(() => channels.activeChannel?.type === 'dm')
const isAdmin = computed(() => {
  const me = users.users.find(u => u.id === auth.user?.id)
  return me ? me.role === 'admin' : workspace.activeWorkspace?.owner_id === auth.user?.id
})

const showMenu   = ref(false)
const showCreate = ref(false)
const showInvite = ref(false)
const copied     = ref(false)
const newName    = ref('')
const newDesc    = ref('')
const headerRef  = ref(null)

function closeMenu(e) {
  if (headerRef.value && !headerRef.value.contains(e.target)) showMenu.value = false
}
onMounted(()  => document.addEventListener('mousedown', closeMenu))
onUnmounted(() => document.removeEventListener('mousedown', closeMenu))

function select(ch) {
  channels.setActive(ch)
  wsJoin(ch.id)
}

async function createChannel() {
  if (!newName.value.trim() || !workspace.activeWorkspace) return
  const ch = await channels.create(
    workspace.activeWorkspace.id,
    newName.value.trim(),
    newDesc.value.trim(),
  )
  showCreate.value = false
  newName.value = ''
  newDesc.value = ''
  select(ch)
}

async function copyCode() {
  await navigator.clipboard.writeText(workspace.activeWorkspace?.invite_code || '')
  copied.value = true
  setTimeout(() => copied.value = false, 2000)
}

async function regen() {
  if (workspace.activeWorkspace) await workspace.regenerateInvite(workspace.activeWorkspace.id)
}

function channelClass(ch) {
  return [
    'w-full text-left px-3 py-1 rounded-md text-sm truncate transition',
    channels.activeChannel?.id === ch.id
      ? 'bg-indigo-600 text-white'
      : 'text-gray-300 hover:bg-gray-700',
  ]
}
</script>
