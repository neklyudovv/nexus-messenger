<template>
  <aside class="w-48 bg-gray-800 border-l border-gray-700 p-3 shrink-0 overflow-y-auto">
    <p class="text-gray-400 text-xs uppercase font-semibold tracking-wide mb-2">Участники</p>
    <div class="space-y-1">
      <div
        v-for="u in visibleUsers"
        :key="u.id"
        class="flex items-center gap-2 px-2 py-1 rounded-md hover:bg-gray-700 cursor-pointer transition"
        @click.stop="openProfile(u, $event)"
      >
        <div class="relative shrink-0">
          <div class="w-7 h-7 rounded-full bg-indigo-500 flex items-center justify-center text-white text-xs font-bold">
            {{ u.username[0].toUpperCase() }}
          </div>
          <span :class="['absolute bottom-0 right-0 w-2.5 h-2.5 rounded-full border-2 border-gray-800', users.isOnline(u.id) ? 'bg-green-400' : 'bg-gray-500']" />
        </div>
        <span class="text-gray-300 text-sm truncate">{{ u.username }}</span>
      </div>
    </div>

    <div v-if="typingNames.length" class="mt-3 text-gray-400 text-xs italic">
      {{ typingNames.join(', ') }} {{ typingNames.length === 1 ? 'печатает' : 'печатают' }}...
    </div>
  </aside>

  <!-- click-outside backdrop -->
  <div v-if="selectedUser" class="fixed inset-0 z-40" @click="selectedUser = null" />

  <!-- profile popup -->
  <div
    v-if="selectedUser"
    class="fixed z-50 w-56 bg-gray-800 border border-gray-700 rounded-xl shadow-xl overflow-hidden"
    :style="popupStyle"
  >
    <!-- header -->
    <div class="px-4 pt-4 pb-3 flex items-center gap-3">
      <div class="relative shrink-0">
        <div class="w-11 h-11 rounded-full bg-indigo-500 flex items-center justify-center text-white font-bold text-base">
          {{ selectedUser.username[0].toUpperCase() }}
        </div>
        <span :class="['absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-gray-800', users.isOnline(selectedUser.id) ? 'bg-green-400' : 'bg-gray-500']" />
      </div>
      <div class="min-w-0">
        <p class="text-white font-semibold text-sm truncate">{{ selectedUser.username }}</p>
        <span :class="roleBadgeClass(selectedUser.role)" class="text-xs font-medium px-1.5 py-0.5 rounded mt-0.5 inline-block">
          {{ roleLabel(selectedUser.role) }}
        </span>
      </div>
    </div>

    <!-- actions -->
    <div class="border-t border-gray-700 p-1.5 space-y-0.5">
      <!-- role change (admin only, not self) -->
      <template v-if="isAdmin && selectedUser.id !== auth.user?.id">
        <button
          v-if="selectedUser.role !== 'admin'"
          @click="changeRole('admin')"
          class="w-full text-left px-3 py-1.5 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition"
        >↑ Сделать администратором</button>
        <button
          v-if="selectedUser.role !== 'member'"
          @click="changeRole('member')"
          class="w-full text-left px-3 py-1.5 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition"
        >↓ Сделать участником</button>
      </template>

      <!-- DM (not self, not already in DM with them) -->
      <button
        v-if="selectedUser.id !== auth.user?.id && !isDM"
        @click="startDM()"
        class="w-full text-left px-3 py-1.5 rounded-lg text-sm text-gray-300 hover:bg-gray-700 transition"
      >💬 Написать сообщение</button>

      <!-- self — no actions hint -->
      <p v-if="selectedUser.id === auth.user?.id && !isAdmin" class="px-3 py-1.5 text-xs text-gray-500">Это вы</p>
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
  const parts   = channels.activeChannel.name.split('-')
  const id1     = Number(parts[1])
  const id2     = Number(parts[2])
  return [id1, id2].map(id => users.getById(id)).filter(Boolean)
})

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

function openProfile(u, event) {
  selectedUser.value = u
  const rect = event.currentTarget.getBoundingClientRect()
  const popupHeight = 180
  const top = Math.min(rect.top, window.innerHeight - popupHeight - 8)
  popupStyle.value = {
    top:   top + 'px',
    right: (window.innerWidth - rect.left + 8) + 'px',
  }
}

function roleLabel(role) {
  return { admin: 'Администратор', member: 'Участник' }[role] ?? role
}

function roleBadgeClass(role) {
  return {
    admin:  'bg-indigo-500/20 text-indigo-400',
    member: 'bg-gray-500/20 text-gray-400',
  }[role] ?? 'bg-gray-500/20 text-gray-400'
}

async function changeRole(newRole) {
  if (!workspace.activeWorkspace || !selectedUser.value) return
  await workspace.updateMemberRole(workspace.activeWorkspace.id, selectedUser.value.id, newRole)
  // Update locally so the badge refreshes without a full refetch
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
