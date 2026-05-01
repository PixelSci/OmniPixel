<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import HigSheet from '@/components/HigSheet.vue'
import { Checkbox } from '@/components/ui/checkbox'
import { useModelSettings } from '@/composables/useModelSettings'

interface Props {
    open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const openModel = useVModel(props, 'open', emit, { passive: true })
const { sections, isModelEnabled, setModelEnabled } = useModelSettings()
</script>

<template>
    <HigSheet v-model:open="openModel" title="模型设置" width="620px">
        <div class="flex flex-col gap-5 px-5 py-5">
            <section
                v-for="section in sections"
                :key="section.key"
                class="flex flex-col gap-3"
            >
                <div class="flex items-center justify-between gap-4">
                    <div class="flex flex-col">
                        <span class="text-[13px] font-[510] text-foreground">{{ section.title }}</span>
                        <span class="text-[12px] text-[var(--hig-secondary-label)]">{{ section.description }}</span>
                    </div>
                    <span class="text-[12px] text-[var(--hig-secondary-label)]">
                        {{ section.models.filter(model => isModelEnabled(section.key, model.id)).length }}/{{ section.models.length }}
                    </span>
                </div>

                <div class="grid grid-cols-2 gap-2">
                    <label
                        v-for="model in section.models"
                        :key="model.id"
                        class="flex min-h-10 cursor-default items-center gap-2 rounded-lg border border-[rgba(0,0,0,0.08)] px-3 py-2 text-[12px] transition-colors hover:bg-[var(--hig-fill)] dark:border-[rgba(255,255,255,0.1)]"
                    >
                        <Checkbox
                            :model-value="isModelEnabled(section.key, model.id)"
                            @update:model-value="value => setModelEnabled(section.key, model.id, value === true)"
                        />
                        <span class="flex min-w-0 flex-1 flex-col">
                            <span class="truncate font-medium text-foreground">{{ model.name }}</span>
                            <span class="truncate text-[11px] text-[var(--hig-secondary-label)]">{{ model.provider }}</span>
                        </span>
                        <span
                            v-if="model.badge"
                            class="rounded px-1 py-0.5 text-[10px] font-medium leading-none"
                            :class="model.badge === 'Fast'
                                ? 'bg-[rgba(52,199,89,0.12)] text-[#34c759] dark:bg-[rgba(48,209,88,0.18)] dark:text-[#30d158]'
                                : 'bg-[rgba(0,119,237,0.10)] text-[#0077ED] dark:bg-[rgba(0,149,255,0.18)] dark:text-[#3aa3ff]'"
                        >
                            {{ model.badge }}
                        </span>
                    </label>
                </div>
            </section>
        </div>
    </HigSheet>
</template>
