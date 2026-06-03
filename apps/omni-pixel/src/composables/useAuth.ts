import { useLocalStorage } from '@vueuse/core'
import { computed, ref } from 'vue'
import {
    signin as signinApi,
    signup as signupApi,
    type AuthUser,
    type SigninRequest,
    type SignupRequest,
} from '@/lib/auth'
import { ApiError } from '@/lib/http'

const TOKEN_KEY = 'omni-pixel:access-token'
const USER_KEY = 'omni-pixel:user'

// ── Shared reactive state (singleton) ────────────────────────────────────

const token = useLocalStorage<string | null>(TOKEN_KEY, null)
const user = useLocalStorage<AuthUser | null>(USER_KEY, null)

const loading = ref(false)
const error = ref<string | null>(null)

// Authorization header is set in http.ts interceptor reading from localStorage

// ── Composable ───────────────────────────────────────────────────────────

export function useAuth() {
    const isAuthenticated = computed(() => !!token.value)
    const currentUser = computed(() => user.value)

    async function signin(request: SigninRequest): Promise<boolean> {
        loading.value = true
        error.value = null
        try {
            const res = await signinApi(request)
            token.value = res.access_token
            user.value = res.user
            return true
        } catch (e) {
            error.value = e instanceof ApiError
                ? (e.status === 401 ? '邮箱或密码错误' : e.message)
                : (e instanceof Error ? e.message : '登录失败')
            return false
        } finally {
            loading.value = false
        }
    }

    async function signup(request: SignupRequest): Promise<boolean> {
        loading.value = true
        error.value = null
        try {
            const res = await signupApi(request)
            token.value = res.access_token
            user.value = res.user
            return true
        } catch (e) {
            error.value = e instanceof ApiError
                ? (e.status === 409 ? '该邮箱已被注册' : e.message)
                : (e instanceof Error ? e.message : '注册失败')
            return false
        } finally {
            loading.value = false
        }
    }

    function logout() {
        token.value = null
        user.value = null
        error.value = null
    }

    function clearError() {
        error.value = null
    }

    return {
        token,
        user: currentUser,
        isAuthenticated,
        loading,
        error,
        signin,
        signup,
        logout,
        clearError,
    }
}
