import { useLocalStorage } from '@vueuse/core'
import { computed, ref } from 'vue'
import { listModels, toChatModel, type ChatModel, type ModelResponse } from '@/lib/model'

const STORAGE_KEY = 'omni-pixel:enabled-models'

const apiModels = ref<ModelResponse[]>([])
const loaded = ref(false)

async function loadModels() {
    if (loaded.value) return
    try {
        apiModels.value = await listModels()
    } catch { /* use empty fallback */ }
    loaded.value = true
}

const enabledIds = useLocalStorage<string[]>(STORAGE_KEY, [])

export function useModelSettings() {
    const chatModels = computed<ChatModel[]>(() => {
        const models = apiModels.value
        if (!models.length) return []
        const enabled = new Set(enabledIds.value.length ? enabledIds.value : models.map(m => m.id))
        return models.filter(m => enabled.has(m.id)).map(toChatModel)
    })

    function isModelEnabled(modelId: string): boolean {
        if (!enabledIds.value.length) return true
        return enabledIds.value.includes(modelId)
    }

    function setModelEnabled(modelId: string, enabled: boolean) {
        const base = enabledIds.value.length ? enabledIds.value : apiModels.value.map(m => m.id)
        const set = new Set(base)
        if (enabled) set.add(modelId)
        else set.delete(modelId)
        enabledIds.value = [...set]
    }

    loadModels()

    return {
        chatModels,
        isModelEnabled,
        setModelEnabled,
    }
}
