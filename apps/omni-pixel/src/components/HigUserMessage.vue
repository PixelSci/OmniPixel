<script setup lang="ts">
interface Props {
    content: string
}

defineProps<Props>()
</script>

<template>
    <div class="flex justify-end">
        <div
            class="max-w-[72%] rounded-[18px] bg-[#0077ED] px-4 py-2.5 text-white shadow-[0_1px_2px_rgba(0,0,119,0.18),0_2px_8px_rgba(0,119,237,0.22)] dark:bg-[#0a7aff]"
        >
            <p class="whitespace-pre-wrap text-[15px] leading-[20px]">{{ content }}</p>
        </div>
    </div>
</template>
