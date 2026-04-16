<script setup lang="ts">
import { ref } from 'vue'

const activeSessionId = ref<string>('1')
const isLoading = ref(false)

function handleSelect(id: string) {
    activeSessionId.value = id
}

function handleNew() {
    activeSessionId.value = ''
}

function handleSend(message: string) {
    console.log('send:', message)
    // placeholder: simulate loading state
    isLoading.value = true
    setTimeout(() => { isLoading.value = false }, 3000)
}

function handleStop() {
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
                        <HigUserMessage content="What is the best way to share state between components in Vue 3?" />
                        <HigAssistantMessage
                            model="Claude Sonnet 4"
                            content="The recommended approach in Vue 3 is to use composables — plain functions that encapsulate reactive state and expose it via the Composition API.

For local, parent-child communication you can rely on props and emits. But when multiple unrelated components need to share the same piece of state, extract it into a composable module:

// useCounter.ts
import { ref } from 'vue'
const count = ref(0)
export function useCounter() {
  return { count, increment: () => count.value++ }
}

Because the ref is declared at module scope, every component that calls useCounter() shares the same reactive instance. For more complex cases, Pinia is the official state-management library and integrates naturally with Vue DevTools."
                        />
                        <HigUserMessage content="Can you show me a Pinia example?" />
                        <HigAssistantMessage
                            model="Claude Sonnet 4"
                            content="Sure! Here is a minimal Pinia store and how to use it in a component."
                        />
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
