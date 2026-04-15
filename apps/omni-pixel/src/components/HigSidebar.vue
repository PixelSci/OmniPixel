<script setup lang="ts">
import type { Component, HTMLAttributes } from 'vue'
import {
    Film,
    Image,
    Layers,
    Music,
    Settings,
    Sparkles,
    Type,
    PanelLeft,
} from 'lucide-vue-next'
import { cn } from '@/lib/utils'

interface SidebarItem {
    id: string
    label: string
    icon?: Component
    badge?: string | number
    active?: boolean
}

interface SidebarSection {
    title?: string
    items: SidebarItem[]
}

interface Props {
    title?: string
    subtitle?: string
    searchPlaceholder?: string
    sections?: SidebarSection[]
    class?: HTMLAttributes['class']
}

const props = withDefaults(defineProps<Props>(), {
    title: 'OmniPixel',
    subtitle: 'Creative Suite',
    searchPlaceholder: 'Search',
    sections: () => [
        {
            items: [
                { id: 'text', label: 'Text', icon: Type, active: true },
                { id: 'image', label: 'Image', icon: Image },
                { id: 'audio', label: 'Audio', icon: Music },
                { id: 'video', label: 'Video', icon: Film },
            ],
        },
        {
            title: 'Library',
            items: [
                { id: 'layers', label: 'Layers', icon: Layers, badge: 12 },
                { id: 'effects', label: 'Effects', icon: Sparkles },
                { id: 'settings', label: 'Settings', icon: Settings },
            ],
        },
    ],
})
</script>

<template>
    <aside
        :class="cn(
            'relative isolate flex h-full min-h-[420px] w-[254px] flex-col overflow-hidden rounded-2xl border-r border-white/60 text-[#1d1d1f] shadow-[inset_-1px_0_0_rgba(255,255,255,0.65),0_0_0_0.5px_rgba(0,0,0,0.04),0_28px_64px_-12px_rgba(0,0,0,0.18)] backdrop-blur-[44px] backdrop-saturate-150 dark:border-white/14 dark:text-[rgba(245,245,247,0.94)] dark:shadow-[inset_-1px_0_0_rgba(255,255,255,0.10),0_0_0_0.5px_rgba(255,255,255,0.06),0_28px_64px_-12px_rgba(0,0,0,0.55)]',
            'bg-[rgba(252,252,254,0.52)] dark:bg-[rgba(34,34,36,0.55)]',
            props.class,
        )"
    >
        <div class="pointer-events-none absolute inset-0 -z-10 bg-white/45 mix-blend-color-dodge dark:hidden" aria-hidden="true" />
        <div class="pointer-events-none absolute inset-0 -z-10 bg-[linear-gradient(105deg,rgba(255,255,255,0.62)_0%,rgba(255,255,255,0.22)_42%,rgba(255,255,255,0.06)_72%,rgba(255,255,255,0.02)_100%)] dark:bg-[linear-gradient(105deg,rgba(255,255,255,0.14)_0%,rgba(255,255,255,0.05)_45%,rgba(255,255,255,0.02)_100%)]" aria-hidden="true" />
        <div class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(120%_80%_at_0%_0%,rgba(255,255,255,0.55),transparent_55%)] opacity-80 dark:bg-[radial-gradient(120%_80%_at_0%_0%,rgba(255,255,255,0.12),transparent_55%)] dark:opacity-100" aria-hidden="true" />

        <header class="flex shrink-0 justify-between items-center gap-3 px-3 pb-3 pt-3">
            <HigTrafficLights />
            <PanelLeft class="size-4" />
        </header>
    </aside>
</template>
