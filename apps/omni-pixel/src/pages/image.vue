<script setup lang="ts">
import { computed, ref, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
    Sparkles,
    Images,
    Scissors,
    Zap,
    Palette,
    Paintbrush,
    ArrowUp,
    ChevronDown,
} from 'lucide-vue-next'
import { useModelSettings } from '@/composables/useModelSettings'
import { Button } from '@/components/ui/button'
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from '@/components/ui/card'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuLabel,
    DropdownMenuRadioGroup,
    DropdownMenuRadioItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const route = useRoute()
const router = useRouter()

const isChildActive = () => route.matched.length > 1

const { chatModels } = useModelSettings()
const selectedModelId = ref(chatModels.value[0]?.id)
const selectedModel = computed(() => chatModels.value.find(m => m.id === selectedModelId.value) || chatModels.value[0])
const modelGroups = computed(() => {
    const groups = new Map<string, typeof chatModels.value>()
    for (const model of chatModels.value) {
        const models = groups.get(model.providerName) || []
        models.push(model)
        groups.set(model.providerName, models)
    }
    return [...groups.entries()].map(([label, models]) => ({ label, models }))
})

watchEffect(() => {
    if (!selectedModel.value && chatModels.value[0])
        selectedModelId.value = chatModels.value[0].id
})

const prompt = ref('')
const canSend = computed(() => prompt.value.trim().length > 0 && !!selectedModel.value)

function send() {
    if (!canSend.value || !selectedModel.value) return
    router.push({
        path: '/image/text-to-image',
        query: { prompt: prompt.value.trim(), model: selectedModel.value.id },
    })
}

function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
        e.preventDefault()
        send()
    }
}

const features = [
    { title: 'Text to Image', description: 'Generate images from text descriptions', icon: Sparkles, path: '/image/text-to-image' },
    { title: 'Image to Image', description: 'Transform existing images with prompts', icon: Images, path: '/image/image-to-image' },
    { title: 'Remove Background', description: 'Remove image backgrounds instantly', icon: Scissors, path: '/image/remove-bg' },
    { title: 'Enhance', description: 'Upscale and enhance image quality', icon: Zap, path: '/image/enhance' },
    { title: 'Style Transfer', description: 'Apply artistic styles to your images', icon: Palette, path: '/image/style-transfer' },
    { title: 'Inpaint', description: 'Edit specific areas of an image', icon: Paintbrush, path: '/image/inpaint' },
]

function goCard(path: string) {
    router.push(path)
}
</script>

<template>
    <RouterView v-if="isChildActive()" />
    <div v-else class="flex h-full flex-col items-center justify-center px-8">
        <h1 class="mb-8 text-[34px] font-semibold tracking-tight text-[var(--hig-label)]">
            Create Images with AI
        </h1>

        <!-- input area -->
        <div
            class="hig-liquid-glass flex w-full max-w-[800px] flex-col rounded-[20px] border border-[rgba(0,0,0,0.08)] bg-white/70 backdrop-blur-[28px] backdrop-saturate-180 dark:border-[rgba(255,255,255,0.12)] dark:bg-[rgba(44,44,46,0.62)]"
        >
            <div class="px-5 pt-4 pb-2">
                <textarea
                    v-model="prompt"
                    placeholder="Describe the image you want to create..."
                    rows="2"
                    class="max-h-[200px] min-h-[48px] w-full resize-none bg-transparent text-[16px] leading-7 text-[var(--hig-label)] placeholder:text-[var(--hig-tertiary-label)] focus:outline-none"
                    @keydown="handleKeydown"
                />
            </div>
            <div class="flex items-center gap-2 border-t border-[rgba(0,0,0,0.06)] px-4 py-3 dark:border-[rgba(255,255,255,0.06)]">
                <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                        <Button
                            variant="ghost"
                            aria-label="Select model"
                            class="h-[32px] gap-1.5 rounded-lg px-3 text-[13px] font-medium text-[var(--hig-secondary-label)] data-[state=open]:bg-accent"
                        >
                            <span class="max-w-[160px] truncate">{{ selectedModel?.name ?? 'No model' }}</span>
                            <ChevronDown class="size-[12px]! opacity-60" :stroke-width="2" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent side="top" align="start" :side-offset="8" class="min-w-[200px]">
                        <DropdownMenuRadioGroup v-model="selectedModelId">
                            <template v-for="(group, gi) in modelGroups" :key="group.label">
                                <DropdownMenuSeparator v-if="gi > 0" />
                                <DropdownMenuLabel class="text-[11px] text-[var(--hig-tertiary-label)]">
                                    {{ group.label }}
                                </DropdownMenuLabel>
                                <DropdownMenuRadioItem
                                    v-for="model in group.models"
                                    :key="model.id"
                                    :value="model.id"
                                    class="text-[13px]"
                                >
                                    <span class="flex-1">{{ model.name }}</span>
                                </DropdownMenuRadioItem>
                            </template>
                        </DropdownMenuRadioGroup>
                    </DropdownMenuContent>
                </DropdownMenu>

                <div class="flex-1" />

                <Button
                    variant="ghost"
                    size="icon-sm"
                    :disabled="!canSend"
                    aria-label="Generate"
                    class="size-[34px] rounded-full active:scale-95"
                    :class="canSend
                        ? 'bg-[#0077ED] text-white hover:bg-[#0068d1]'
                        : 'bg-[var(--hig-quaternary-fill)] text-[var(--hig-tertiary-label)]'"
                    @click="send"
                >
                    <ArrowUp class="size-[16px]!" :stroke-width="2.2" />
                </Button>
            </div>
        </div>

        <!-- cards -->
        <div class="mt-10 grid w-full max-w-[880px] grid-cols-3 gap-4">
            <Card
                v-for="feature in features"
                :key="feature.title"
                class="cursor-pointer transition-all duration-150 hover:scale-[1.02] hover:shadow-md"
                @click="goCard(feature.path)"
            >
                <CardHeader class="pb-2">
                    <component :is="feature.icon" class="mb-1 size-6 text-[var(--hig-secondary-label)]" :stroke-width="1.8" />
                    <CardTitle class="text-[16px]">{{ feature.title }}</CardTitle>
                </CardHeader>
                <CardContent>
                    <CardDescription class="text-[13px]">{{ feature.description }}</CardDescription>
                </CardContent>
            </Card>
        </div>
    </div>
</template>
