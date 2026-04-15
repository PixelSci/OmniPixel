<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
    class?: HTMLAttributes['class']
}

const props = defineProps<Props>()
</script>

<template>
    <div
        :class="cn('flex h-3 items-center gap-2.5 px-1', props.class)"
        aria-hidden="true"
    >
        <span class="size-3 rounded-full bg-[#ff5f57] shadow-[inset_0_0_0_0.5px_rgba(0,0,0,0.18),0_0.5px_1px_rgba(0,0,0,0.12)]" />
        <span class="size-3 rounded-full bg-[#ffbd2e] shadow-[inset_0_0_0_0.5px_rgba(0,0,0,0.16),0_0.5px_1px_rgba(0,0,0,0.10)]" />
        <span class="size-3 rounded-full bg-[#28c840] shadow-[inset_0_0_0_0.5px_rgba(0,0,0,0.16),0_0.5px_1px_rgba(0,0,0,0.10)]" />
    </div>
</template>
