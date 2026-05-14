import { api } from './api'

export interface SessionItem {
    id: string
    title: string
    preview?: string
    model?: string
    updated_at: string
    last_chat_at?: string | null
}

export interface SessionListResponse {
    sessions: SessionItem[]
}

export interface ChatMessage {
    id: string
    role: 'user' | 'assistant'
    content: string
    model?: string
}

export interface SessionDetail {
    id: string
    title: string
    preview?: string
    model?: string
    messages: ChatMessage[]
    created_at: string
    updated_at: string
}

export interface SendChatPayload {
    prompt: string
    provider: string
    model: string
    api_key: string
}

export interface SendChatResponse {
    session_id: string
    created_session: boolean
    message: ChatMessage
    assistant_message: ChatMessage
    messages: ChatMessage[]
}

export function listSessions(): Promise<SessionListResponse> {
    return api.get<SessionListResponse>('/sessions')
}

export function getSession(id: string): Promise<SessionDetail> {
    return api.get<SessionDetail>(`/sessions/${id}`)
}

export function createSession(): Promise<SessionDetail> {
    const id = crypto.randomUUID()
    return api.post<SessionDetail>('/sessions', {
        session_id: id,
        title: 'New Chat',
        preview: '',
        model: '',
    })
}

export function deleteSession(id: string): Promise<void> {
    return api.delete<void>(`/sessions/${id}`)
}

export function sendChat(sessionId: string, payload: SendChatPayload): Promise<SendChatResponse> {
    return api.post<SendChatResponse>(`/sessions/${sessionId}/chat`, {
        session_id: sessionId,
        ...payload,
    })
}
