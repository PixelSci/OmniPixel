import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'

const BASE_URL = '/api/v1'

const instance: AxiosInstance = axios.create({
    baseURL: BASE_URL,
    timeout: 30000,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
})

instance.interceptors.request.use(
    (config) => {
        const raw = localStorage.getItem('omni-pixel:access-token')
        if (raw) {
            try {
                const token = JSON.parse(raw)
                if (token) config.headers.Authorization = `Bearer ${token}`
            } catch { /* ignore parse errors */ }
        }
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

export const http = {
    get: <T>(url: string, config?: AxiosRequestConfig) =>
        instance.get<T>(url, config).then((r) => r.data),
    post: <T>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
        instance.post<T>(url, data, config).then((r) => r.data),
    delete: <T>(url: string, config?: AxiosRequestConfig) =>
        instance.delete<T>(url, config).then((r) => r.data),
}

export default instance
