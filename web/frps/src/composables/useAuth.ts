import { ref, readonly } from 'vue'
import { getLoginStatus, login as apiLogin, logout as apiLogout } from '../api/auth'
import type { LoginRequest } from '../api/auth'

const isLoggedIn = ref(false)
const authRequired = ref(false)
const currentUser = ref('')
const initialized = ref(false)

let initPromise: Promise<void> | null = null

export function useAuth() {
  async function checkAuth(): Promise<void> {
    try {
      const status = await getLoginStatus()
      authRequired.value = status.authRequired
      isLoggedIn.value = status.loggedIn
      currentUser.value = status.user || ''
    } catch {
      authRequired.value = true
      isLoggedIn.value = false
    }
    initialized.value = true
  }

  function ensureInitialized(): Promise<void> {
    if (initialized.value) return Promise.resolve()
    if (!initPromise) {
      initPromise = checkAuth().finally(() => {
        initPromise = null
      })
    }
    return initPromise
  }

  async function login(req: LoginRequest): Promise<void> {
    await apiLogin(req)
    isLoggedIn.value = true
    currentUser.value = req.user
  }

  async function logout(): Promise<void> {
    try {
      await apiLogout()
    } finally {
      isLoggedIn.value = false
      currentUser.value = ''
    }
  }

  function markUnauthorized(): void {
    isLoggedIn.value = false
    currentUser.value = ''
  }

  return {
    isLoggedIn: readonly(isLoggedIn),
    authRequired: readonly(authRequired),
    currentUser: readonly(currentUser),
    initialized: readonly(initialized),
    checkAuth,
    ensureInitialized,
    login,
    logout,
    markUnauthorized,
  }
}
