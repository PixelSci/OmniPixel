<script setup lang="ts">
import { useDark } from '@vueuse/core'
import MarkdownRender from 'markstream-vue'
import 'markstream-vue/index.css'

interface Props {
    content: string
    streaming?: boolean
    model?: string
}

const props = withDefaults(defineProps<Props>(), {
    streaming: false,
})

const isDark = useDark()
</script>

<template>
    <div>
        <div
            class="rounded-[14px] border border-[rgba(0,0,0,0.06)] bg-white/60 px-5 py-4 backdrop-blur-[20px] backdrop-saturate-150 dark:border-[rgba(255,255,255,0.08)] dark:bg-[rgba(44,44,46,0.48)]"
            style="box-shadow: 0 1px 0 rgba(255,255,255,0.72) inset, 0 2px 12px rgba(0,0,0,0.06);"
        >
            <div class="mb-3 flex items-center gap-2">
                <div class="flex size-[22px] shrink-0 items-center justify-center rounded-full bg-[linear-gradient(135deg,#0077ed,#9b51e0)] shadow-[0_1px_3px_rgba(0,0,0,0.18)]">
                    <svg width="11" height="11" viewBox="0 0 13 13" fill="none" aria-hidden="true">
                        <path
                            d="M6.5 1L7.72 4.78H11.7L8.49 7.07L9.71 10.85L6.5 8.56L3.29 10.85L4.51 7.07L1.3 4.78H5.28L6.5 1Z"
                            fill="white"
                            fill-opacity="0.95"
                        />
                    </svg>
                </div>
                <span class="text-[12px] font-medium text-[var(--hig-secondary-label)]">{{ model ?? 'Assistant' }}</span>
            </div>

            <div class="mb-3.5 h-px bg-[rgba(0,0,0,0.06)] dark:bg-[rgba(255,255,255,0.07)]" />

            <div class="text-[14px] leading-[22px] text-[var(--hig-label)]">
                <MarkdownRender
                    v-if="streaming"
                    :content="content"
                    :final="false"
                    :max-live-nodes="0"
                    :batch-rendering="true"
                    :render-batch-size="16"
                    :render-batch-delay="8"
                    :typewriter="true"
                    :is-dark="isDark"
                />
                <MarkdownRender
                    v-else
                    :content="content"
                    :final="true"
                    :max-live-nodes="320"
                    :is-dark="isDark"
                />
            </div>
        </div>
    </div>
</template>
