<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { computed, onBeforeUnmount, onMounted, watch } from 'vue'
import HigTrafficLights from '@/components/HigTrafficLights.vue'

interface Props {
    open: boolean
    title?: string
    width?: string
}

const props = withDefaults(defineProps<Props>(), {
    title: '',
    width: '520px',
})

const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const openModel = useVModel(props, 'open', emit, { passive: true })

const panelStyle = computed(() => ({
    width: props.width,
}))

function close() {
    openModel.value = false
}

function handleKey(e: KeyboardEvent) {
    if (e.key === 'Escape' && openModel.value)
        close()
}

onMounted(() => document.addEventListener('keydown', handleKey))
onBeforeUnmount(() => document.removeEventListener('keydown', handleKey))

watch(openModel, (v) => {
    document.body.style.overflow = v ? 'hidden' : ''
})
</script>

<template>
    <Teleport to="body">
        <Transition
            enter-active-class="transition-opacity duration-150"
            leave-active-class="transition-opacity duration-100"
            enter-from-class="opacity-0"
            leave-to-class="opacity-0"
        >
            <div
                v-if="openModel"
                class="fixed inset-0 z-[2000] flex items-center justify-center bg-black/30 backdrop-blur-sm"
                @click.self="close"
            >
                <div
                    class="relative inline-block max-w-[92vw] overflow-hidden rounded-xl bg-white shadow-[0px_0px_0px_1px_rgba(255,255,255,0.23),0px_16px_48px_0px_rgba(0,0,0,0.35)] dark:bg-[var(--hig-secondary-grouped-background)] dark:text-[var(--hig-label)] dark:shadow-[0px_0px_0px_1px_rgba(255,255,255,0.14),0px_16px_48px_0px_rgba(0,0,0,0.58)]"
                    :style="panelStyle"
                >
                    <div class="relative flex items-center h-11 px-3 border-b border-[rgba(0,0,0,0.06)] dark:border-[rgba(255,255,255,0.08)]">
                        <button
                            type="button"
                            aria-label="关闭"
                            class="group flex items-center outline-none"
                            @click="close"
                        >
                            <HigTrafficLights />
                        </button>
                        <div
                            v-if="title"
                            class="pointer-events-none absolute inset-x-0 text-center text-[13px] font-semibold tracking-tight"
                        >
                            {{ title }}
                        </div>
                    </div>

                    <div class="flex flex-col">
                        <slot />
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>
