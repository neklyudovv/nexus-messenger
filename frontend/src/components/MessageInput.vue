<template>
  <div class="px-4 pb-4 pt-1 shrink-0">
    <div class="flex items-end gap-2 bg-overlay rounded-xl px-4 py-3 border border-white/8 focus-within:border-white/15 transition-colors duration-150">
      <textarea
        ref="textareaEl"
        v-model="text"
        :placeholder="placeholder"
        rows="1"
        @keydown.enter.exact.prevent="send"
        @input="onInput"
        class="flex-1 bg-transparent text-tx text-sm placeholder-tx-subtle resize-none focus:outline-none leading-relaxed min-h-[22px] max-h-36 overflow-y-auto"
      />
      <button
        @click="send"
        :disabled="!text.trim()"
        :class="[
          'shrink-0 w-8 h-8 rounded-lg flex items-center justify-center transition-all duration-150',
          text.trim()
            ? 'bg-brand text-white cursor-pointer hover:bg-brand-dark'
            : 'bg-transparent text-tx-subtle cursor-default opacity-40',
        ]"
      >
        <svg class="w-4 h-4 translate-x-px" fill="currentColor" viewBox="0 0 24 24">
          <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
        </svg>
      </button>
    </div>
    <p class="text-tx-subtle text-[11px] mt-1.5 px-1">
      <kbd class="font-sans">Enter</kbd> — отправить &nbsp;·&nbsp; <kbd class="font-sans">Shift+Enter</kbd> — новая строка
    </p>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useChannelsStore } from '../stores/channels'
import { useAuthStore }     from '../stores/auth'
import { useUsersStore }    from '../stores/users'
import { sendMessage, sendTyping } from '../api/ws'

const channels = useChannelsStore()
const auth     = useAuthStore()
const users    = useUsersStore()

const channel  = computed(() => channels.activeChannel)
const isDM     = computed(() => channel.value?.type === 'dm')

const dmPartnerName = computed(() => {
  if (!isDM.value) return ''
  const parts   = channel.value.name.split('-')
  const id1     = Number(parts[1])
  const id2     = Number(parts[2])
  const otherId = id1 === auth.user?.id ? id2 : id1
  return users.getById(otherId)?.username || ''
})

const placeholder = computed(() => {
  if (!channel.value) return 'Выберите канал…'
  return isDM.value
    ? `Написать ${dmPartnerName.value || 'собеседнику'}…`
    : `Написать в #${channel.value.name}…`
})

const text       = ref('')
const textareaEl = ref(null)
let typingTimer  = null

function autoResize() {
  const el = textareaEl.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 144) + 'px'
}

function send() {
  const content = text.value.trim()
  if (!content || !channels.activeChannel) return
  sendMessage(channels.activeChannel.id, content)
  text.value = ''
  nextTick(autoResize)
}

function onInput() {
  autoResize()
  if (!channels.activeChannel) return
  if (typingTimer) clearTimeout(typingTimer)
  sendTyping(channels.activeChannel.id)
  typingTimer = setTimeout(() => { typingTimer = null }, 2000)
}
</script>
