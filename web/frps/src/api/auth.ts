import { http } from './http'

export interface LoginRequest {
  user: string
  password: string
  rememberMe?: boolean
}

export interface LoginResponse {
  token: string
  expiresAt: number
}

export interface LoginStatusResponse {
  authRequired: boolean
  loggedIn: boolean
  user?: string
}

export const login = (req: LoginRequest) => {
  return http.post<LoginResponse>('../api/login', req)
}

export const logout = () => {
  return http.post<{ Code: number; Msg: string }>('../api/logout')
}

export const getLoginStatus = () => {
  return http.get<LoginStatusResponse>('../api/login/status')
}
