<script setup lang="ts">
import { ref } from 'vue'

const activeSessionId = ref<string>('1')

function handleSelect(id: string) {
    activeSessionId.value = id
}

function handleNew() {
    activeSessionId.value = ''
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

            <!-- main content area -->
            <div class="flex flex-1 items-center justify-center text-[var(--hig-secondary-label)] text-[13px]">
                <span v-if="!activeSessionId">Start a new chat</span>
                <span v-else>Session {{ activeSessionId }}</span>
            </div>
        </div>
    </HigBox>
</template>
