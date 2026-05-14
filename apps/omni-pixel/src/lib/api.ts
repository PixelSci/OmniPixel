const BASE_URL = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'

interface RequestOptions {
    method: string
    path: string
    body?: unknown
}

async function request<T>({ method, path, body }: RequestOptions): Promise<T> {
    const url = `${BASE_URL}${path}`

    const res = await fetch(url, {
        method,
        headers: body ? { 'Content-Type': 'application/json' } : undefined,
        body: body ? JSON.stringify(body) : undefined,
        credentials: 'include',
    })

    if (!res.ok) {
        const error = await res.json().catch(() => ({}))
        throw new ApiError(res.status, error?.message ?? `Request failed: ${res.status}`)
    }

    return res.json() as Promise<T>
}

export class ApiError extends Error {
    status: number
    constructor(status: number, message: string) {
        super(message)
        this.name = 'ApiError'
        this.status = status
    }
}

export const api = {
    get: <T>(path: string) => request<T>({ method: 'GET', path }),
    post: <T>(path: string, body?: unknown) => request<T>({ method: 'POST', path, body }),
    delete: <T>(path: string) => request<T>({ method: 'DELETE', path }),
}
