<template>
  <!-- Grouped: compact, no avatar/name -->
  <div
    v-if="grouped"
    class="group flex items-start gap-3 px-4 py-0.5 hover:bg-white/[0.025] rounded transition-colors"
  >
    <div class="w-10 shrink-0 flex justify-end items-center h-full opacity-0 group-hover:opacity-100 transition-opacity pt-0.5">
      <span class="text-tx-subtle text-[10px] font-medium">{{ time }}</span>
    </div>
    <div class="flex-1 min-w-0 py-0.5">
      <p v-if="!message.is_deleted" class="text-tx text-sm leading-relaxed break-words">{{ message.content }}</p>
      <p v-else class="text-tx-subtle text-sm italic">сообщение удалено</p>
    </div>
    <button
      v-if="isMine && !message.is_deleted"
      @click="remove"
      class="opacity-0 group-hover:opacity-100 text-tx-subtle hover:text-danger transition-all duration-100 cursor-pointer shrink-0 mt-1 p-1 rounded"
      title="Удалить"
    >
      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
      </svg>
    </button>
  </div>

  <!-- First message in group: full, with avatar and name -->
  <div
    v-else
    class="group flex items-start gap-3 px-4 pt-4 pb-0.5 hover:bg-white/[0.025] rounded transition-colors"
  >
    <div
      class="w-10 h-10 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0 mt-0.5 select-none"
      :style="{ background: authorColor }"
    >{{ authorInitial }}</div>

    <div class="flex-1 min-w-0">
      <div class="flex items-baseline gap-2 mb-0.5">
        <span class="text-tx-strong text-sm font-semibold">{{ authorName }}</span>
        <span class="text-tx-subtle text-[11px]">{{ time }}</span>
      </div>
      <p v-if="!message.is_deleted" class="text-tx text-sm leading-relaxed break-words">{{ message.content }}</p>
      <p v-else class="text-tx-subtle text-sm italic">сообщение удалено</p>
    </div>

    <button
      v-if="isMine && !message.is_deleted"
      @click="remove"
      class="opacity-0 group-hover:opacity-100 text-tx-subtle hover:text-danger transition-all duration-100 cursor-pointer shrink-0 mt-1.5 p-1 rounded"
      title="Удалить"
    >
      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
      </svg>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { api } from '../api/http'
import { useAuthStore }     from '../stores/auth'
import { useUsersStore }    from '../stores/users'
import { useMessagesStore } from '../stores/messages'

const props = defineProps({ message: Object, grouped: Boolean })

const auth     = useAuthStore()
const users    = useUsersStore()
const msgStore = useMessagesStore()

const PALETTE = ['#1D6FE8','#0891B2','#0D9488','#16A34A','#D97706','#EA580C','#0E7490','#059669']
function getAvatarColor(str = '') {
  let h = 0
  for (const c of str) h = (h << 5) - h + c.charCodeAt(0)
  return PALETTE[Math.abs(h) % PALETTE.length]
}

const author        = computed(() => users.getById(props.message.user_id))
const authorName    = computed(() => author.value?.username || `User ${props.message.user_id}`)
const authorInitial = computed(() => authorName.value[0]?.toUpperCase())
const authorColor   = computed(() => getAvatarColor(authorName.value))
const isMine        = computed(() => props.message.user_id === auth.user?.id)

const time = computed(() => {
  const d = new Date(props.message.created_at)
  return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
})

async function remove() {
  await api.delete(`/messages/${props.message.id}`)
  msgStore.deleteMessage(props.message.id)
}
</script>
