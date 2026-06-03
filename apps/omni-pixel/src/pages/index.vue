<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useConversations } from '@/composables/useConversations'
import { getConversation, streamChat, type ChatMessage } from '@/lib/conversation'
import HigSigninModal from '@/components/HigSigninModal.vue'

const { isAuthenticated } = useAuth()
const showSignin = computed(() => !isAuthenticated.value)

const {
    conversations: conversationList,
    activeId,
    loading: conversationsLoading,
    error: conversationsError,
    fetchConversations,
    create,
    remove,
    setActive,
} = useConversations()

const messages = ref<ChatMessage[]>([])
const isLoading = ref(false)
const streamBuffer = ref('')
const messagesEl = ref<HTMLElement | null>(null)

let abortController: AbortController | null = null

onMounted(() => {
    fetchConversations()
})

async function scrollToBottom() {
    await nextTick()
    if (messagesEl.value) {
        messagesEl.value.scrollTop = messagesEl.value.scrollHeight
    }
}

async function handleSelect(id: string) {
    setActive(id)
    try {
        const detail = await getConversation(id)
        messages.value = detail.messages ?? []
    } catch {
        messages.value = []
    }
}

async function handleNew() {
    const session = await create()
    if (session) {
        setActive('')
        messages.value = []
    }
}

async function handleDelete(id: string) {
    await remove(id)
}

async function handleSend(payload: { message: string, model: { id: string } }) {
    const convId = activeId.value || undefined
    const userMsg: ChatMessage = {
        id: `user-${Date.now()}`,
        conversation_id: convId ?? '',
        user_id: '',
        content: payload.message,
        model_id: payload.model.id,
        type: 0,
        created_at: new Date().toISOString(),
    }
    messages.value.push(userMsg)
    await scrollToBottom()

    isLoading.value = true
    streamBuffer.value = ''
    abortController = new AbortController()

    try {
        await streamChat(
            {
                conversation_id: convId ?? null,
                message: payload.message,
                model_id: payload.model.id,
            },
            {
                onToken(token) {
                    streamBuffer.value += token
                },
                onDone(conversationId, messageId) {
                    const aiMsg: ChatMessage = {
                        id: messageId,
                        conversation_id: conversationId,
                        user_id: '',
                        content: streamBuffer.value,
                        model_id: payload.model.id,
                        type: 1,
                        created_at: new Date().toISOString(),
                    }
                    messages.value.push(aiMsg)
                    if (!activeId.value) {
                        setActive(conversationId)
                        fetchConversations()
                    }
                    streamBuffer.value = ''
                    scrollToBottom()
                },
                onError(err) {
                    messages.value.push({
                        id: `err-${Date.now()}`,
                        conversation_id: convId ?? '',
                        user_id: '',
                        content: err.message,
                        model_id: payload.model.id,
                        type: 1,
                        created_at: new Date().toISOString(),
                    })
                },
            },
            abortController.signal,
        )
    } finally {
        isLoading.value = false
        abortController = null
    }
}

function handleStop() {
    abortController?.abort()
    isLoading.value = false
}
</script>

<template>
    <HigBox class="w-full h-full">
        <div class="flex h-full w-full">
            <HigSidebar>
                <HigChatSessionList
                    :sessions="conversationList"
                    :active-id="activeId ?? undefined"
                    :loading="conversationsLoading"
                    @select="handleSelect"
                    @new="handleNew"
                    @delete="handleDelete"
                />
            </HigSidebar>

            <div class="flex flex-1 flex-col overflow-hidden">
                <div
                    ref="messagesEl"
                    class="flex flex-1 flex-col overflow-y-auto py-6"
                >
                    <div v-if="!activeId && conversationsError" class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                        {{ conversationsError }}
                    </div>
                    <div v-else-if="!activeId && conversationList.length === 0" class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                        Start a new chat to begin
                    </div>
                    <div v-else class="mx-auto flex w-full max-w-3xl flex-col gap-5 px-4">
                        <template v-for="message in messages" :key="message.id">
                            <HigUserMessage
                                v-if="message.type === 0"
                                :content="message.content"
                            />
                            <HigAssistantMessage
                                v-else
                                :content="message.content"
                                :streaming="false"
                            />
                        </template>
                        <HigAssistantMessage
                            v-if="isLoading && streamBuffer"
                            :content="streamBuffer"
                            :streaming="true"
                        />
                    </div>
                </div>

                <HigPromptInput
                    :loading="isLoading"
                    @send="handleSend"
                    @stop="handleStop"
                />
            </div>
        </div>
    </HigBox>

    <HigSigninModal :open="showSignin" @close="() => {}" />
</template>
