<template>
  <aside class="w-[220px] bg-sidebar flex flex-col h-full shrink-0 border-r border-white/5">

    <!-- Workspace name / header -->
    <div class="relative shrink-0" ref="headerRef">
      <button
        @click="!isDM && (showMenu = !showMenu)"
        :class="[
          'w-full px-4 py-3.5 flex items-center justify-between gap-2 text-left border-b border-white/5 transition-colors',
          isDM ? 'cursor-default' : 'hover:bg-white/4 cursor-pointer',
        ]"
      >
        <span :class="['font-semibold text-sm truncate', isDM ? 'invisible' : 'text-tx-strong']">
          {{ workspace.activeWorkspace?.name }}
        </span>
        <svg
          v-if="!isDM"
          class="w-3.5 h-3.5 text-tx-muted shrink-0 transition-transform duration-150"
          :class="showMenu ? 'rotate-180' : ''"
          fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      <div
        v-if="showMenu && !isDM"
        class="absolute left-2 right-2 top-full mt-1 z-50 bg-surface border border-white/10 rounded-xl shadow-2xl p-1"
      >
        <button
          v-if="isAdmin"
          @click="showInvite = true; showMenu = false"
          class="w-full text-left px-3 py-2 rounded-lg text-sm text-tx hover:bg-overlay transition-colors cursor-pointer flex items-center gap-2.5"
        >
          <svg class="w-4 h-4 text-tx-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
          Invite-код
        </button>
        <p v-if="!isAdmin" class="px-3 py-2 text-xs text-tx-subtle">Нет доступных действий</p>
      </div>
    </div>

    <!-- Channel list -->
    <div class="flex-1 overflow-y-auto py-3 px-2">
      <template v-if="!isDM">
        <div class="flex items-center justify-between px-2 mb-2">
          <span class="text-tx-subtle text-[11px] font-semibold uppercase tracking-widest">Каналы</span>
          <button
            @click="showCreate = true"
            class="w-5 h-5 flex items-center justify-center text-tx-subtle hover:text-tx-muted transition-colors cursor-pointer rounded"
          >
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>

        <button
          v-for="ch in channels.myChannels"
          :key="ch.id"
          @click="select(ch)"
          :class="channelClass(ch)"
        >
          <span class="text-[13px] font-bold mr-1.5 opacity-60">#</span>
          <span class="truncate">{{ ch.name }}</span>
        </button>
      </template>
    </div>
  </aside>

  <!-- Modal: create channel -->
  <div
    v-if="showCreate"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    @click.self="showCreate = false"
  >
    <div class="bg-surface border border-white/10 rounded-2xl p-6 w-80 shadow-2xl">
      <h2 class="text-tx-strong font-semibold text-base mb-5">Новый канал</h2>
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
        <button @click="createChannel" class="flex-1 py-2.5 bg-brand hover:bg-brand-dark text-white rounded-lg text-sm font-medium transition-colors cursor-pointer">Создать</button>
        <button @click="showCreate = false" class="flex-1 py-2.5 bg-overlay hover:bg-white/8 text-tx rounded-lg text-sm transition-colors cursor-pointer">Отмена</button>
      </div>
    </div>
  </div>

  <!-- Modal: invite code -->
  <div
    v-if="showInvite"
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    @click.self="showInvite = false"
  >
    <div class="bg-surface border border-white/10 rounded-2xl p-6 w-80 shadow-2xl">
      <h2 class="text-tx-strong font-semibold text-base mb-5">Invite-код</h2>
      <div class="flex items-center gap-3 bg-overlay rounded-xl px-4 py-3 border border-white/8">
        <span class="text-tx-strong font-mono flex-1 text-sm tracking-wider">{{ workspace.activeWorkspace?.invite_code }}</span>
        <button
          @click="copyCode"
          class="text-brand hover:text-brand-dark text-xs font-semibold transition-colors cursor-pointer shrink-0"
        >{{ copied ? '✓ Готово' : 'Копировать' }}</button>
      </div>
      <button @click="regen" class="w-full py-2.5 bg-overlay hover:bg-white/8 text-tx text-sm rounded-xl transition-colors mt-3 cursor-pointer">
        Сгенерировать новый
      </button>
      <button @click="showInvite = false" class="w-full py-2.5 bg-brand hover:bg-brand-dark text-white text-sm rounded-xl transition-colors mt-2 cursor-pointer font-medium">
        Закрыть
      </button>
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

function channelClass(ch) {
  const isActive = channels.activeChannel?.id === ch.id
  return [
    'w-full text-left px-3 py-1.5 rounded-lg text-sm flex items-center truncate transition-all duration-100 cursor-pointer',
    isActive
      ? 'bg-brand/15 text-tx-strong font-medium'
      : 'text-tx-muted hover:bg-white/5 hover:text-tx',
  ]
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
</script>
