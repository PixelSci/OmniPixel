<script setup lang="ts">
interface Props {
    variant?: 'glass' | 'filled' | 'inset'
    padding?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'glass',
    padding: true,
})
</script>

<template>
    <div
        class="rounded-xl border" :class="[
            props.padding && 'p-4',
            props.variant === 'glass' && 'hig-glass',
            props.variant === 'filled' && 'bg-(--hig-secondary-background) border-(--hig-separator) shadow-sm',
            props.variant === 'inset' && 'bg-[var(--hig-grouped-background)] border-[var(--hig-opaque-separator)] shadow-xs',
        ]"
    >
        <slot />
    </div>
</template>
