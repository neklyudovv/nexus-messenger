<template>
  <div class="flex h-screen bg-base overflow-hidden">
    <AppRail @switch="loadWorkspace" />

    <template v-if="workspace.activeWorkspace || activeChannel">
      <Sidebar v-if="!isDM" />

      <main class="flex-1 flex flex-col min-w-0">
        <!-- Channel header -->
        <header class="h-12 px-5 border-b border-white/5 flex items-center gap-2.5 shrink-0 bg-base">
          <template v-if="activeChannel">
            <svg v-if="isDM" class="w-4 h-4 text-tx-muted shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z" clip-rule="evenodd" />
            </svg>
            <span v-else class="text-tx-subtle font-bold text-[15px] leading-none shrink-0">#</span>
            <span class="text-tx-strong font-semibold text-sm">{{ isDM ? dmPartnerName : activeChannel.name }}</span>
            <template v-if="!isDM && activeChannel.description">
              <span class="w-px h-4 bg-white/10 shrink-0 ml-1" />
              <span class="text-tx-subtle text-xs truncate">{{ activeChannel.description }}</span>
            </template>
          </template>
          <span v-else class="text-tx-muted text-sm">Выберите канал</span>
        </header>

        <template v-if="activeChannel">
          <MessageList />
          <MessageInput />
        </template>
        <div v-else class="flex-1 flex items-center justify-center">
          <p class="text-tx-subtle text-sm">Выберите канал или личное сообщение слева</p>
        </div>
      </main>

      <UserList />
    </template>

    <!-- No workspace -->
    <div v-else class="flex-1 flex items-center justify-center">
      <div class="text-center max-w-xs px-6">
        <div class="w-16 h-16 rounded-2xl bg-surface border border-white/8 flex items-center justify-center mx-auto mb-5">
          <svg class="w-8 h-8 text-tx-muted" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <p class="text-tx-strong font-semibold text-base mb-2">Добро пожаловать в Nexus</p>
        <p class="text-tx-muted text-sm leading-relaxed">
          Создайте или вступите в workspace — нажмите
          <span class="text-tx font-semibold">+</span> на левой панели.
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
