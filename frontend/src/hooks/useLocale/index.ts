import { useEffect, useState } from 'react'

import { DEFAULT_LOCALE, LOCALES } from './consts'
import { loadTranslation } from './utils'

export const useLocale = (loc: string) => {
  const prefLocale = LOCALES.find((l) => l === loc) || DEFAULT_LOCALE
  const [locale, setLocale] = useState(prefLocale)
  const [translations, setTranslations] = useState<
    Record<string, Record<string, string>>
  >({})

  useEffect(() => {
    loadTranslation(loc, setTranslations)
  }, [loc])

  return {
    locale,
    setLocale,
    t: translations[locale]
  }
}
