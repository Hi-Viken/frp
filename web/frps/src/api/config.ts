import type { FrpsConfig, ConfigAPIResponse } from '../types/config'

const BASE = '../api/config'

async function request<T>(url: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(url, {
    credentials: 'include',
    ...options,
  })
  if (response.status === 401) {
    window.location.hash = '#/login'
    throw new Error('Unauthorized')
  }
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }
  return response.json()
}

export const configApi = {
  get: () => request<ConfigAPIResponse>(BASE),

  save: (cfg: FrpsConfig) =>
    request<ConfigAPIResponse>(BASE, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(cfg),
    }),

  reload: () =>
    request<ConfigAPIResponse>(`${BASE}/reload`, { method: 'POST' }),

  restart: () =>
    request<ConfigAPIResponse>(`${BASE}/restart`, { method: 'POST' }),
}
