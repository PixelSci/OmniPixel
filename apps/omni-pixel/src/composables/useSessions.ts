import { ref, type Ref } from 'vue'
import {
    listSessions,
    deleteSession as deleteSessionApi,
    type SessionItem,
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
            sessions.value = await listSessions()
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

    async function create(): Promise<SessionItem | null> {
        // Conversations are created on first chat via POST /conversation.
        // For now, create an optimistic placeholder that will be replaced
        // when the first message stream completes.
        const optimistic: SessionItem = {
            id: '',
            title: 'New Chat',
            is_visible: true,
            is_archived: false,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
        }
        return optimistic
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
