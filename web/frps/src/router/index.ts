import { createRouter, createWebHashHistory } from 'vue-router'
import ServerOverview from '../views/ServerOverview.vue'
import Clients from '../views/Clients.vue'
import ClientDetail from '../views/ClientDetail.vue'
import Proxies from '../views/Proxies.vue'
import ProxyDetail from '../views/ProxyDetail.vue'
import ConfigManager from '../views/ConfigManager.vue'
import Login from '../views/Login.vue'
import { useAuth } from '../composables/useAuth'

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior() {
    return { top: 0 }
  },
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      name: 'ServerOverview',
      component: ServerOverview,
      meta: { requiresAuth: true },
    },
    {
      path: '/clients',
      name: 'Clients',
      component: Clients,
      meta: { requiresAuth: true },
    },
    {
      path: '/clients/:key',
      name: 'ClientDetail',
      component: ClientDetail,
      meta: { requiresAuth: true },
    },
    {
      path: '/proxies/:type?',
      name: 'Proxies',
      component: Proxies,
      meta: { requiresAuth: true },
    },
    {
      path: '/proxy/:name',
      name: 'ProxyDetail',
      component: ProxyDetail,
      meta: { requiresAuth: true },
    },
    {
      path: '/config',
      name: 'ConfigManager',
      component: ConfigManager,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to) => {
  if (to.meta.requiresAuth === false) return true

  const { ensureInitialized, authRequired, isLoggedIn } = useAuth()
  await ensureInitialized()

  if (!authRequired.value) return true

  if (!isLoggedIn.value) {
    return { name: 'Login' }
  }
  return true
})

export default router
