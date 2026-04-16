<script setup lang="ts">
interface Props {
    content: string
    model?: string
}

defineProps<Props>()
</script>

<template>
    <div>
        <!-- article card -->
        <div
            class="rounded-[14px] border border-[rgba(0,0,0,0.06)] bg-white/60 px-5 py-4 backdrop-blur-[20px] backdrop-saturate-150 dark:border-[rgba(255,255,255,0.08)] dark:bg-[rgba(44,44,46,0.48)]"
            style="box-shadow: 0 1px 0 rgba(255,255,255,0.72) inset, 0 2px 12px rgba(0,0,0,0.06);"
        >
            <!-- header: avatar + model -->
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

            <!-- separator -->
            <div class="mb-3.5 h-px bg-[rgba(0,0,0,0.06)] dark:bg-[rgba(255,255,255,0.07)]" />

            <!-- prose -->
            <div class="text-[14px] leading-[22px] text-[var(--hig-label)]">
                <p class="whitespace-pre-wrap">{{ content }}</p>
            </div>
        </div>
    </div>
</template>
