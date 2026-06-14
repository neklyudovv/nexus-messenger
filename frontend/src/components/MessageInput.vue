<template>
  <div class="px-4 py-3 border-t border-gray-700">
    <div class="flex items-end gap-2 bg-gray-700 rounded-xl px-4 py-2">
      <textarea
        v-model="text"
        :placeholder="`Написать в #${channel?.name || '...'}`"
        rows="1"
        @keydown.enter.exact.prevent="send"
        @input="onInput"
        class="flex-1 bg-transparent text-white placeholder-gray-400 resize-none focus:outline-none text-sm max-h-32"
      />
      <button
        @click="send"
        :disabled="!text.trim()"
        class="text-indigo-400 hover:text-indigo-300 disabled:opacity-30 transition pb-0.5"
      >
        ➤
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useChannelsStore } from '../stores/channels'
import { sendMessage, sendTyping } from '../api/ws'

const channels = useChannelsStore()
const channel  = computed(() => channels.activeChannel)

const text = ref('')
let typingTimer = null

function send() {
  const content = text.value.trim()
  if (!content || !channels.activeChannel) return
  sendMessage(channels.activeChannel.id, content)
  text.value = ''
}

function onInput() {
  if (!channels.activeChannel) return
  if (typingTimer) clearTimeout(typingTimer)
  sendTyping(channels.activeChannel.id)
  typingTimer = setTimeout(() => { typingTimer = null }, 2000)
}
</script>
