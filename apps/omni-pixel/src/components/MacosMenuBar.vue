<script setup lang="ts">
import { useDark, useNow, useToggle } from '@vueuse/core'
import { Moon, Sun, Volume2 } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import HigBox from '@/components/HigBox.vue'
import { Button } from '@/components/ui/button'

interface Props {
    appName?: string
}

const props = withDefaults(defineProps<Props>(), {
    appName: 'OmniPixel',
})

const isDark = useDark()
const toggleDark = useToggle(isDark)

const now = useNow({ interval: 1000 })

const timeStr = computed(() =>
    now.value.toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
    }),
)

const dateStr = computed(() =>
    now.value.toLocaleDateString('zh-CN', {
        month: 'long',
        day: 'numeric',
        weekday: 'short',
    }),
)

// ── menu bar nav items ──────────────────────────────────────────────────
const menuItems = [
    { label: 'Text' },
    { label: 'Image' },
    { label: 'Audio' },
    { label: 'Video' },
]

// ── Image submenu ───────────────────────────────────────────────────────
type SubItem = { type: 'item', label: string } | { type: 'separator' }

const imageSubItems: SubItem[] = [
    { type: 'item', label: 'Import Image...' },
    { type: 'item', label: 'Export Image...' },
    { type: 'separator' },
    { type: 'item', label: 'Resize...' },
    { type: 'item', label: 'Crop' },
    { type: 'item', label: 'Rotate Left' },
    { type: 'item', label: 'Rotate Right' },
    { type: 'separator' },
    { type: 'item', label: 'Convert Format...' },
]

const imageMenuOpen = ref(false)
const menuBarEl = ref<HTMLElement | null>(null)
const imageMenuLeft = ref(0)
const menuBarHeight = ref(32)
let closeTimer: ReturnType<typeof setTimeout> | null = null

const imageMenuStyle = computed(() => ({
    left: `${imageMenuLeft.value}px`,
    top: `1px`,
}))

function updateImageMenuPosition(target: HTMLElement) {
    const menuBarRect = menuBarEl.value?.getBoundingClientRect()
    const targetRect = target.getBoundingClientRect()

    menuBarHeight.value = menuBarRect?.height ?? menuBarHeight.value
    imageMenuLeft.value = menuBarRect ? targetRect.left - menuBarRect.left : target.offsetLeft
}

function openImageMenu(event?: MouseEvent) {
    keepImageMenuOpen()

    if (event?.currentTarget instanceof HTMLElement) {
        updateImageMenuPosition(event.currentTarget)
    }

    imageMenuOpen.value = true
}

function keepImageMenuOpen() {
    if (closeTimer) {
        clearTimeout(closeTimer)
        closeTimer = null
    }
}

function scheduleCloseImageMenu() {
    closeTimer = setTimeout(() => {
        imageMenuOpen.value = false
    }, 80)
}

function closeMenus() {
    if (closeTimer) {
        clearTimeout(closeTimer)
        closeTimer = null
    }

    imageMenuOpen.value = false
}

function handleDocumentPointerDown(event: PointerEvent) {
    if (event.target instanceof Node && menuBarEl.value?.contains(event.target))
        return

    closeMenus()
}

onMounted(() => {
    document.addEventListener('pointerdown', handleDocumentPointerDown)
})

onBeforeUnmount(() => {
    if (closeTimer)
        clearTimeout(closeTimer)
    document.removeEventListener('pointerdown', handleDocumentPointerDown)
})
</script>

<template>
    <div
        ref="menuBarEl"
        class="relative z-[100] h-8 w-full overflow-visible select-none"
    >
        <div
            class="hig-glass relative flex h-8 w-full items-center border-b px-1 [box-shadow:inset_0_1px_0_var(--hig-glass-highlight)]"
            style="z-index: 100;"
        >
            <!-- Left: Apple logo + app name + nav menus -->
            <div class="flex flex-1 items-center">
                <!-- Apple logo -->
                <Button
                    variant="ghost"
                    aria-label="Apple menu"
                    class="inline-flex h-[22px] cursor-default items-center justify-center rounded px-2.5 text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    <svg width="13" height="16" viewBox="0 0 814 1000" fill="currentColor">
                        <path d="M788.1 340.9c-5.8 4.5-108.2 62.2-108.2 190.5 0 148.4 130.3 200.9 134.2 202.2-.6 3.2-20.7 71.9-68.7 141.9-42.8 61.6-87.5 123.1-155.5 123.1s-85.5-39.5-164-39.5c-76 0-103.7 40.8-165.9 40.8s-105-37.5-155.5-127.4C46.3 747.6.5 641.2.5 539.9c0-182.4 116.5-278.6 231.3-278.6 63.3 0 115.9 41.4 155 41.4 37.5 0 96.9-43.7 165-43.7 26.7 0 108.2 2.6 166.7 98.5zm-209.5-222c30.7-36.5 52.4-87 52.4-137.5 0-7.1-.6-14.3-1.9-20.1-49.7 1.9-107.7 33.1-142.5 75.7-27.5 32.4-53.2 83.5-53.2 134.7 0 7.7 1.3 15.4 1.9 17.9 3.2.6 8.4 1.3 13.6 1.3 44.9 0 96.8-29.8 129.7-72z" />
                    </svg>
                </Button>

                <!-- App name -->
                <Button
                    variant="ghost"
                    class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] font-semibold tracking-tight text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    {{ props.appName }}
                </Button>

                <!-- Nav menu items -->
                <template v-for="item in menuItems" :key="item.label">
                    <!-- Image: drop-down submenu -->
                    <div
                        v-if="item.label === 'Image'"
                        class="flex h-8 items-center"
                        @mouseenter="openImageMenu"
                        @mouseleave="scheduleCloseImageMenu"
                    >
                        <Button
                            variant="ghost"
                            class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                            :class="imageMenuOpen && 'bg-[var(--hig-fill)]'"
                        >
                            Image
                        </Button>
                    </div>

                    <!-- Regular menu items -->
                    <Button
                        v-else
                        variant="ghost"
                        class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                    >
                        {{ item.label }}
                    </Button>
                </template>
            </div>

            <!-- Center: Dynamic Island -->
            <div class="pointer-events-none absolute top-0 left-1/2 flex h-7 -translate-x-1/2 items-center">
                <div
                    class="pointer-events-auto flex min-w-20 items-center justify-center gap-[5px] rounded-full px-3.5 h-[22px] bg-foreground shadow-[0_1px_3px_rgba(0,0,0,0.28),inset_0_1px_0_rgba(255,255,255,0.10)] transition-[min-width] duration-300 ease-[var(--ease-default)]"
                >
                    <slot>
                        <span class="block size-[5px] rounded-full bg-[var(--hig-secondary-label)] opacity-50" />
                        <span class="block size-[5px] rounded-full bg-[var(--hig-secondary-label)] opacity-50" />
                    </slot>
                </div>
            </div>

            <!-- Right: Status menu -->
            <div class="flex flex-1 items-center justify-end">
                <!-- Theme toggle -->
                <Button
                    variant="ghost"
                    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
                    class="inline-flex h-[22px] w-7 cursor-default items-center justify-center rounded text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                    @click="toggleDark()"
                >
                    <Sun v-if="isDark" :size="13" :stroke-width="1.8" />
                    <Moon v-else :size="13" :stroke-width="1.8" />
                </Button>

                <!-- Volume -->
                <Button
                    variant="ghost"
                    aria-label="Volume"
                    class="inline-flex h-[22px] w-7 cursor-default items-center justify-center rounded text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    <Volume2 :size="13" :stroke-width="1.8" />
                </Button>

                <!-- Clock -->
                <div
                    class="flex h-[22px] cursor-default items-center gap-1 rounded px-1.5 text-[13px] transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    <span class="text-[var(--hig-secondary-label)]">{{ dateStr }}</span>
                    <span class="font-medium tabular-nums text-foreground">{{ timeStr }}</span>
                </div>
            </div>
        </div>

        <HigBox
            v-if="imageMenuOpen"
            :padding="false"
            class="absolute z-[1000] min-w-[180px]"
            :style="imageMenuStyle"
            @mouseenter="keepImageMenuOpen"
            @mouseleave="scheduleCloseImageMenu"
        >
            <!-- px-[12px] py-[5px] matches Figma menu inset -->
            <div class="px-3 py-[5px]">
                <template v-for="(sub, i) in imageSubItems" :key="i">
                    <!-- separator: #e6e6e6 = Figma vibrant secondary fill -->
                    <div
                        v-if="sub.type === 'separator'"
                        class="-mx-3 flex h-[11px] items-center"
                    >
                        <div class="h-px w-full bg-[#e6e6e6] dark:bg-[rgba(84,84,88,0.65)]" />
                    </div>

                    <!-- item -->
                    <button
                        v-else
                        class="group relative flex h-[24px] w-full cursor-default items-center outline-none"
                    >
                        <!--
                            Hover BG — three-layer composite from Figma:
                            1. rgba(0,0,0,0.05)  subtle darkening
                            2. rgba(255,255,255,0.65) color-dodge  vibrancy brightening
                            3. #0088ff  accent blue on top
                            Extends 7px past item edges to reach menu padding boundary
                        -->
                        <span
                            class="absolute inset-y-0 left-[-7px] right-[-7px] rounded-[8px] opacity-0 group-hover:opacity-100"
                            aria-hidden="true"
                        >
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(0,0,0,0.05)]" />
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(255,255,255,0.65)] mix-blend-color-dodge" />
                            <span class="absolute inset-0 rounded-[8px] bg-[#0088ff]" />
                        </span>
                        <!-- #1a1a1a = Figma vibrant primary label (not pure black) -->
                        <span class="relative text-[13px] font-[510] text-[#1a1a1a] group-hover:text-white dark:text-[rgba(235,235,245,0.9)] dark:group-hover:text-white">
                            {{ sub.label }}
                        </span>
                    </button>
                </template>
            </div>
        </HigBox>
    </div>
</template>
