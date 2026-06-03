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
const { chatModels, isModelEnabled, setModelEnabled } = useModelSettings()
</script>

<template>
    <HigSheet v-model:open="openModel" title="模型设置" width="620px">
        <div class="flex flex-col gap-5 px-5 py-5">
            <div class="flex items-center justify-between">
                <span class="text-[13px] font-[510] text-foreground">Chat 模型</span>
                <span class="text-[12px] text-[var(--hig-secondary-label)]">
                    {{ chatModels.length }} 个可用
                </span>
            </div>

            <div class="grid grid-cols-2 gap-2">
                <label
                    v-for="model in chatModels"
                    :key="model.id"
                    class="flex min-h-10 cursor-default items-center gap-2 rounded-lg border border-[rgba(0,0,0,0.08)] px-3 py-2 text-[12px] transition-colors hover:bg-[var(--hig-fill)] dark:border-[rgba(255,255,255,0.1)]"
                >
                    <Checkbox
                        :model-value="isModelEnabled(model.id)"
                        @update:model-value="value => setModelEnabled(model.id, value === true)"
                    />
                    <span class="flex min-w-0 flex-1 flex-col">
                        <span class="truncate font-medium text-foreground">{{ model.name }}</span>
                        <span class="truncate text-[11px] text-[var(--hig-secondary-label)]">{{ model.providerName }}</span>
                    </span>
                </label>
            </div>
        </div>
    </HigSheet>
</template>
