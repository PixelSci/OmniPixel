<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { Eye, EyeOff, Plus, Trash2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import HigSheet from '@/components/HigSheet.vue'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Separator } from '@/components/ui/separator'

interface Props {
    open: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
    'update:open': [value: boolean]
}>()

const openModel = useVModel(props, 'open', emit, { passive: true })

const user = {
    name: 'Omni User',
    email: 'omni@pixel.app',
    plan: 'Pro',
}

const initials = user.name
    .split(' ')
    .map(part => part.charAt(0))
    .join('')
    .slice(0, 2)
    .toUpperCase()

const providerOptions = [
    { value: 'openai', label: 'OpenAI' },
    { value: 'anthropic', label: 'Anthropic' },
    { value: 'google', label: 'Google' },
    { value: 'deepseek', label: 'DeepSeek' },
    { value: 'other', label: '其他' },
]

interface ApiKeyEntry {
    id: string
    provider: string
    label: string
    key: string
    show: boolean
}

let idSeed = 0
function nextId() {
    idSeed += 1
    return `key-${idSeed}`
}

const keys = reactive<ApiKeyEntry[]>([
    { id: nextId(), provider: 'openai', label: '工作号', key: 'sk-************1234', show: false },
])

const draftProvider = ref(providerOptions[0].value)

function addKey() {
    keys.push({
        id: nextId(),
        provider: draftProvider.value,
        label: '',
        key: '',
        show: false,
    })
}

function removeKey(id: string) {
    const idx = keys.findIndex(k => k.id === id)
    if (idx !== -1)
        keys.splice(idx, 1)
}
</script>

<template>
    <HigSheet v-model:open="openModel" title="个人中心" width="560px">
        <div class="flex flex-col gap-5 px-5 py-5">
            <section class="flex items-center gap-4">
                <Avatar class="size-14 bg-[var(--hig-fill)]">
                    <AvatarFallback class="text-[15px] font-semibold text-[var(--hig-label)]">
                        {{ initials }}
                    </AvatarFallback>
                </Avatar>
                <div class="flex flex-col gap-0.5">
                    <span class="text-[15px] font-semibold tracking-tight text-foreground">{{ user.name }}</span>
                    <span class="text-[12px] text-[var(--hig-secondary-label)]">{{ user.email }}</span>
                    <span class="mt-1 inline-flex w-fit items-center rounded-full border border-[rgba(0,0,0,0.08)] bg-[var(--hig-fill)] px-2 py-[1px] text-[11px] font-medium text-[var(--hig-secondary-label)] dark:border-[rgba(255,255,255,0.1)]">
                        {{ user.plan }} 用户
                    </span>
                </div>
            </section>

            <Separator />

            <section class="flex flex-col gap-3">
                <div class="flex items-center justify-between">
                    <div class="flex flex-col">
                        <span class="text-[13px] font-[510] text-foreground">Model API Keys</span>
                        <span class="text-[12px] text-[var(--hig-secondary-label)]">管理你的模型服务凭证</span>
                    </div>
                    <Button
                        variant="outline"
                        size="sm"
                        class="h-7 text-[12px]"
                        @click="addKey"
                    >
                        <Plus :size="14" />
                        添加
                    </Button>
                </div>

                <div
                    v-if="keys.length === 0"
                    class="flex flex-col items-center justify-center gap-2 rounded-lg border border-dashed border-[rgba(0,0,0,0.1)] px-4 py-8 text-[12px] text-[var(--hig-secondary-label)] dark:border-[rgba(255,255,255,0.12)]"
                >
                    <span>尚未添加任何 API Key</span>
                    <Button variant="outline" size="sm" class="h-7 text-[12px]" @click="addKey">
                        <Plus :size="14" />
                        添加第一条
                    </Button>
                </div>

                <div v-else class="flex flex-col gap-2">
                    <div
                        v-for="entry in keys"
                        :key="entry.id"
                        class="flex items-center gap-2"
                    >
                        <NativeSelect v-model="entry.provider" class="w-[120px] shrink-0">
                            <NativeSelectOption
                                v-for="option in providerOptions"
                                :key="option.value"
                                :value="option.value"
                            >
                                {{ option.label }}
                            </NativeSelectOption>
                        </NativeSelect>
                        <Input
                            v-model="entry.label"
                            placeholder="名称"
                            class="h-8 w-[120px] shrink-0 text-[12px]"
                        />
                        <div class="relative flex-1">
                            <Input
                                v-model="entry.key"
                                :type="entry.show ? 'text' : 'password'"
                                placeholder="sk-..."
                                class="h-8 pr-8 text-[12px]"
                            />
                            <button
                                type="button"
                                :aria-label="entry.show ? '隐藏' : '显示'"
                                class="absolute inset-y-0 right-1.5 flex w-6 items-center justify-center text-[var(--hig-secondary-label)] outline-none hover:text-foreground"
                                @click="entry.show = !entry.show"
                            >
                                <EyeOff v-if="entry.show" :size="14" />
                                <Eye v-else :size="14" />
                            </button>
                        </div>
                        <Button
                            variant="ghost"
                            size="icon-sm"
                            aria-label="删除"
                            class="shrink-0 text-[var(--hig-secondary-label)] hover:text-destructive"
                            @click="removeKey(entry.id)"
                        >
                            <Trash2 :size="14" />
                        </Button>
                    </div>
                </div>
            </section>
        </div>
    </HigSheet>
</template>
