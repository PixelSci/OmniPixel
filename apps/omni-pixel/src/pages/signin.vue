<script setup lang="ts">
import { ArrowRight, Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import logo from '@/assets/logo.svg'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { isAuthenticated, loading, error, signin, signup, clearError } = useAuth()

// Redirect if already authenticated
if (isAuthenticated.value) {
    router.replace('/')
}

const isRegistering = ref(false)
const showPassword = ref(false)

const form = ref({
    username: '',
    email: '',
    password: '',
})

const canSubmit = computed(() => {
    const { email, password, username } = form.value
    if (!email.trim() || !password.trim()) return false
    if (isRegistering.value && !username.trim()) return false
    return !loading.value
})

function toggleMode() {
    isRegistering.value = !isRegistering.value
    clearError()
}

async function handleSubmit() {
    if (!canSubmit.value) return
    clearError()

    const success = isRegistering.value
        ? await signup({
            username: form.value.username.trim(),
            email: form.value.email.trim(),
            password: form.value.password,
        })
        : await signin({
            email: form.value.email.trim(),
            password: form.value.password,
        })

    if (success) {
        router.push('/')
    }
}
</script>

<template>
    <div class="flex min-h-screen items-center justify-center">
        <div
            class="hig-liquid-glass relative w-full max-w-[400px] overflow-hidden rounded-2xl border border-[rgba(0,0,0,0.08)] bg-white/70 px-8 py-10 backdrop-blur-[28px] backdrop-saturate-180 dark:border-[rgba(255,255,255,0.12)] dark:bg-[rgba(44,44,46,0.62)]"
        >
            <!-- logo -->
            <div class="mb-6 flex flex-col items-center gap-3">
                <img :src="logo" alt="OmniPixel" class="size-12">
                <div class="text-center">
                    <h1 class="text-[20px] font-semibold tracking-tight text-foreground">
                        OmniPixel
                    </h1>
                    <p class="mt-1 text-[13px] text-[var(--hig-secondary-label)]">
                        {{ isRegistering ? '创建你的账户' : '登录以继续' }}
                    </p>
                </div>
            </div>

            <!-- error banner -->
            <div
                v-if="error"
                class="mb-4 flex items-center gap-2 rounded-lg bg-[rgba(255,59,48,0.08)] px-3 py-2 text-[12px] text-[#ff3b30] dark:bg-[rgba(255,69,58,0.12)] dark:text-[#ff6961]"
            >
                <span>{{ error }}</span>
            </div>

            <!-- form -->
            <form class="flex flex-col gap-3" @submit.prevent="handleSubmit">
                <div v-if="isRegistering" class="flex flex-col gap-1.5">
                    <Label class="text-[12px] text-[var(--hig-secondary-label)]" for="username">
                        用户名
                    </Label>
                    <Input
                        id="username"
                        v-model="form.username"
                        type="text"
                        placeholder="你的用户名"
                        autocomplete="username"
                        :disabled="loading"
                        class="h-9 text-[13px]"
                    />
                </div>

                <div class="flex flex-col gap-1.5">
                    <Label class="text-[12px] text-[var(--hig-secondary-label)]" for="email">
                        邮箱
                    </Label>
                    <Input
                        id="email"
                        v-model="form.email"
                        type="email"
                        placeholder="name@example.com"
                        autocomplete="email"
                        :disabled="loading"
                        class="h-9 text-[13px]"
                    />
                </div>

                <div class="flex flex-col gap-1.5">
                    <Label class="text-[12px] text-[var(--hig-secondary-label)]" for="password">
                        密码
                    </Label>
                    <div class="relative">
                        <Input
                            id="password"
                            v-model="form.password"
                            :type="showPassword ? 'text' : 'password'"
                            placeholder="••••••••"
                            autocomplete="current-password"
                            :disabled="loading"
                            class="h-9 pr-9 text-[13px]"
                        />
                        <button
                            type="button"
                            :aria-label="showPassword ? '隐藏密码' : '显示密码'"
                            class="absolute inset-y-0 right-2 flex w-7 items-center justify-center text-[var(--hig-secondary-label)] outline-none hover:text-foreground"
                            @click="showPassword = !showPassword"
                        >
                            <EyeOff v-if="showPassword" :size="15" />
                            <Eye v-else :size="15" />
                        </button>
                    </div>
                </div>

                <Button
                    type="submit"
                    :disabled="!canSubmit"
                    class="mt-2 h-10 w-full gap-2 text-[14px] font-medium"
                >
                    <Loader2 v-if="loading" :size="16" class="animate-spin" />
                    <span>{{ isRegistering ? '创建账户' : '登录' }}</span>
                    <ArrowRight v-if="!loading" :size="16" />
                </Button>
            </form>

            <!-- toggle mode -->
            <p class="mt-5 text-center text-[12px] text-[var(--hig-secondary-label)]">
                {{ isRegistering ? '已有账户？' : '没有账户？' }}
                <button
                    type="button"
                    class="font-medium text-[#0077ED] hover:underline dark:text-[#3aa3ff]"
                    @click="toggleMode"
                >
                    {{ isRegistering ? '去登录' : '注册' }}
                </button>
            </p>
        </div>
    </div>
</template>
