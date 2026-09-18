import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

export const SUPPORTED_LOCALES = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
] as const

export type LocaleCode = (typeof SUPPORTED_LOCALES)[number]['value']

export const DEFAULT_LOCALE: LocaleCode = 'en'

const STORAGE_KEY = 'frp-locale'

export function loadLocale(): LocaleCode {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && SUPPORTED_LOCALES.some((l) => l.value === stored)) {
    return stored as LocaleCode
  }
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) return 'zh'
  return DEFAULT_LOCALE
}

export function saveLocale(locale: LocaleCode): void {
  localStorage.setItem(STORAGE_KEY, locale)
}

const i18n = createI18n({
  legacy: false,
  locale: loadLocale(),
  fallbackLocale: DEFAULT_LOCALE,
  messages: { en, zh },
})

export default i18n
