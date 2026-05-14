<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSessions } from '@/composables/useSessions'
import { useModelApiKeys } from '@/composables/useModelApiKeys'
import { sendChat, getSession, type ChatMessage } from '@/lib/session'
import type { FeatureModel } from '@/composables/useModelSettings'

const {
    sessions: sessionList,
    activeId,
    loading: sessionsLoading,
    error: sessionsError,
    fetchSessions,
    create,
    remove,
    setActive,
} = useSessions()
const { keyForProvider } = useModelApiKeys()

const messages = ref<ChatMessage[]>([])
const isLoading = ref(false)
const isCreatingSession = ref(false)

let abortController: AbortController | null = null

function nextMessageId() {
    return `msg-${Date.now()}-${Math.random().toString(16).slice(2)}`
}

onMounted(() => {
    fetchSessions()
})

async function handleSelect(id: string) {
    setActive(id)
    try {
        const detail = await getSession(id)
        messages.value = detail.messages ?? []
    } catch {
        messages.value = []
    }
}

async function handleNew() {
    isCreatingSession.value = true
    const session = await create()
    isCreatingSession.value = false
    if (session) {
        setActive(session.id)
        messages.value = []
    }
}

async function handleDelete(id: string) {
    await remove(id)
}

async function handleSend(payload: { message: string, model: FeatureModel }) {
    const sessionId = activeId.value
    if (!sessionId) return

    const userMessage: ChatMessage = {
        id: nextMessageId(),
        role: 'user',
        content: payload.message,
    }
    messages.value.push(userMessage)

    const apiKey = keyForProvider(payload.model.providerId)
    if (!apiKey) {
        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content: `请先在个人中心添加 ${payload.model.provider} API Key，然后再发送消息。`,
        })
        return
    }

    isLoading.value = true
    abortController = new AbortController()

    try {
        const res = await sendChat(sessionId, {
            prompt: payload.message,
            provider: payload.model.providerId,
            model: payload.model.id,
            api_key: apiKey,
        })
        messages.value.push(res.assistant_message)
    } catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError') return

        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content: error instanceof Error ? error.message : '调用模型时发生未知错误。',
        })
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
                    :sessions="sessionList"
                    :active-id="activeId ?? undefined"
                    :loading="sessionsLoading"
                    @select="handleSelect"
                    @new="handleNew"
                    @delete="handleDelete"
                />
            </HigSidebar>

            <!-- chat area -->
            <div class="flex flex-1 flex-col overflow-hidden">
                <!-- messages -->
                <div class="flex flex-1 flex-col overflow-y-auto py-6">
                    <div v-if="!activeId && sessionsError" class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                        {{ sessionsError }}
                    </div>
                    <div v-else-if="!activeId" class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                        Select or start a new chat
                    </div>
                    <div v-else class="mx-auto flex w-full max-w-3/5 flex-col gap-5 px-4">
                        <template v-for="message in messages" :key="message.id">
                            <HigUserMessage
                                v-if="message.role === 'user'"
                                :content="message.content"
                            />
                            <HigAssistantMessage
                                v-else
                                :model="message.model"
                                :content="message.content"
                            />
                        </template>
                    </div>
                </div>

                <!-- prompt input -->
                <HigPromptInput
                    :loading="isLoading"
                    :disabled="isCreatingSession"
                    @send="handleSend"
                    @stop="handleStop"
                />
            </div>
        </div>
    </HigBox>
</template>
