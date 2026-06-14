import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/http'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces       = ref([])
  const activeWorkspace  = ref(null)

  async function fetchMine() {
    const { data } = await api.get('/workspaces/me')
    workspaces.value = data || []
  }

  async function create(name, description) {
    const { data } = await api.post('/workspaces', { name, description })
    workspaces.value.push(data)
    return data
  }

  async function join(inviteCode) {
    const { data } = await api.post('/workspaces/join', { invite_code: inviteCode })
    if (!workspaces.value.find(w => w.id === data.id)) {
      workspaces.value.push(data)
    }
    return data
  }

  async function getInviteCode(workspaceId) {
    const ws = workspaces.value.find(w => w.id === workspaceId)
    return ws?.invite_code || null
  }

  async function updateMemberRole(workspaceId, userId, role) {
    await api.patch(`/workspaces/${workspaceId}/members/${userId}`, { role })
  }

  async function regenerateInvite(workspaceId) {
    const { data } = await api.post(`/workspaces/${workspaceId}/invite/regenerate`)
    const ws = workspaces.value.find(w => w.id === workspaceId)
    if (ws) ws.invite_code = data.invite_code
    return data.invite_code
  }

  function setActive(ws) {
    activeWorkspace.value = ws
    if (ws) localStorage.setItem('active_workspace_id', ws.id)
    else localStorage.removeItem('active_workspace_id')
  }

  function restoreActive() {
    const savedId = Number(localStorage.getItem('active_workspace_id'))
    if (!savedId) return null
    return workspaces.value.find(w => w.id === savedId) || null
  }

  return { workspaces, activeWorkspace, fetchMine, create, join, getInviteCode, regenerateInvite, updateMemberRole, setActive, restoreActive }
})
