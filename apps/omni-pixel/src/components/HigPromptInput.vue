<script setup lang="ts">
import { useTextareaAutosize } from '@vueuse/core'
import { ArrowUp, ChevronDown, Globe, Mic, Paperclip, Square } from 'lucide-vue-next'
import { computed, nextTick, ref } from 'vue'
import { Button } from '@/components/ui/button'

interface Props {
    placeholder?: string
    disabled?: boolean
    loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    placeholder: 'Message…',
    disabled: false,
    loading: false,
})

const emit = defineEmits<{
    send: [message: string]
    stop: []
    attach: []
}>()

const { textarea, input: text, triggerResize } = useTextareaAutosize({
    styleProp: 'height',
})

interface Model {
    id: string
    name: string
    provider: string
    badge?: string
}

const modelGroups: { label: string, models: Model[] }[] = [
    {
        label: 'Anthropic',
        models: [
            { id: 'claude-sonnet-4', name: 'Claude Sonnet 4', provider: 'anthropic', badge: 'Smart' },
            { id: 'claude-haiku-4', name: 'Claude Haiku 4', provider: 'anthropic', badge: 'Fast' },
        ],
    },
    {
        label: 'OpenAI',
        models: [
            { id: 'gpt-4o', name: 'GPT-4o', provider: 'openai' },
            { id: 'gpt-4o-mini', name: 'GPT-4o mini', provider: 'openai', badge: 'Fast' },
        ],
    },
]

const allModels = modelGroups.flatMap(g => g.models)
const selectedModelId = ref('claude-sonnet-4')
const selectedModel = computed(() => allModels.find(m => m.id === selectedModelId.value)!)

const webSearchEnabled = ref(false)
const isRecording = ref(false)

const canSend = computed(() => text.value.trim().length > 0 && !props.disabled && !props.loading)

function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
        e.preventDefault()
        send()
    }
}

function send() {
    if (!canSend.value)
        return
    emit('send', text.value.trim())
    text.value = ''
    nextTick(() => triggerResize())
}

function toggleWebSearch() {
    webSearchEnabled.value = !webSearchEnabled.value
}

function toggleVoice() {
    isRecording.value = !isRecording.value
}
</script>

<template>
    <div class="w-full px-4 pb-4 pt-2">
        <div class="mx-auto w-full max-w-1/2">
            <!-- input pill -->
            <div
                class="hig-liquid-glass relative flex flex-col rounded-[18px] border border-[rgba(0,0,0,0.08)] bg-white/70 backdrop-blur-[28px] backdrop-saturate-180 dark:border-[rgba(255,255,255,0.12)] dark:bg-[rgba(44,44,46,0.62)]"
            >
                <!-- textarea -->
                <div class="px-4 pt-3 pb-1">
                    <textarea
                        ref="textarea"
                        v-model="text"
                        :placeholder="placeholder"
                        :disabled="disabled"
                        rows="1"
                        class="max-h-[200px] min-h-[24px] w-full resize-none bg-transparent text-[14px] leading-6 text-[var(--hig-label)] placeholder:text-[var(--hig-tertiary-label)] focus:outline-none disabled:opacity-50"
                        @keydown="handleKeydown"
                    />
                </div>

                <!-- toolbar -->
                <div class="flex items-center gap-0.5 border-t border-[rgba(0,0,0,0.06)] px-2.5 py-2 dark:border-[rgba(255,255,255,0.06)]">
                    <!-- model selector -->
                    <UiDropdownMenuDropdownMenu>
                        <UiDropdownMenuDropdownMenuTrigger as-child>
                            <Button
                                variant="ghost"
                                size="sm"
                                :disabled="disabled"
                                aria-label="Select model"
                                class="h-[30px] gap-1 rounded-lg px-2 text-[12px] font-medium text-[var(--hig-secondary-label)] data-[state=open]:bg-accent"
                            >
                                <span class="max-w-[120px] truncate">{{ selectedModel.name }}</span>
                                <ChevronDown class="size-[11px]! opacity-60" :stroke-width="2" />
                            </Button>
                        </UiDropdownMenuDropdownMenuTrigger>

                        <UiDropdownMenuDropdownMenuContent
                            side="top"
                            align="start"
                            :side-offset="6"
                            class="min-w-[180px]"
                        >
                            <UiDropdownMenuDropdownMenuRadioGroup v-model="selectedModelId">
                                <template v-for="(group, gi) in modelGroups" :key="group.label">
                                    <UiDropdownMenuDropdownMenuSeparator v-if="gi > 0" />
                                    <UiDropdownMenuDropdownMenuLabel class="text-[11px] text-[var(--hig-tertiary-label)]">
                                        {{ group.label }}
                                    </UiDropdownMenuDropdownMenuLabel>
                                    <UiDropdownMenuDropdownMenuRadioItem
                                        v-for="model in group.models"
                                        :key="model.id"
                                        :value="model.id"
                                        class="text-[13px]"
                                    >
                                        <span class="flex-1">{{ model.name }}</span>
                                        <span
                                            v-if="model.badge"
                                            class="ml-2 rounded px-1 py-0.5 text-[10px] font-medium leading-none"
                                            :class="model.badge === 'Fast'
                                                ? 'bg-[rgba(52,199,89,0.12)] text-[#34c759] dark:bg-[rgba(48,209,88,0.18)] dark:text-[#30d158]'
                                                : 'bg-[rgba(0,119,237,0.10)] text-[#0077ED] dark:bg-[rgba(0,149,255,0.18)] dark:text-[#3aa3ff]'"
                                        >
                                            {{ model.badge }}
                                        </span>
                                    </UiDropdownMenuDropdownMenuRadioItem>
                                </template>
                            </UiDropdownMenuDropdownMenuRadioGroup>
                        </UiDropdownMenuDropdownMenuContent>
                    </UiDropdownMenuDropdownMenu>

                    <div class="mx-1 h-[14px] w-px bg-[var(--hig-separator)] opacity-40" />

                    <!-- attach -->
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        :disabled="disabled"
                        aria-label="Attach file"
                        class="size-[30px] rounded-lg text-[var(--hig-secondary-label)]"
                        @click="emit('attach')"
                    >
                        <Paperclip class="size-[15px]!" :stroke-width="1.7" />
                    </Button>

                    <!-- web search toggle -->
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        :disabled="disabled"
                        :aria-pressed="webSearchEnabled"
                        aria-label="Toggle web search"
                        class="size-[30px] rounded-lg"
                        :class="webSearchEnabled
                            ? 'bg-[rgba(0,119,237,0.10)] text-[#0077ED] hover:bg-[rgba(0,119,237,0.15)] dark:bg-[rgba(0,149,255,0.18)] dark:text-[#3aa3ff] dark:hover:bg-[rgba(0,149,255,0.24)]'
                            : 'text-[var(--hig-secondary-label)]'"
                        @click="toggleWebSearch"
                    >
                        <Globe class="size-[15px]!" :stroke-width="1.7" />
                    </Button>

                    <!-- voice input -->
                    <Button
                        variant="ghost"
                        size="icon-sm"
                        :disabled="disabled"
                        :aria-pressed="isRecording"
                        aria-label="Voice input"
                        class="size-[30px] rounded-lg"
                        :class="isRecording
                            ? 'bg-[rgba(255,59,48,0.10)] text-[#ff3b30] hover:bg-[rgba(255,59,48,0.15)] dark:bg-[rgba(255,69,58,0.18)] dark:text-[#ff6961] dark:hover:bg-[rgba(255,69,58,0.24)]'
                            : 'text-[var(--hig-secondary-label)]'"
                        @click="toggleVoice"
                    >
                        <Mic class="size-[15px]!" :stroke-width="1.7" />
                    </Button>

                    <div class="flex-1" />

                    <!-- stop -->
                    <Button
                        v-if="loading"
                        variant="ghost"
                        size="icon-sm"
                        aria-label="Stop generating"
                        class="size-[30px] rounded-full bg-[var(--hig-label)] text-[var(--hig-background)] hover:bg-[var(--hig-label)] hover:opacity-80 active:scale-95"
                        @click="emit('stop')"
                    >
                        <Square class="size-[11px]! fill-current" :stroke-width="0" />
                    </Button>

                    <!-- send -->
                    <Button
                        v-else
                        variant="ghost"
                        size="icon-sm"
                        :disabled="!canSend"
                        aria-label="Send message"
                        class="size-[30px] rounded-full active:scale-95"
                        :class="canSend
                            ? 'bg-[#0077ED] text-white hover:bg-[#0068d1]'
                            : 'bg-[var(--hig-quaternary-fill)] text-[var(--hig-tertiary-label)]'"
                        @click="send"
                    >
                        <ArrowUp class="size-[15px]!" :stroke-width="2.2" />
                    </Button>
                </div>
            </div>
        </div>
    </div>
</template>
