import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/http'

export const useChannelsStore = defineStore('channels', () => {
  const myChannels    = ref([])
  const allDMs        = ref([])
  const activeChannel = ref(null)

  async function fetchMine(workspaceId) {
    const { data } = await api.get(`/workspaces/${workspaceId}/channels`)
    myChannels.value = data || []
  }

  async function fetchAllDMs() {
    const { data } = await api.get('/channels/dms')
    allDMs.value = data || []
  }

  async function create(workspaceId, name, description) {
    const { data } = await api.post(`/workspaces/${workspaceId}/channels`, { name, description, type: 'public' })
    myChannels.value.push(data)
    return data
  }

  async function openDM(workspaceId, userId) {
    const { data } = await api.post(`/workspaces/${workspaceId}/dm/${userId}`)
    if (!allDMs.value.find(c => c.id === data.id)) {
      allDMs.value.push(data)
    }
    return data
  }

  function setActive(channel) {
    activeChannel.value = channel
    if (channel) localStorage.setItem('active_channel_id', channel.id)
    else localStorage.removeItem('active_channel_id')
  }

  function restoreActive() {
    const savedId = Number(localStorage.getItem('active_channel_id'))
    if (!savedId) return null
    return myChannels.value.find(c => c.id === savedId) || null
  }

  function addChannel(ch) {
    if (!myChannels.value.find(c => c.id === ch.id)) {
      myChannels.value.push(ch)
    }
  }

  function reset() {
    myChannels.value    = []
    activeChannel.value = null
    localStorage.removeItem('active_channel_id')
    // allDMs intentionally not reset — DMs are global across workspaces
  }

  return { myChannels, allDMs, activeChannel, fetchMine, fetchAllDMs, create, openDM, addChannel, setActive, restoreActive, reset }
})
