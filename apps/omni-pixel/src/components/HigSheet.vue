<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { onBeforeUnmount, onMounted, reactive, watch } from 'vue'
import HigTrafficLights from '@/components/HigTrafficLights.vue'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'

interface Props {
    open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const openModel = useVModel(props, 'open', emit, { passive: true })

interface ModelSection {
    key: 'chat' | 'image' | 'audio' | 'video'
    title: string
    description: string
    options: string[]
}

const sections: ModelSection[] = [
    {
        key: 'chat',
        title: 'Chat 模型',
        description: '用于对话与文本生成',
        options: ['Claude Opus 4.7', 'Claude Sonnet 4.6', 'GPT-4o', 'DeepSeek V3'],
    },
    {
        key: 'image',
        title: 'Image 模型',
        description: '用于图像生成与编辑',
        options: ['DALL·E 3', 'Midjourney v6', 'Stable Diffusion XL', 'Flux.1'],
    },
    {
        key: 'audio',
        title: 'Audio 模型',
        description: '用于语音合成与识别',
        options: ['ElevenLabs v2', 'OpenAI TTS', 'Whisper Large v3'],
    },
    {
        key: 'video',
        title: 'Video 模型',
        description: '用于视频生成',
        options: ['Sora', 'Runway Gen-3', 'Pika 1.5', 'Veo 2'],
    },
]

const selected = reactive<Record<ModelSection['key'], string>>({
    chat: sections[0].options[0],
    image: sections[1].options[0],
    audio: sections[2].options[0],
    video: sections[3].options[0],
})

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
                    class="relative inline-block w-[520px] max-w-[92vw] overflow-hidden rounded-xl bg-white shadow-[0px_0px_0px_1px_rgba(255,255,255,0.23),0px_16px_48px_0px_rgba(0,0,0,0.35)] dark:bg-[var(--hig-secondary-grouped-background)] dark:text-[var(--hig-label)] dark:shadow-[0px_0px_0px_1px_rgba(255,255,255,0.14),0px_16px_48px_0px_rgba(0,0,0,0.58)]"
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
                        <div class="pointer-events-none absolute inset-x-0 text-center text-[13px] font-semibold tracking-tight">
                            个人中心
                        </div>
                    </div>

                    <div class="flex flex-col gap-5 px-5 py-5">
                        <section
                            v-for="section in sections"
                            :key="section.key"
                            class="flex items-center justify-between gap-4"
                        >
                            <div class="flex flex-col">
                                <span class="text-[13px] font-[510] text-foreground">{{ section.title }}</span>
                                <span class="text-[12px] text-[var(--hig-secondary-label)]">{{ section.description }}</span>
                            </div>
                            <NativeSelect v-model="selected[section.key]" class="w-[200px]">
                                <NativeSelectOption
                                    v-for="option in section.options"
                                    :key="option"
                                    :value="option"
                                >
                                    {{ option }}
                                </NativeSelectOption>
                            </NativeSelect>
                        </section>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>
