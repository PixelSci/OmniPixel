<script setup lang="ts">
interface Props {
    padding?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    padding: true,
})
</script>

<template>
    <div
        class="relative inline-block rounded-xl bg-white overflow-hidden shadow-[0px_0px_0px_1px_rgba(255,255,255,0.23),0px_16px_48px_0px_rgba(0,0,0,0.35)] dark:bg-[var(--hig-secondary-grouped-background)] dark:text-[var(--hig-label)] dark:shadow-[0px_0px_0px_1px_rgba(255,255,255,0.14),0px_16px_48px_0px_rgba(0,0,0,0.58)]"
    >
        <div class="relative" :class="props.padding && 'p-4'">
            <slot />
        </div>
    </div>
</template>
