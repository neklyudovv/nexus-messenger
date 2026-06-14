<template>
  <div class="flex h-screen bg-gray-900 overflow-hidden">
    <AppRail @switch="loadWorkspace" />

    <template v-if="workspace.activeWorkspace || activeChannel">
      <Sidebar />

      <main class="flex-1 flex flex-col min-w-0">
        <div class="px-4 py-3 border-b border-gray-700 flex items-center gap-2 shrink-0">
          <template v-if="activeChannel">
            <span class="text-gray-400 text-lg leading-none">{{ isDM ? '💬' : '#' }}</span>
            <span class="text-white font-semibold">{{ isDM ? dmPartnerName : activeChannel.name }}</span>
            <span v-if="!isDM && activeChannel.description" class="text-gray-500 text-sm ml-2">{{ activeChannel.description }}</span>
          </template>
          <span v-else class="text-gray-500">Выберите канал</span>
        </div>

        <template v-if="activeChannel">
          <MessageList />
          <MessageInput />
        </template>
        <div v-else class="flex-1 flex items-center justify-center">
          <p class="text-gray-600 text-sm">Выберите канал или личное сообщение слева</p>
        </div>
      </main>

      <UserList />
    </template>

    <div v-else class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-3 max-w-xs px-6">
        <p class="text-5xl">💬</p>
        <p class="text-white font-semibold text-lg">Добро пожаловать!</p>
        <p class="text-gray-500 text-sm leading-relaxed">
          Создайте workspace или вступите по invite-коду — нажмите
          <span class="text-gray-300 font-bold">+</span> на левой панели.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import AppRail      from '../components/AppRail.vue'
import Sidebar      from '../components/Sidebar.vue'
import MessageList  from '../components/MessageList.vue'
import MessageInput from '../components/MessageInput.vue'
import UserList     from '../components/UserList.vue'
import { useChannelsStore }  from '../stores/channels'
import { useWorkspaceStore } from '../stores/workspace'
import { useUsersStore }     from '../stores/users'
import { useAuthStore }      from '../stores/auth'
import { connectWS, disconnectWS, joinChannel } from '../api/ws'

const channels  = useChannelsStore()
const workspace = useWorkspaceStore()
const users     = useUsersStore()
const auth      = useAuthStore()

const activeChannel = computed(() => channels.activeChannel)
const isDM          = computed(() => activeChannel.value?.type === 'dm')

const dmPartnerName = computed(() => {
  if (!isDM.value) return ''
  const parts   = activeChannel.value.name.split('-')
  const id1     = Number(parts[1])
  const id2     = Number(parts[2])
  const otherId = id1 === auth.user?.id ? id2 : id1
  return users.getById(otherId)?.username || activeChannel.value.name
})

onMounted(async () => {
  await workspace.fetchMine()
  await channels.fetchAllDMs()
  const saved = workspace.restoreActive()
  if (saved) await loadWorkspace(saved)
  connectWS()
})

onUnmounted(() => disconnectWS())

async function loadWorkspace(ws) {
  workspace.setActive(ws)
  channels.reset()
  await channels.fetchMine(ws.id)
  await users.fetchAll(ws.id)
  const saved = channels.restoreActive()
  if (saved) {
    channels.setActive(saved)
    joinChannel(saved.id)
  }
}
</script>
