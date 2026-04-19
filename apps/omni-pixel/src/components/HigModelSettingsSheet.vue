<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { reactive } from 'vue'
import HigSheet from '@/components/HigSheet.vue'
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
</script>

<template>
    <HigSheet v-model:open="openModel" title="模型设置">
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
    </HigSheet>
</template>
