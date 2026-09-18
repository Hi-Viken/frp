import { useI18n } from 'vue-i18n'
import { SUPPORTED_LOCALES, saveLocale } from '../i18n'
import type { LocaleCode } from '../i18n'

export function useLanguage() {
  const { locale } = useI18n()

  function setLocale(code: LocaleCode): void {
    locale.value = code
    saveLocale(code)
    document.documentElement.lang = code
  }

  function currentLocale(): LocaleCode {
    return locale.value as LocaleCode
  }

  return {
    locale,
    setLocale,
    currentLocale,
    supportedLocales: SUPPORTED_LOCALES,
  }
}
