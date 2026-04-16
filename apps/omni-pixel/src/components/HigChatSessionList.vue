<script lang="ts">
export interface ChatSession {
    id: string
    title: string
    preview?: string
    updatedAt: Date
    model?: string
}

// module-level constant — accessible to defineProps default factory
const defaultSessions: ChatSession[] = [
    {
        id: '1',
        title: 'Tailwind v4 migration guide',
        preview: 'How do I migrate my existing Tailwind CSS v3 project to v4?',
        updatedAt: new Date(Date.now() - 1000 * 60 * 14),
        model: 'GPT-4o',
    },
    {
        id: '2',
        title: 'Vue composable pattern',
        preview: 'What is the best way to share state between components using composables?',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 2),
        model: 'Claude 3.5',
    },
    {
        id: '3',
        title: 'Regex for email validation',
        preview: 'Write a regex that validates common email address formats.',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 5),
        model: 'GPT-4o',
    },
    {
        id: '4',
        title: 'Docker multi-stage build',
        preview: 'Explain how multi-stage Docker builds reduce final image size.',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24),
        model: 'Claude 3.5',
    },
    {
        id: '5',
        title: 'Postgres full-text search',
        preview: 'How do I implement full-text search with ranking in PostgreSQL?',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 26),
        model: 'GPT-4o',
    },
    {
        id: '6',
        title: 'TypeScript generics deep dive',
        preview: 'Explain conditional types and infer keyword with real examples.',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3),
        model: 'Claude 3.5',
    },
    {
        id: '7',
        title: 'Rust ownership explained',
        preview: 'I keep getting borrow checker errors. Can you explain ownership?',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5),
        model: 'GPT-4o',
    },
    {
        id: '8',
        title: 'CSS grid vs flexbox',
        preview: 'When should I prefer CSS Grid over Flexbox for layout?',
        updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 10),
        model: 'Claude 3.5',
    },
]
</script>

<script setup lang="ts">
import { computed } from 'vue'
import { Plus } from 'lucide-vue-next'

interface Props {
    sessions?: ChatSession[]
    activeId?: string
}

const props = withDefaults(defineProps<Props>(), {
    sessions: () => defaultSessions,
    activeId: undefined,
})

const emit = defineEmits<{
    select: [id: string]
    new: []
}>()

// ── date grouping ────────────────────────────────────────────────────
type Group = { label: string, sessions: ChatSession[] }

const groups = computed<Group[]>(() => {
    const now = new Date()
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const startOfYesterday = new Date(startOfToday.getTime() - 86400000)
    const startOf7Days = new Date(startOfToday.getTime() - 6 * 86400000)
    const startOf30Days = new Date(startOfToday.getTime() - 29 * 86400000)

    const buckets: Record<string, ChatSession[]> = {
        Today: [],
        Yesterday: [],
        'Previous 7 Days': [],
        'Previous 30 Days': [],
        Older: [],
    }

    for (const s of props.sessions) {
        const t = s.updatedAt.getTime()
        if (t >= startOfToday.getTime())
            buckets['Today'].push(s)
        else if (t >= startOfYesterday.getTime())
            buckets['Yesterday'].push(s)
        else if (t >= startOf7Days.getTime())
            buckets['Previous 7 Days'].push(s)
        else if (t >= startOf30Days.getTime())
            buckets['Previous 30 Days'].push(s)
        else
            buckets['Older'].push(s)
    }

    return Object.entries(buckets)
        .filter(([, items]) => items.length > 0)
        .map(([label, sessions]) => ({ label, sessions }))
})

function formatTime(date: Date): string {
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

        <!-- grouped session list -->
        <template v-for="group in groups" :key="group.label">
            <!-- section label -->
            <p class="mt-2 mb-0.5 px-2 text-[11px] font-semibold uppercase tracking-wide text-[var(--hig-tertiary-label,rgba(60,60,67,0.3))] dark:text-[rgba(235,235,245,0.3)] first:mt-0">
                {{ group.label }}
            </p>

            <!-- session rows -->
            <button
                v-for="session in group.sessions"
                :key="session.id"
                class="group relative flex w-full flex-col items-start rounded-lg px-2 py-1.5 text-left transition-colors"
                :class="activeId === session.id
                    ? 'bg-[#0077ED] text-white'
                    : 'text-[var(--hig-label)] hover:bg-[rgba(0,0,0,0.06)] dark:hover:bg-[rgba(255,255,255,0.06)]'"
                @click="emit('select', session.id)"
            >
                <div class="flex w-full items-baseline gap-1">
                    <!-- title -->
                    <span
                        class="flex-1 truncate text-[13px] font-medium leading-snug"
                        :class="activeId !== session.id && 'text-[var(--hig-label)]'"
                    >
                        {{ session.title }}
                    </span>
                    <!-- timestamp -->
                    <span
                        class="shrink-0 text-[11px] tabular-nums"
                        :class="activeId === session.id
                            ? 'text-white/70'
                            : 'text-[var(--hig-tertiary-label,rgba(60,60,67,0.3))] dark:text-[rgba(235,235,245,0.3)]'"
                    >
                        {{ formatTime(session.updatedAt) }}
                    </span>
                </div>

                <!-- preview -->
                <span
                    v-if="session.preview"
                    class="mt-0.5 line-clamp-2 w-full text-[12px] leading-snug"
                    :class="activeId === session.id
                        ? 'text-white/75'
                        : 'text-[var(--hig-secondary-label)]'"
                >
                    {{ session.preview }}
                </span>
            </button>
        </template>
    </div>
</template>
