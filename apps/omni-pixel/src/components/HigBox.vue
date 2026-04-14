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
        class="relative inline-block rounded-xl overflow-hidden shadow-[0px_0px_2px_0px_rgba(0,0,0,0.10),0px_0px_25px_0px_rgba(0,0,0,0.16)]"
    >
        <!-- backdrop blur -->
        <div class="absolute inset-0 pointer-events-none backdrop-blur-[28px] saturate-180" aria-hidden="true" />
        <!-- color-dodge vibrancy highlight (light mode) -->
        <div class="absolute inset-0 pointer-events-none bg-white/40 mix-blend-color-dodge dark:hidden" aria-hidden="true" />
        <!-- translucent fill -->
        <div class="absolute inset-0 pointer-events-none bg-[rgba(250,250,252,0.86)] dark:bg-[rgba(44,44,46,0.78)]" aria-hidden="true" />
        <!-- content -->
        <div class="relative" :class="props.padding && 'p-4'">
            <slot />
        </div>
    </div>
</template>
