<script setup lang="ts">
</script>

<template>
    <div class="flex h-screen w-screen flex-col overflow-hidden rounded-2xl">
        <HigMenuBar class="shrink-0" app-name="OmniPixel" />
        <main class="min-h-0 flex-1 overflow-hidden">
            <slot />
        </main>
    </div>
</template>
