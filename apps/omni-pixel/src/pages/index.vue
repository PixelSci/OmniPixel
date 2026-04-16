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
                <div class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                    <span v-if="!activeSessionId">Select or start a new chat</span>
                    <span v-else>Session {{ activeSessionId }}</span>
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
