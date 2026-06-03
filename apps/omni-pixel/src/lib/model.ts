import { http } from './http'

export interface ModelResponse {
    id: string
    provider_id: string
    provider_name: string
    model_name: string
    is_enabled: boolean
    expire_time: string | null
    created_at: string
    updated_at: string
}

export interface ChatModel {
    id: string
    name: string
    providerName: string
}

export function listModels(): Promise<ModelResponse[]> {
    return http.get<{ models: ModelResponse[] }>('/models').then((r) => r.models)
}

export function toChatModel(m: ModelResponse): ChatModel {
    return {
        id: m.id,
        name: m.model_name,
        providerName: m.provider_name,
    }
}
