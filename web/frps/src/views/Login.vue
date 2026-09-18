<template>
  <div class="login-container">
    <div class="login-lang-bar">
      <el-dropdown trigger="click" @command="handleLocaleChange">
        <span class="login-lang-trigger">
          {{ currentLocaleLabel }}
          <el-icon class="login-lang-arrow"><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="l in supportedLocales"
              :key="l.value"
              :command="l.value"
              :class="{ 'is-active': locale === l.value }"
            >
              {{ l.label }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <div class="login-card">
      <div class="login-header">
        <div class="login-logo">
          <LogoIcon class="logo-icon" />
        </div>
        <h1 class="login-title">{{ $t('login.title') }}</h1>
        <p class="login-subtitle">{{ $t('login.subtitle') }}</p>
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="user">
          <el-input
            v-model="form.user"
            :placeholder="$t('login.username')"
            :prefix-icon="User"
            size="large"
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="$t('login.password')"
            :prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="form.rememberMe">{{ $t('login.rememberMe') }}</el-checkbox>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            {{ $t('login.signIn') }}
          </el-button>
        </el-form-item>

        <div v-if="errorMsg" class="login-error">
          <el-alert :title="errorMsg" type="error" show-icon :closable="false" />
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { User, Lock, ArrowDown } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import LogoIcon from '../assets/icons/logo.svg?component'
import { useAuth } from '../composables/useAuth'
import { useLanguage } from '../composables/useLanguage'
import type { LocaleCode } from '../i18n'

const router = useRouter()
const { login } = useAuth()
const { t } = useI18n()
const { locale, setLocale, supportedLocales } = useLanguage()

const formRef = ref<FormInstance>()
const loading = ref(false)
const errorMsg = ref('')

const currentLocaleLabel = computed(() => {
  const found = supportedLocales.find((l) => l.value === locale.value)
  return found ? found.label : locale.value
})

function handleLocaleChange(code: LocaleCode) {
  setLocale(code)
}

const form = reactive({
  user: '',
  password: '',
  rememberMe: false,
})

const rules = computed<FormRules>(() => ({
  user: [{ required: true, message: t('login.usernameRequired'), trigger: 'blur' }],
  password: [{ required: true, message: t('login.passwordRequired'), trigger: 'blur' }],
}))

async function handleLogin() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  errorMsg.value = ''

  try {
    await login({
      user: form.user,
      password: form.password,
      rememberMe: form.rememberMe,
    })
    router.push('/')
  } catch (err: any) {
    errorMsg.value = err?.message || t('login.loginFailed')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-height: 100dvh;
  background: var(--content-bg);
  padding: 20px;
  position: relative;
}

.login-lang-bar {
  position: absolute;
  top: 16px;
  right: 20px;
}

.login-lang-trigger {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  padding: 6px 10px;
  border-radius: 6px;
  transition: all 0.15s ease;
  user-select: none;
}

.login-lang-trigger:hover {
  color: var(--text-primary);
  background: var(--hover-bg);
}

.login-lang-arrow {
  font-size: 12px;
}

.login-lang-bar :deep(.el-dropdown-menu__item.is-active) {
  color: var(--el-color-primary);
  font-weight: 600;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: var(--header-bg);
  border: 1px solid var(--header-border);
  border-radius: 16px;
  padding: 40px 32px 32px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.logo-icon {
  width: 48px;
  height: 48px;
}

.login-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
}

.login-form {
  margin-top: 8px;
}

.login-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.login-btn {
  width: 100%;
}

.login-error {
  margin-top: 8px;
}
</style>
