import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/http'

export const useUsersStore = defineStore('users', () => {
  const users     = ref([])
  const usersById = ref(new Map())
  const online    = ref(new Set())
  const typing    = ref(new Map()) // channelId → Set<userId>

  async function fetchAll(workspaceId) {
    const { data: allUsers } = await api.get('/users')
    const all = allUsers || []

    // Accumulate all users in lookup cache — never shrink it, so DM names resolve
    // correctly even when the other person is not in the current workspace.
    const map = new Map(usersById.value)
    for (const u of all) map.set(u.id, u)
    usersById.value = map

    if (workspaceId) {
      const { data: members } = await api.get(`/workspaces/${workspaceId}/members`)
      const memberMap = new Map((members || []).map(m => [m.user_id, m]))
      users.value = all
        .filter(u => memberMap.has(u.id))
        .map(u => ({ ...u, role: memberMap.get(u.id).role }))
    } else {
      users.value = all
    }
  }

  function setOnline(userId) {
    const s = new Set(online.value)
    s.add(userId)
    online.value = s
  }

  function setOffline(userId) {
    const s = new Set(online.value)
    s.delete(userId)
    online.value = s
  }

  function setTyping(channelId, userId) {
    updateTypingSet(channelId, s => s.add(userId))
    setTimeout(() => clearTyping(channelId, userId), 3000)
  }

  function clearTyping(channelId, userId) {
    updateTypingSet(channelId, s => s.delete(userId))
  }

  function updateTypingSet(channelId, fn) {
    const map = new Map(typing.value)
    const set = new Set(map.get(channelId) || [])
    fn(set)
    map.set(channelId, set)
    typing.value = map
  }

  function isOnline(userId) {
    return online.value.has(userId)
  }

  function getById(userId) {
    return usersById.value.get(userId)
  }

  return { users, online, typing, fetchAll, setOnline, setOffline, setTyping, isOnline, getById }
})
