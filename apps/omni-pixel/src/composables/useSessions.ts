import { ref, type Ref } from 'vue'
import {
    listSessions,
    createSession,
    deleteSession as deleteSessionApi,
    type SessionItem,
    type SessionDetail,
} from '@/lib/session'
import { ApiError } from '@/lib/http'

export function useSessions() {
    const sessions: Ref<SessionItem[]> = ref([])
    const activeId: Ref<string | null> = ref(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    async function fetchSessions() {
        loading.value = true
        error.value = null
        try {
            const res = await listSessions()
            sessions.value = res.sessions
        } catch (e) {
            if (e instanceof ApiError && e.status === 401) {
                sessions.value = []
                error.value = 'Please sign in'
            } else {
                error.value = e instanceof Error ? e.message : 'Failed to load sessions'
            }
        } finally {
            loading.value = false
        }
    }

    async function create(): Promise<SessionDetail | null> {
        try {
            const session = await createSession()
            sessions.value.unshift({
                id: session.id,
                title: session.title,
                preview: session.preview,
                model: session.model,
                updated_at: session.updated_at,
                last_chat_at: null,
            })
            return session
        } catch (e) {
            error.value = e instanceof Error ? e.message : 'Failed to create session'
            return null
        }
    }

    async function remove(id: string) {
        try {
            await deleteSessionApi(id)
            sessions.value = sessions.value.filter(s => s.id !== id)
            if (activeId.value === id) {
                activeId.value = null
            }
        } catch (e) {
            error.value = e instanceof Error ? e.message : 'Failed to delete session'
        }
    }

    function setActive(id: string | null) {
        activeId.value = id
    }

    return {
        sessions,
        activeId,
        loading,
        error,
        fetchSessions,
        create,
        remove,
        setActive,
    }
}
