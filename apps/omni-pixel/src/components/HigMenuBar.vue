<script setup lang="ts">
import { useDark, useNow } from '@vueuse/core'
import { Moon, Sun, Volume2 } from 'lucide-vue-next'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import logo from '@/assets/logo.svg'
import HigBox from '@/components/HigBox.vue'
import HigModelSettingsSheet from '@/components/HigModelSettingsSheet.vue'
import HigProfileSheet from '@/components/HigProfileSheet.vue'
import { Button } from '@/components/ui/button'

interface Props {
    appName?: string
}

const props = withDefaults(defineProps<Props>(), {
    appName: 'OmniPixel',
})

const isDark = useDark()
const isAnimating = ref(false)

async function toggleTheme(event: MouseEvent) {
    if (isAnimating.value) return

    const isAppearanceTransition = 'startViewTransition' in document
        && !window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (!isAppearanceTransition) {
        isDark.value = !isDark.value
        return
    }

    isAnimating.value = true

    const x = event.clientX
    const y = event.clientY
    const endRadius = Math.hypot(
        Math.max(x, innerWidth - x),
        Math.max(y, innerHeight - y),
    )

    const transition = document.startViewTransition(() => {
        isDark.value = !isDark.value
    })

    transition.ready.then(() => {
        const clipPath = [
            `circle(0px at ${x}px ${y}px)`,
            `circle(${endRadius}px at ${x}px ${y}px)`,
        ]
        document.documentElement.animate(
            {
                clipPath: isDark.value
                    ? [...clipPath].reverse()
                    : clipPath,
            },
            {
                duration: 400,
                easing: 'ease-out',
                fill: 'forwards',
                pseudoElement: isDark.value
                    ? '::view-transition-old(root)'
                    : '::view-transition-new(root)',
            },
        )
    })

    transition.finished.then(() => {
        isAnimating.value = false
    })
}

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

// ── Logo submenu ────────────────────────────────────────────────────────
type LogoAction = 'profile' | 'models' | 'logout'
type LogoSubItem
    = | { type: 'item', label: string, action: LogoAction }
        | { type: 'separator' }

const logoSubItems: LogoSubItem[] = [
    { type: 'item', label: '个人中心', action: 'profile' },
    { type: 'item', label: '模型设置', action: 'models' },
    { type: 'separator' },
    { type: 'item', label: '退出登录', action: 'logout' },
]

const profileOpen = ref(false)
const modelSettingsOpen = ref(false)

function handleLogoAction(action: LogoAction) {
    closeMenus()
    if (action === 'profile')
        profileOpen.value = true
    else if (action === 'models')
        modelSettingsOpen.value = true
}

const imageMenuOpen = ref(false)
const logoMenuOpen = ref(false)
const menuBarEl = ref<HTMLElement | null>(null)
const imageMenuLeft = ref(0)
const logoMenuLeft = ref(0)
const menuBarHeight = ref(32)
let closeTimer: ReturnType<typeof setTimeout> | null = null

const imageMenuStyle = computed(() => ({
    left: `${imageMenuLeft.value}px`,
    top: `1px`,
}))

const logoMenuStyle = computed(() => ({
    left: `${logoMenuLeft.value}px`,
    top: `1px`,
}))

function menuOffsetLeft(target: HTMLElement) {
    const menuBarRect = menuBarEl.value?.getBoundingClientRect()
    const targetRect = target.getBoundingClientRect()

    menuBarHeight.value = menuBarRect?.height ?? menuBarHeight.value
    return menuBarRect ? targetRect.left - menuBarRect.left : target.offsetLeft
}

function openImageMenu(event?: MouseEvent) {
    keepMenusOpen()

    if (event?.currentTarget instanceof HTMLElement) {
        imageMenuLeft.value = menuOffsetLeft(event.currentTarget)
    }

    logoMenuOpen.value = false
    imageMenuOpen.value = true
}

function openLogoMenu(event?: MouseEvent) {
    keepMenusOpen()

    if (event?.currentTarget instanceof HTMLElement) {
        logoMenuLeft.value = menuOffsetLeft(event.currentTarget)
    }

    imageMenuOpen.value = false
    logoMenuOpen.value = true
}

function keepMenusOpen() {
    if (closeTimer) {
        clearTimeout(closeTimer)
        closeTimer = null
    }
}

function scheduleCloseMenus() {
    closeTimer = setTimeout(() => {
        imageMenuOpen.value = false
        logoMenuOpen.value = false
    }, 80)
}

function closeMenus() {
    if (closeTimer) {
        clearTimeout(closeTimer)
        closeTimer = null
    }

    imageMenuOpen.value = false
    logoMenuOpen.value = false
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
            <div class="flex flex-1 items-center">
                <div
                    class="flex h-8 items-center"
                    @mouseenter="openLogoMenu"
                    @mouseleave="scheduleCloseMenus"
                >
                    <Button
                        variant="ghost"
                        aria-label="Apple menu"
                        class="inline-flex h-[22px] cursor-default items-center justify-center rounded px-2.5 text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                        :class="logoMenuOpen && 'bg-[var(--hig-fill)]'"
                    >
                        <img :src="logo" alt="logo" class="size-4">
                    </Button>
                </div>

                <Button
                    variant="ghost"
                    class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] font-semibold tracking-tight text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    {{ props.appName }}
                </Button>

                <template v-for="item in menuItems" :key="item.label">
                    <div
                        v-if="item.label === 'Image'"
                        class="flex h-8 items-center"
                        @mouseenter="openImageMenu"
                        @mouseleave="scheduleCloseMenus"
                    >
                        <Button
                            variant="ghost"
                            class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                            :class="imageMenuOpen && 'bg-[var(--hig-fill)]'"
                        >
                            Image
                        </Button>
                    </div>

                    <Button
                        v-else
                        variant="ghost"
                        class="inline-flex h-[22px] cursor-default items-center rounded px-2 text-[13px] text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                    >
                        {{ item.label }}
                    </Button>
                </template>
            </div>

            <div class="flex flex-1 items-center justify-end">
                <Button
                    variant="ghost"
                    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
                    class="inline-flex h-[22px] w-7 cursor-default items-center justify-center rounded text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                    @click="toggleTheme"
                >
                    <Sun v-if="isDark" :size="13" :stroke-width="1.8" />
                    <Moon v-else :size="13" :stroke-width="1.8" />
                </Button>

                <Button
                    variant="ghost"
                    aria-label="Volume"
                    class="inline-flex h-[22px] w-7 cursor-default items-center justify-center rounded text-foreground transition-colors duration-100 hover:bg-[var(--hig-fill)]"
                >
                    <Volume2 :size="13" :stroke-width="1.8" />
                </Button>

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
            @mouseenter="keepMenusOpen"
            @mouseleave="scheduleCloseMenus"
        >
            <div class="px-3 py-[5px]">
                <template v-for="(sub, i) in imageSubItems" :key="i">
                    <div
                        v-if="sub.type === 'separator'"
                        class="-mx-3 flex h-[11px] items-center"
                    >
                        <div class="h-px w-full bg-[#e6e6e6] dark:bg-[rgba(84,84,88,0.65)]" />
                    </div>

                    <button
                        v-else
                        class="group relative flex h-[24px] w-full cursor-default items-center outline-none"
                    >
                        <span
                            class="absolute inset-y-0 left-[-7px] right-[-7px] rounded-[8px] opacity-0 group-hover:opacity-100"
                            aria-hidden="true"
                        >
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(0,0,0,0.05)]" />
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(255,255,255,0.65)] mix-blend-color-dodge" />
                            <span class="absolute inset-0 rounded-[8px] bg-[#0088ff]" />
                        </span>
                        <span class="relative text-[13px] font-[510] text-[#1a1a1a] group-hover:text-white dark:text-[rgba(235,235,245,0.9)] dark:group-hover:text-white">
                            {{ sub.label }}
                        </span>
                    </button>
                </template>
            </div>
        </HigBox>

        <HigBox
            v-if="logoMenuOpen"
            :padding="false"
            class="absolute z-[1000] min-w-[160px]"
            :style="logoMenuStyle"
            @mouseenter="keepMenusOpen"
            @mouseleave="scheduleCloseMenus"
        >
            <div class="px-3 py-[5px]">
                <template v-for="(sub, i) in logoSubItems" :key="i">
                    <div
                        v-if="sub.type === 'separator'"
                        class="-mx-3 flex h-[11px] items-center"
                    >
                        <div class="h-px w-full bg-[#e6e6e6] dark:bg-[rgba(84,84,88,0.65)]" />
                    </div>

                    <button
                        v-else
                        class="group relative flex h-[24px] w-full cursor-default items-center outline-none"
                        @click="handleLogoAction(sub.action)"
                    >
                        <span
                            class="absolute inset-y-0 left-[-7px] right-[-7px] rounded-[8px] opacity-0 group-hover:opacity-100"
                            aria-hidden="true"
                        >
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(0,0,0,0.05)]" />
                            <span class="absolute inset-0 rounded-[8px] bg-[rgba(255,255,255,0.65)] mix-blend-color-dodge" />
                            <span class="absolute inset-0 rounded-[8px] bg-[#0088ff]" />
                        </span>
                        <span class="relative text-[13px] font-[510] text-[#1a1a1a] group-hover:text-white dark:text-[rgba(235,235,245,0.9)] dark:group-hover:text-white">
                            {{ sub.label }}
                        </span>
                    </button>
                </template>
            </div>
        </HigBox>

        <HigProfileSheet v-model:open="profileOpen" />
        <HigModelSettingsSheet v-model:open="modelSettingsOpen" />
    </div>
</template>
