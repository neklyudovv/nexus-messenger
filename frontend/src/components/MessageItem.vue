<template>
  <div class="flex items-start gap-3 group px-2 py-1 rounded-lg hover:bg-gray-700/40">
    <div class="w-8 h-8 rounded-full bg-indigo-500 flex items-center justify-center text-white text-xs font-bold shrink-0 mt-0.5">
      {{ authorInitial }}
    </div>
    <div class="flex-1 min-w-0">
      <div class="flex items-baseline gap-2">
        <span class="text-white text-sm font-semibold">{{ authorName }}</span>
        <span class="text-gray-500 text-xs">{{ time }}</span>
      </div>
      <p v-if="!message.is_deleted" class="text-gray-200 text-sm break-words">{{ message.content }}</p>
      <p v-else class="text-gray-500 text-sm italic">сообщение удалено</p>
    </div>
    <button
      v-if="isMine && !message.is_deleted"
      @click="remove"
      class="opacity-0 group-hover:opacity-100 text-gray-500 hover:text-red-400 text-xs px-2 py-1 transition"
    >✕</button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { api } from '../api/http'
import { useAuthStore }    from '../stores/auth'
import { useUsersStore }   from '../stores/users'
import { useMessagesStore } from '../stores/messages'

const props = defineProps({ message: Object })

const auth     = useAuthStore()
const users    = useUsersStore()
const msgStore = useMessagesStore()

const author = computed(() => users.getById(props.message.user_id))
const authorName    = computed(() => author.value?.username || `User ${props.message.user_id}`)
const authorInitial = computed(() => authorName.value[0]?.toUpperCase())
const isMine = computed(() => props.message.user_id === auth.user?.id)

const time = computed(() => {
  const d = new Date(props.message.created_at)
  return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
})

async function remove() {
  await api.delete(`/messages/${props.message.id}`)
  msgStore.deleteMessage(props.message.id)
}
</script>
