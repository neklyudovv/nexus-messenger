import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/http'
import { useWorkspaceStore } from './workspace'

export const useMessagesStore = defineStore('messages', () => {
  const messages = ref([])
  const hasMore  = ref(true)
  const loading  = ref(false)

  function reset() {
    messages.value = []
    hasMore.value  = true
    loading.value  = false
  }

  async function fetchHistory(channelId, before = null) {
    if (loading.value || !hasMore.value) return
    const workspaceId = useWorkspaceStore().activeWorkspace?.id
    if (!workspaceId) return
    loading.value = true
    try {
      const { data } = await api.get(`/workspaces/${workspaceId}/channels/${channelId}/messages`, {
        params: before ? { before } : {},
      })
      if (!data || data.length === 0) {
        hasMore.value = false
      } else {
        messages.value = [...data, ...messages.value]
        if (data.length < 50) hasMore.value = false
      }
    } finally {
      loading.value = false
    }
  }

  function addMessage(msg) {
    messages.value.push(msg)
  }

  function deleteMessage(msgId) {
    const idx = messages.value.findIndex(m => m.id === msgId)
    if (idx !== -1) {
      messages.value[idx] = { ...messages.value[idx], is_deleted: true, content: '' }
    }
  }

  return { messages, hasMore, loading, reset, fetchHistory, addMessage, deleteMessage }
})
