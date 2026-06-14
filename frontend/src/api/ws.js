import { useMessagesStore }  from '../stores/messages'
import { useUsersStore }     from '../stores/users'
import { useChannelsStore }  from '../stores/channels'
import { useWorkspaceStore } from '../stores/workspace'

let socket = null
let shouldReconnect = false

export function connectWS() {
  const token = localStorage.getItem('access_token')
  if (!token) return
  if (socket?.readyState === WebSocket.OPEN) return

  shouldReconnect = true
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  socket = new WebSocket(`${proto}://${location.host}/ws`)

  socket.onopen = () => {
    // токен передаётся первым сообщением, а не в URL (защита от утечки в логи)
    socket.send(JSON.stringify({ type: 'auth', token }))
  }

  socket.onmessage = ({ data }) => {
    const ev = JSON.parse(data)
    const msgStore      = useMessagesStore()
    const usersStore    = useUsersStore()
    const channelStore  = useChannelsStore()
    const workspaceStore = useWorkspaceStore()

    switch (ev.type) {
      case 'new_message':
        msgStore.addMessage(ev.payload)
        break
      case 'message_deleted':
        msgStore.deleteMessage(ev.payload?.id)
        break
      case 'user_online':
        usersStore.setOnline(ev.user_id)
        break
      case 'user_offline':
        usersStore.setOffline(ev.user_id)
        break
      case 'typing':
        usersStore.setTyping(ev.channel_id, ev.user_id)
        break
      case 'channel_created':
        if (ev.payload?.workspace_id === workspaceStore.activeWorkspace?.id) {
          channelStore.addChannel(ev.payload)
        }
        break
    }
  }

  socket.onclose = () => {
    socket = null
    if (shouldReconnect) setTimeout(connectWS, 3000)
  }
}

export function disconnectWS() {
  shouldReconnect = false
  socket?.close()
  socket = null
}

export function sendWS(payload) {
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(payload))
  }
}

export function joinChannel(channelId) {
  sendWS({ type: 'join_channel', channel_id: channelId })
}

export function sendMessage(channelId, content) {
  sendWS({ type: 'send_message', channel_id: channelId, content })
}

export function sendTyping(channelId) {
  sendWS({ type: 'typing', channel_id: channelId })
}
