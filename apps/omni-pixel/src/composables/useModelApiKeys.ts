import { useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

export interface ModelApiKeyEntry {
    id: string
    provider: string
    label: string
    key: string
    show: boolean
}

const STORAGE_KEY = 'omni-pixel:model-api-keys'

let idSeed = Date.now()

function nextId() {
    idSeed += 1
    return `key-${idSeed}`
}

const keys = useLocalStorage<ModelApiKeyEntry[]>(STORAGE_KEY, [], {
    mergeDefaults: true,
})

export function createModelApiKey(provider = 'openai'): ModelApiKeyEntry {
    return {
        id: nextId(),
        provider,
        label: '',
        key: '',
        show: false,
    }
}

export function useModelApiKeys() {
    function keyForProvider(provider: string) {
        return keys.value.find(entry => entry.provider === provider && entry.key.trim())?.key.trim() || ''
    }

    const openAIKey = computed(() => keyForProvider('openai'))

    function addKey(provider = 'openai') {
        keys.value.push(createModelApiKey(provider))
    }

    function removeKey(id: string) {
        const idx = keys.value.findIndex(entry => entry.id === id)
        if (idx !== -1)
            keys.value.splice(idx, 1)
    }

    return {
        keys,
        openAIKey,
        keyForProvider,
        addKey,
        removeKey,
    }
}
