<script setup lang="ts">
import { Ellipsis, Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { computed, nextTick, ref } from 'vue'
import type { ConversationItem } from '@/lib/conversation'
import { updateConversationTitle } from '@/lib/conversation'

const props = defineProps<Props>()

const emit = defineEmits<{
    select: [id: string]
    new: []
    rename: [id: string, title: string]
    delete: [id: string]
}>()

interface Props {
    sessions: ConversationItem[]
    activeId?: string
    loading?: boolean
}

const editingId = ref<string | null>(null)
const editTitle = ref('')
const editInput = ref<HTMLInputElement | null>(null)

async function startEdit(id: string, currentTitle: string) {
    editingId.value = id
    editTitle.value = currentTitle
    await nextTick()
    setTimeout(() => {
        editInput.value?.focus()
        editInput.value?.select()
    })
}

async function commitEdit(id: string) {
    const title = editTitle.value.trim()
    editingId.value = null
    if (title) {
        await updateConversationTitle(id, title)
        emit('rename', id, title)
    }
}

function onEditKeydown(e: KeyboardEvent, id: string) {
    if (e.key === 'Enter') {
        e.preventDefault()
        commitEdit(id)
    } else if (e.key === 'Escape') {
        editingId.value = null
    }
}

// ── date grouping ────────────────────────────────────────────────────
interface Group { label: string, sessions: ConversationItem[] }

const groups = computed<Group[]>(() => {
    const now = new Date()
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const startOfYesterday = new Date(startOfToday.getTime() - 86400000)
    const startOf7Days = new Date(startOfToday.getTime() - 6 * 86400000)
    const startOf30Days = new Date(startOfToday.getTime() - 29 * 86400000)

    const buckets: Record<string, ConversationItem[]> = {
        'Today': [],
        'Yesterday': [],
        'Previous 7 Days': [],
        'Previous 30 Days': [],
        'Older': [],
    }

    for (const s of props.sessions) {
        const t = new Date(s.updated_at).getTime()
        if (t >= startOfToday.getTime())
            buckets.Today.push(s)
        else if (t >= startOfYesterday.getTime())
            buckets.Yesterday.push(s)
        else if (t >= startOf7Days.getTime())
            buckets['Previous 7 Days'].push(s)
        else if (t >= startOf30Days.getTime())
            buckets['Previous 30 Days'].push(s)
        else
            buckets.Older.push(s)
    }

    return Object.entries(buckets)
        .filter(([, items]) => items.length > 0)
        .map(([label, sessions]) => ({ label, sessions }))
})

function formatTime(dateStr: string): string {
    const date = new Date(dateStr)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const mins = Math.floor(diff / 60000)
    if (mins < 1) return 'Just now'
    if (mins < 60) return `${mins}m ago`
    const hrs = Math.floor(mins / 60)
    if (hrs < 24) return `${hrs}h ago`
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}
</script>

<template>
    <div class="flex flex-col gap-0.5">
        <!-- new chat button -->
        <button
            class="group mb-1 flex w-full items-center gap-2 rounded-lg px-2 py-1.5 text-left text-[13px] text-[var(--hig-secondary-label)] transition-colors hover:bg-[rgba(0,0,0,0.06)] dark:hover:bg-[rgba(255,255,255,0.06)]"
            @click="emit('new')"
        >
            <Plus class="size-3.5 shrink-0" />
            <span>New Chat</span>
        </button>

        <div v-if="loading && sessions.length === 0" class="px-2 py-3 text-center text-[12px] text-[var(--hig-tertiary-label)]">
            Loading sessions…
        </div>
        <template v-for="group in groups" :key="group.label">
            <p class="mt-2 mb-0.5 px-2 text-[11px] font-semibold uppercase tracking-wide text-[var(--hig-tertiary-label,rgba(60,60,67,0.3))] dark:text-[rgba(235,235,245,0.3)] first:mt-0">
                {{ group.label }}
            </p>

            <button
                v-for="session in group.sessions"
                :key="session.id"
                class="group relative flex w-full flex-col items-start rounded-lg px-2 py-1.5 text-left transition-colors"
                :class="activeId === session.id
                    ? 'bg-[#0077ED] text-white'
                    : 'text-[var(--hig-label)] hover:bg-[rgba(0,0,0,0.06)] dark:hover:bg-[rgba(255,255,255,0.06)]'"
                @click="editingId !== session.id && emit('select', session.id)"
            >
                <div class="flex w-full items-baseline gap-1">
                    <!-- title / inline edit -->
                    <input
                        v-if="editingId === session.id"
                        ref="editInput"
                        v-model="editTitle"
                        type="text"
                        class="h-[22px] flex-1 border-0 bg-transparent p-0 text-[13px] font-medium leading-snug outline-none"
                        :class="activeId === session.id ? 'text-white placeholder-white/50' : 'text-[var(--hig-label)]'"
                        @click.stop
                        @blur="commitEdit(session.id)"
                        @keydown="onEditKeydown($event, session.id)"
                    />
                    <span
                        v-else
                        class="flex-1 truncate text-[13px] font-medium leading-snug"
                        :class="activeId !== session.id && 'text-[var(--hig-label)]'"
                    >
                        {{ session.title }}
                    </span>
                    <!-- timestamp + ellipsis overlay -->
                    <span v-if="editingId !== session.id" class="relative shrink-0">
                        <span
                            class="text-[11px] tabular-nums transition-opacity group-hover:opacity-0"
                            :class="activeId === session.id
                                ? 'text-white/70'
                                : 'text-[var(--hig-tertiary-label,rgba(60,60,67,0.3))] dark:text-[rgba(235,235,245,0.3)]'"
                        >
                            {{ formatTime(session.updated_at) }}
                        </span>
                        <UiDropdownMenuDropdownMenu>
                            <UiDropdownMenuDropdownMenuTrigger as-child>
                                <button
                                    class="absolute inset-0 flex items-center justify-center rounded p-0.5 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-black/10 dark:hover:bg-white/10"
                                    :class="activeId === session.id ? 'text-white/70' : 'text-[var(--hig-tertiary-label)]'"
                                    @click.stop
                                    aria-label="Actions"
                                >
                                    <Ellipsis class="size-4" />
                                </button>
                            </UiDropdownMenuDropdownMenuTrigger>
                            <UiDropdownMenuDropdownMenuContent side="right" align="start" class="min-w-[130px]">
                                <UiDropdownMenuDropdownMenuItem @click.stop="startEdit(session.id, session.title)">
                                    <Pencil class="size-3.5" />
                                    <span>修改标题</span>
                                </UiDropdownMenuDropdownMenuItem>
                                <UiDropdownMenuDropdownMenuItem @click.stop="emit('delete', session.id)">
                                    <Trash2 class="size-3.5" />
                                    <span>删除对话</span>
                                </UiDropdownMenuDropdownMenuItem>
                            </UiDropdownMenuDropdownMenuContent>
                        </UiDropdownMenuDropdownMenu>
                    </span>
                </div>
            </button>
        </template>
    </div>
</template>
