import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

export type ModelFeature = 'chat' | 'image' | 'audio' | 'video'

export interface FeatureModel {
    id: string
    name: string
    provider: string
    providerId: string
    badge?: string
}

export interface ModelFeatureSection {
    key: ModelFeature
    title: string
    description: string
    models: FeatureModel[]
}

const STORAGE_KEY = 'omni-pixel:model-feature-settings'

export const modelFeatureSections: ModelFeatureSection[] = [
    {
        key: 'chat',
        title: 'Chat 模型',
        description: '用于对话与文本生成',
        models: [
            { id: 'claude-sonnet-4', name: 'Claude Sonnet 4', provider: 'Anthropic', providerId: 'anthropic', badge: 'Smart' },
            { id: 'claude-haiku-4', name: 'Claude Haiku 4', provider: 'Anthropic', providerId: 'anthropic', badge: 'Fast' },
            { id: 'gpt-4o', name: 'GPT-4o', provider: 'OpenAI', providerId: 'openai' },
            { id: 'gpt-4o-mini', name: 'GPT-4o mini', provider: 'OpenAI', providerId: 'openai', badge: 'Fast' },
            { id: 'deepseek-chat', name: 'DeepSeek V4', provider: 'DeepSeek', providerId: 'deepseek' },
        ],
    },
    {
        key: 'image',
        title: 'Image 模型',
        description: '用于图像生成与编辑',
        models: [
            { id: 'dall-e-3', name: 'DALL·E 3', provider: 'OpenAI', providerId: 'openai' },
            { id: 'midjourney-v6', name: 'Midjourney v6', provider: 'Midjourney', providerId: 'other' },
            { id: 'stable-diffusion-xl', name: 'Stable Diffusion XL', provider: 'Stability AI', providerId: 'other' },
            { id: 'flux-1', name: 'Flux.1', provider: 'Black Forest Labs', providerId: 'other' },
        ],
    },
    {
        key: 'audio',
        title: 'Audio 模型',
        description: '用于语音合成与识别',
        models: [
            { id: 'elevenlabs-v2', name: 'ElevenLabs v2', provider: 'ElevenLabs', providerId: 'other' },
            { id: 'openai-tts', name: 'OpenAI TTS', provider: 'OpenAI', providerId: 'openai' },
            { id: 'whisper-large-v3', name: 'Whisper Large v3', provider: 'OpenAI', providerId: 'openai' },
        ],
    },
    {
        key: 'video',
        title: 'Video 模型',
        description: '用于视频生成',
        models: [
            { id: 'sora', name: 'Sora', provider: 'OpenAI', providerId: 'openai' },
            { id: 'runway-gen-3', name: 'Runway Gen-3', provider: 'Runway', providerId: 'other' },
            { id: 'pika-1-5', name: 'Pika 1.5', provider: 'Pika', providerId: 'other' },
            { id: 'veo-2', name: 'Veo 2', provider: 'Google', providerId: 'google' },
        ],
    },
]

function defaultSettings(): Record<ModelFeature, string[]> {
    return modelFeatureSections.reduce((settings, section) => {
        settings[section.key] = section.models.map(model => model.id)
        return settings
    }, {} as Record<ModelFeature, string[]>)
}

const settings = useLocalStorage<Record<ModelFeature, string[]>>(STORAGE_KEY, defaultSettings(), {
    mergeDefaults: true,
})

function sectionFor(feature: ModelFeature) {
    return modelFeatureSections.find(section => section.key === feature)!
}

export function useModelSettings() {
    function enabledModelIds(feature: ModelFeature) {
        const enabledIds = settings.value[feature]
        return Array.isArray(enabledIds)
            ? enabledIds.map(id => id === 'deepseek-v3' ? 'deepseek-chat' : id)
            : sectionFor(feature).models.map(model => model.id)
    }

    function enabledModels(feature: ModelFeature) {
        const enabledIds = enabledModelIds(feature)
        return sectionFor(feature).models.filter(model => enabledIds.includes(model.id))
    }

    function isModelEnabled(feature: ModelFeature, modelId: string) {
        return enabledModelIds(feature).includes(modelId)
    }

    function setModelEnabled(feature: ModelFeature, modelId: string, enabled: boolean) {
        const enabledIds = new Set(enabledModelIds(feature))

        if (enabled)
            enabledIds.add(modelId)
        else
            enabledIds.delete(modelId)

        settings.value[feature] = [...enabledIds]
    }

    const chatModels = computed(() => enabledModels('chat'))

    return {
        sections: modelFeatureSections,
        settings,
        chatModels,
        enabledModels,
        isModelEnabled,
        setModelEnabled,
    }
}
