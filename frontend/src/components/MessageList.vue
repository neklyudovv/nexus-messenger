<template>
  <div ref="listEl" class="flex-1 overflow-y-auto px-4 py-4 space-y-1" @scroll="onScroll">
    <div v-if="msgStore.loading" class="text-center text-gray-500 text-sm py-2">Загрузка...</div>
    <MessageItem
      v-for="msg in msgStore.messages"
      :key="msg.id"
      :message="msg"
    />
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { useMessagesStore }   from '../stores/messages'
import { useChannelsStore }   from '../stores/channels'
import MessageItem from './MessageItem.vue'

const msgStore      = useMessagesStore()
const channelsStore = useChannelsStore()
const listEl        = ref(null)

watch(() => channelsStore.activeChannel, async (ch) => {
  if (!ch) return
  msgStore.reset()
  await msgStore.fetchHistory(ch.id)
  await nextTick()
  scrollToBottom()
}, { immediate: true })

watch(() => msgStore.messages.length, async (newLen, oldLen) => {
  if (newLen > oldLen) {
    await nextTick()
    const el = listEl.value
    if (el && el.scrollTop + el.clientHeight >= el.scrollHeight - 100) {
      scrollToBottom()
    }
  }
})

async function onScroll() {
  const el = listEl.value
  if (!el || el.scrollTop > 50) return
  const ch = channelsStore.activeChannel
  if (!ch || !msgStore.hasMore) return
  const oldest = msgStore.messages[0]
  if (oldest) {
    const prevHeight = el.scrollHeight
    await msgStore.fetchHistory(ch.id, oldest.created_at)
    await nextTick()
    el.scrollTop = el.scrollHeight - prevHeight
  }
}

function scrollToBottom() {
  if (listEl.value) listEl.value.scrollTop = listEl.value.scrollHeight
}
</script>
