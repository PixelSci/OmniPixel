import { ref, type Ref } from 'vue'
import {
    listConversations,
    deleteConversation as deleteConversationApi,
    type ConversationItem,
} from '@/lib/conversation'
import { ApiError } from '@/lib/http'

export function useConversations() {
    const conversations: Ref<ConversationItem[]> = ref([])
    const activeId: Ref<string | null> = ref(null)
    const loading = ref(false)
    const error = ref<string | null>(null)

    async function fetchConversations() {
        loading.value = true
        error.value = null
        try {
            conversations.value = await listConversations()
        } catch (e) {
            if (e instanceof ApiError && e.status === 401) {
                conversations.value = []
                error.value = 'Please sign in'
            } else {
                error.value = e instanceof Error ? e.message : 'Failed to load conversations'
            }
        } finally {
            loading.value = false
        }
    }

    async function create(): Promise<ConversationItem | null> {
        const optimistic: ConversationItem = {
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
            await deleteConversationApi(id)
            conversations.value = conversations.value.filter(s => s.id !== id)
            if (activeId.value === id) {
                activeId.value = null
            }
        } catch (e) {
            error.value = e instanceof Error ? e.message : 'Failed to delete conversation'
        }
    }

    function setActive(id: string | null) {
        activeId.value = id
    }

    return {
        conversations,
        activeId,
        loading,
        error,
        fetchConversations,
        create,
        remove,
        setActive,
    }
}
