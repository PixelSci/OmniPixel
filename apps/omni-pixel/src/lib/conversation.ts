import { http, getToken } from './http'

export interface ConversationItem {
    id: string
    title: string
    is_visible: boolean
    is_archived: boolean
    created_at: string
    updated_at: string
}

export interface ChatMessage {
    id: string
    conversation_id: string
    user_id: string
    content: string
    model_id: string
    type: number // 0 = user, 1 = assistant
    created_at: string
}

export interface ConversationDetail {
    id: string
    title: string
    is_visible: boolean
    is_archived: boolean
    created_at: string
    updated_at: string
    messages: ChatMessage[]
}

export interface ChatRequest {
    conversation_id?: string | null
    message: string
    model_id: string
}

export function listConversations(): Promise<ConversationItem[]> {
    return http.get<ConversationItem[]>('/conversations')
}

export function getConversation(id: string): Promise<ConversationDetail> {
    return http.get<ConversationDetail>(`/conversations/${id}`)
}

export function deleteConversation(id: string): Promise<void> {
    return http.delete<void>(`/conversations/${id}`)
}

export function streamChat(
    payload: ChatRequest,
    callbacks: {
        onToken: (token: string) => void
        onDone: (conversationId: string, messageId: string) => void
        onError: (error: Error) => void
    },
    signal?: AbortSignal,
): Promise<void> {
    const baseUrl = '/api/v1'
    const resolved = `${window.location.origin}${baseUrl}/conversation`

    const headers: Record<string, string> = { 'Content-Type': 'application/json' }
    const token = getToken()
    if (token) headers['Authorization'] = `Bearer ${token}`

    return fetch(resolved, {
        method: 'POST',
        headers,
        body: JSON.stringify(payload),
        credentials: 'include',
        signal,
    }).then(async (response) => {
        if (!response.ok) {
            const errorBody = await response.json().catch(() => ({}))
            throw new Error(errorBody?.message ?? `Stream request failed: ${response.status}`)
        }

        const reader = response.body?.getReader()
        if (!reader) throw new Error('No readable stream available')

        const decoder = new TextDecoder()
        let buffer = ''

        while (true) {
            const { done, value } = await reader.read()
            if (done) break

            buffer += decoder.decode(value, { stream: true })
            const lines = buffer.split('\n')
            buffer = lines.pop() ?? ''

            for (const line of lines) {
                if (!line.startsWith('data: ')) continue
                const jsonStr = line.slice(6).trim()
                if (!jsonStr) continue

                let parsed: Record<string, unknown>
                try { parsed = JSON.parse(jsonStr) } catch { continue }

                if (parsed.token) {
                    callbacks.onToken(parsed.token as string)
                } else if (parsed.done) {
                    callbacks.onDone(
                        parsed.conversation_id as string,
                        parsed.message_id as string,
                    )
                } else if (parsed.error) {
                    callbacks.onError(new Error(parsed.error as string))
                }
            }
        }
    })
}
