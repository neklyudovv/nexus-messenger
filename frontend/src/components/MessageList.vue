<template>
  <div ref="listEl" class="flex-1 overflow-y-auto py-2" @scroll="onScroll">
    <div v-if="msgStore.loading" class="text-center text-tx-muted text-xs py-4">Загрузка…</div>
    <MessageItem
      v-for="(msg, i) in msgStore.messages"
      :key="msg.id"
      :message="msg"
      :grouped="isGrouped(msg, i)"
    />
    <div class="h-2" />
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { useMessagesStore } from '../stores/messages'
import { useChannelsStore } from '../stores/channels'
import MessageItem from './MessageItem.vue'

const msgStore      = useMessagesStore()
const channelsStore = useChannelsStore()
const listEl        = ref(null)

function isGrouped(msg, i) {
  if (i === 0) return false
  const prev = msgStore.messages[i - 1]
  if (prev.user_id !== msg.user_id || prev.is_deleted || msg.is_deleted) return false
  const gap = new Date(msg.created_at).getTime() - new Date(prev.created_at).getTime()
  return gap < 5 * 60 * 1000
}

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
    if (el && el.scrollTop + el.clientHeight >= el.scrollHeight - 120) scrollToBottom()
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
