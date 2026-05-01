<script setup lang="ts">
import { ref } from 'vue'
import { useModelApiKeys } from '@/composables/useModelApiKeys'

const activeSessionId = ref<string>('1')
const isLoading = ref(false)
const { keyForProvider } = useModelApiKeys()

interface ChatModel {
    id: string
    name: string
    provider: string
    providerId: string
}

interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    model?: string
}

const messages = ref<ChatMessage[]>([
    {
        id: 'seed-user-1',
        role: 'user',
        content: 'What is the best way to share state between components in Vue 3?',
    },
    {
        id: 'seed-assistant-1',
        role: 'assistant',
        model: 'Claude Sonnet 4',
        content: `The recommended approach in Vue 3 is to use composables — plain functions that encapsulate reactive state and expose it via the Composition API.

For local, parent-child communication you can rely on props and emits. But when multiple unrelated components need to share the same piece of state, extract it into a composable module:

// useCounter.ts
import { ref } from 'vue'
const count = ref(0)
export function useCounter() {
  return { count, increment: () => count.value++ }
}

Because the ref is declared at module scope, every component that calls useCounter() shares the same reactive instance. For more complex cases, Pinia is the official state-management library and integrates naturally with Vue DevTools.`,
    },
    {
        id: 'seed-user-2',
        role: 'user',
        content: 'Can you show me a Pinia example?',
    },
    {
        id: 'seed-assistant-2',
        role: 'assistant',
        model: 'Claude Sonnet 4',
        content: 'Sure! Here is a minimal Pinia store and how to use it in a component.',
    },
])

let abortController: AbortController | null = null

const chatProviderConfig: Record<string, { endpoint: string, providerName: string }> = {
    openai: {
        endpoint: 'https://api.openai.com/v1/chat/completions',
        providerName: 'OpenAI',
    },
    deepseek: {
        endpoint: 'https://api.deepseek.com/chat/completions',
        providerName: 'DeepSeek',
    },
}

function nextMessageId() {
    return `msg-${Date.now()}-${Math.random().toString(16).slice(2)}`
}

function handleSelect(id: string) {
    activeSessionId.value = id
}

function handleNew() {
    activeSessionId.value = ''
}

async function handleSend(payload: { message: string, model: ChatModel }) {
    messages.value.push({
        id: nextMessageId(),
        role: 'user',
        content: payload.message,
    })

    const providerConfig = chatProviderConfig[payload.model.providerId]

    if (!providerConfig) {
        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content: `当前还没有接入 ${payload.model.provider} 的 Chat 接口。`,
        })
        return
    }

    const apiKey = keyForProvider(payload.model.providerId)

    if (!apiKey) {
        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content: `请先在个人中心添加 ${providerConfig.providerName} API Key，然后再发送消息。`,
        })
        return
    }

    isLoading.value = true
    abortController = new AbortController()

    try {
        const response = await fetch(providerConfig.endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${apiKey}`,
            },
            body: JSON.stringify({
                model: payload.model.id,
                messages: messages.value
                    .filter(message => message.role === 'user' || message.role === 'assistant')
                    .map(message => ({
                        role: message.role,
                        content: message.content,
                    })),
            }),
            signal: abortController.signal,
        })

        if (!response.ok) {
            const error = await response.json().catch(() => null)
            throw new Error(error?.error?.message || `${providerConfig.providerName} request failed: ${response.status}`)
        }

        const data = await response.json()
        const content = data?.choices?.[0]?.message?.content?.trim() || '模型没有返回内容。'

        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content,
        })
    }
    catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError')
            return

        messages.value.push({
            id: nextMessageId(),
            role: 'assistant',
            model: payload.model.name,
            content: error instanceof Error ? error.message : '调用模型时发生未知错误。',
        })
    }
    finally {
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
                    :active-id="activeSessionId"
                    @select="handleSelect"
                    @new="handleNew"
                />
            </HigSidebar>

            <!-- chat area -->
            <div class="flex flex-1 flex-col overflow-hidden">
                <!-- messages -->
                <div class="flex flex-1 flex-col overflow-y-auto py-6">
                    <div v-if="!activeSessionId" class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
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
                    @send="handleSend"
                    @stop="handleStop"
                />
            </div>
        </div>
    </HigBox>
</template>
