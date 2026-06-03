import { http } from './http'

export interface SigninRequest {
    email: string
    password: string
}

export interface SignupRequest {
    username: string
    email: string
    password: string
}

export interface AuthUser {
    id: string
    username: string
    email: string
    display_name?: string | null
    avatar_url?: string | null
    status: string
    created_at: string
    updated_at: string
}

export interface SigninResponse {
    access_token: string
    token_type: string
    expires_in: number
    user: AuthUser
}

export interface SignupResponse {
    access_token: string
    token_type: string
    expires_in: number
    user: AuthUser
}

export function signin(data: SigninRequest): Promise<SigninResponse> {
    return http.post<SigninResponse>('/signin', data)
}

export function signup(data: SignupRequest): Promise<SignupResponse> {
    return http.post<SignupResponse>('/signup', data)
}
