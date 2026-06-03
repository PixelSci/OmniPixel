import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'

const BASE_URL = '/api/v1'
const TOKEN_KEY = 'omni-pixel:access-token'

const instance: AxiosInstance = axios.create({
    baseURL: BASE_URL,
    timeout: 30000,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
})

export function getToken(): string | null {
    const raw = localStorage.getItem(TOKEN_KEY)
    if (!raw || raw === 'null') return null
    return raw
}

instance.interceptors.request.use(
    (config) => {
        if ((config as any).anonymous) return config
        const token = getToken()
        if (token) config.headers.Authorization = `Bearer ${token}`
        return config
    },
    (error) => Promise.reject(error),
)

instance.interceptors.response.use(
    (response: AxiosResponse) => response,
    (error) => {
        if (axios.isCancel(error) || error.code === 'ERR_CANCELED') {
            return Promise.reject(error)
        }
        const message = error.response?.data?.message ?? error.message ?? 'Request failed'
        return Promise.reject(new ApiError(error.response?.status ?? 0, message))
    },
)

export class ApiError extends Error {
    status: number
    constructor(status: number, message: string) {
        super(message)
        this.name = 'ApiError'
        this.status = status
    }
}

export interface RequestOptions extends AxiosRequestConfig {
    anonymous?: boolean
}

export const http = {
    get: <T>(url: string, config?: RequestOptions) =>
        instance.get<T>(url, config).then((r) => r.data),
    post: <T>(url: string, data?: unknown, config?: RequestOptions) =>
        instance.post<T>(url, data, config).then((r) => r.data),
    delete: <T>(url: string, config?: RequestOptions) =>
        instance.delete<T>(url, config).then((r) => r.data),
}

export default instance
