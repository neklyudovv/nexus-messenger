<template>
  <div class="w-16 bg-gray-900 flex flex-col items-center py-3 gap-1.5 shrink-0 border-r border-gray-700/40">
    <!-- workspace icons -->
    <button
      v-for="w in workspace.workspaces"
      :key="w.id"
      @click="emit('switch', w)"
      :title="w.name"
      :class="[
        'w-10 h-10 flex items-center justify-center text-white text-sm font-bold transition-all duration-150 shrink-0',
        workspace.activeWorkspace?.id === w.id && channels.activeChannel?.type !== 'dm'
          ? 'bg-indigo-600 rounded-2xl'
          : 'bg-gray-700 rounded-xl hover:bg-indigo-500 hover:rounded-2xl',
      ]"
    >{{ w.name[0].toUpperCase() }}</button>

    <!-- separator + DMs (global, not workspace-scoped) -->
    <template v-if="channels.allDMs.length">
      <div class="w-8 border-t border-gray-600 my-0.5" />
      <button
        v-for="ch in channels.allDMs"
        :key="ch.id"
        @click="selectDM(ch)"
        :title="dmName(ch)"
        :class="[
          'w-10 h-10 flex items-center justify-center text-white text-xs font-bold transition-all duration-150 shrink-0',
          channels.activeChannel?.id === ch.id
            ? 'bg-indigo-600 rounded-2xl'
            : 'bg-gray-600 rounded-full hover:bg-indigo-500 hover:rounded-2xl',
        ]"
      >{{ dmInitial(ch) }}</button>
    </template>

    <div class="w-8 border-t border-gray-600 my-0.5" />

    <!-- add workspace -->
    <div class="relative" ref="addRef">
      <button
        @click="showAdd = !showAdd"
        title="Добавить workspace"
        class="w-10 h-10 rounded-xl bg-gray-700 hover:bg-green-600 hover:rounded-2xl flex items-center justify-center text-gray-300 hover:text-white text-2xl font-light transition-all duration-150"
      >+</button>

      <div
        v-if="showAdd"
        class="absolute left-full ml-2 top-0 z-50 bg-gray-800 border border-gray-700 rounded-xl shadow-xl p-2 w-52"
      >
        <button @click="showCreate = true; showAdd = false" class="w-full text-left px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition">
          + Создать workspace
        </button>
        <button @click="showJoin = true; showAdd = false" class="w-full text-left px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition">
          → Войти по коду
        </button>
      </div>
    </div>

    <div class="flex-1" />

    <!-- user avatar -->
    <div class="relative" ref="userRef">
      <button
        @click="showUser = !showUser"
        :title="auth.user?.username"
        class="w-10 h-10 rounded-full bg-indigo-500 hover:bg-indigo-400 flex items-center justify-center text-white text-sm font-bold transition"
      >{{ auth.user?.username?.[0]?.toUpperCase() }}</button>

      <div
        v-if="showUser"
        class="absolute left-full ml-2 bottom-0 z-50 bg-gray-800 border border-gray-700 rounded-xl shadow-xl p-2 w-44"
      >
        <p class="px-3 py-1.5 text-gray-400 text-sm truncate">{{ auth.user?.username }}</p>
        <hr class="border-gray-700 my-1" />
        <button @click="logout" class="w-full text-left px-3 py-1.5 text-sm text-gray-300 hover:text-white hover:bg-gray-700 rounded-lg transition">
          Выйти
        </button>
      </div>
    </div>
  </div>

  <!-- Modal: create workspace -->
  <div v-if="showCreate" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50" @click.self="showCreate = false">
    <div class="bg-gray-800 rounded-xl p-6 w-80 space-y-4">
      <h2 class="text-white font-semibold text-lg">Создать workspace</h2>
      <input v-model="newName" placeholder="Название" class="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" />
      <input v-model="newDesc" placeholder="Описание (необязательно)" class="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" />
      <div class="flex gap-2">
        <button @click="createWS" class="flex-1 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition">Создать</button>
        <button @click="showCreate = false" class="flex-1 py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg transition">Отмена</button>
      </div>
    </div>
  </div>

  <!-- Modal: join by code -->
  <div v-if="showJoin" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50" @click.self="showJoin = false">
    <div class="bg-gray-800 rounded-xl p-6 w-80 space-y-4">
      <h2 class="text-white font-semibold text-lg">Войти по invite-коду</h2>
      <input v-model="joinCode" placeholder="Код приглашения" class="w-full px-3 py-2 bg-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500" />
      <p v-if="joinError" class="text-red-400 text-sm">{{ joinError }}</p>
      <div class="flex gap-2">
        <button @click="joinWS" class="flex-1 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition">Войти</button>
        <button @click="showJoin = false" class="flex-1 py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg transition">Отмена</button>
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

function dmName(ch) {
  const parts = ch.name.split('-')
  const otherId = Number(parts[1]) === auth.user?.id ? Number(parts[2]) : Number(parts[1])
  return users.getById(otherId)?.username || ch.name
}

function dmInitial(ch) {
  return dmName(ch)[0]?.toUpperCase() || '?'
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
