<template>
  <aside class="w-[196px] bg-sidebar border-l border-white/5 flex flex-col shrink-0 overflow-y-auto">
    <div class="p-3">

      <!-- Online section -->
      <template v-if="onlineUsers.length">
        <p class="px-2 mb-1.5 text-tx-subtle text-[11px] font-semibold uppercase tracking-widest">
          Онлайн · {{ onlineUsers.length }}
        </p>
        <div
          v-for="u in onlineUsers"
          :key="u.id"
          class="flex items-center gap-2.5 px-2 py-1.5 rounded-lg hover:bg-white/5 cursor-pointer transition-colors"
          @click.stop="openProfile(u, $event)"
        >
          <div class="relative shrink-0">
            <div
              class="w-7 h-7 rounded-full flex items-center justify-center text-white text-[11px] font-bold select-none"
              :style="{ background: avatarColor(u.username) }"
            >{{ u.username[0].toUpperCase() }}</div>
            <span class="absolute -bottom-px -right-px w-2.5 h-2.5 rounded-full border-[2px] border-sidebar bg-online" />
          </div>
          <span class="text-tx text-sm truncate">{{ u.username }}</span>
        </div>
      </template>

      <!-- Offline section -->
      <template v-if="offlineUsers.length">
        <p class="px-2 mb-1.5 text-tx-subtle text-[11px] font-semibold uppercase tracking-widest" :class="onlineUsers.length ? 'mt-4' : ''">
          Не в сети · {{ offlineUsers.length }}
        </p>
        <div
          v-for="u in offlineUsers"
          :key="u.id"
          class="flex items-center gap-2.5 px-2 py-1.5 rounded-lg hover:bg-white/5 cursor-pointer transition-colors opacity-40"
          @click.stop="openProfile(u, $event)"
        >
          <div class="relative shrink-0">
            <div
              class="w-7 h-7 rounded-full flex items-center justify-center text-white text-[11px] font-bold grayscale select-none"
              :style="{ background: avatarColor(u.username) }"
            >{{ u.username[0].toUpperCase() }}</div>
            <span class="absolute -bottom-px -right-px w-2.5 h-2.5 rounded-full border-[2px] border-sidebar bg-tx-subtle" />
          </div>
          <span class="text-tx-muted text-sm truncate">{{ u.username }}</span>
        </div>
      </template>

      <!-- Typing indicator -->
      <p v-if="typingNames.length" class="mt-3 px-2 text-tx-muted text-xs italic">
        {{ typingNames.join(', ') }} {{ typingNames.length === 1 ? 'печатает' : 'печатают' }}…
      </p>
    </div>
  </aside>

  <!-- Backdrop -->
  <div v-if="selectedUser" class="fixed inset-0 z-40" @click="selectedUser = null" />

  <!-- Profile popup -->
  <div
    v-if="selectedUser"
    class="fixed z-50 w-52 bg-surface border border-white/12 rounded-2xl shadow-2xl overflow-hidden"
    :style="popupStyle"
  >
    <!-- Header -->
    <div class="px-4 pt-4 pb-3 flex items-center gap-3 border-b border-white/8">
      <div class="relative shrink-0">
        <div
          class="w-11 h-11 rounded-full flex items-center justify-center text-white font-bold text-sm select-none"
          :style="{ background: avatarColor(selectedUser.username) }"
        >{{ selectedUser.username[0].toUpperCase() }}</div>
        <span :class="[
          'absolute -bottom-px -right-px w-3 h-3 rounded-full border-[2px] border-surface',
          users.isOnline(selectedUser.id) ? 'bg-online' : 'bg-tx-subtle',
        ]" />
      </div>
      <div class="min-w-0">
        <p class="text-tx-strong font-semibold text-sm truncate">{{ selectedUser.username }}</p>
        <span :class="roleBadgeClass(selectedUser.role)" class="text-[11px] font-medium px-1.5 py-0.5 rounded mt-0.5 inline-block">
          {{ roleLabel(selectedUser.role) }}
        </span>
      </div>
    </div>

    <!-- Actions -->
    <div class="p-1.5 space-y-0.5">
      <template v-if="isAdmin && selectedUser.id !== auth.user?.id">
        <button
          v-if="selectedUser.role !== 'admin'"
          @click="changeRole('admin')"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2"
        >
          <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" />
          </svg>
          Сделать администратором
        </button>
        <button
          v-if="selectedUser.role !== 'member'"
          @click="changeRole('member')"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2"
        >
          <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
          Сделать участником
        </button>
      </template>

      <button
        v-if="selectedUser.id !== auth.user?.id && !isDM"
        @click="startDM()"
        class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2"
      >
        <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
        </svg>
        Написать сообщение
      </button>

      <p v-if="selectedUser.id === auth.user?.id && !isAdmin" class="px-3 py-2 text-xs text-tx-subtle">Это вы</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUsersStore }     from '../stores/users'
import { useChannelsStore }  from '../stores/channels'
import { useWorkspaceStore } from '../stores/workspace'
import { useAuthStore }      from '../stores/auth'
import { joinChannel }       from '../api/ws'

const users     = useUsersStore()
const channels  = useChannelsStore()
const workspace = useWorkspaceStore()
const auth      = useAuthStore()

const selectedUser = ref(null)
const popupStyle   = ref({})

const isDM = computed(() => channels.activeChannel?.type === 'dm')

const visibleUsers = computed(() => {
  if (!isDM.value) return users.users
  const parts = channels.activeChannel.name.split('-')
  const id1   = Number(parts[1])
  const id2   = Number(parts[2])
  return [id1, id2].map(id => users.getById(id)).filter(Boolean)
})

const onlineUsers  = computed(() => visibleUsers.value.filter(u =>  users.isOnline(u.id)))
const offlineUsers = computed(() => visibleUsers.value.filter(u => !users.isOnline(u.id)))

const isAdmin = computed(() => {
  const me = users.users.find(u => u.id === auth.user?.id)
  return me ? me.role === 'admin' : workspace.activeWorkspace?.owner_id === auth.user?.id
})

const typingNames = computed(() => {
  const ch = channels.activeChannel
  if (!ch) return []
  const typingSet = users.typing.get(ch.id) || new Set()
  return [...typingSet]
    .filter(id => id !== auth.user?.id)
    .map(id => users.getById(id)?.username || `User ${id}`)
})

const PALETTE = ['#1D6FE8','#0891B2','#0D9488','#16A34A','#D97706','#EA580C','#0E7490','#059669']
function avatarColor(str = '') {
  let h = 0
  for (const c of str) h = (h << 5) - h + c.charCodeAt(0)
  return PALETTE[Math.abs(h) % PALETTE.length]
}

function openProfile(u, event) {
  selectedUser.value = u
  const rect        = event.currentTarget.getBoundingClientRect()
  const popupHeight = 220
  const top         = Math.min(rect.top, window.innerHeight - popupHeight - 8)
  popupStyle.value  = {
    top:   top + 'px',
    right: (window.innerWidth - rect.left + 8) + 'px',
  }
}

function roleLabel(role) {
  return { admin: 'Администратор', member: 'Участник' }[role] ?? role
}

function roleBadgeClass(role) {
  return {
    admin:  'bg-brand/20 text-brand',
    member: 'bg-white/8 text-tx-muted',
  }[role] ?? 'bg-white/8 text-tx-muted'
}

async function changeRole(newRole) {
  if (!workspace.activeWorkspace || !selectedUser.value) return
  await workspace.updateMemberRole(workspace.activeWorkspace.id, selectedUser.value.id, newRole)
  const u = users.users.find(u => u.id === selectedUser.value.id)
  if (u) u.role = newRole
  selectedUser.value = { ...selectedUser.value, role: newRole }
}

async function startDM() {
  if (!workspace.activeWorkspace || !selectedUser.value) return
  const ch = await channels.openDM(workspace.activeWorkspace.id, selectedUser.value.id)
  channels.setActive(ch)
  joinChannel(ch.id)
  selectedUser.value = null
}
</script>
