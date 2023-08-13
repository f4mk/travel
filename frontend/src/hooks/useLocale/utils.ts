import { Dict } from './types'

export const loadTranslation = async (
  langCode: string,
  cb: (translations: any) => void
) => {
  try {
    switch (langCode) {
      default: {
        const translations = await import('#/translations/en.json')
        return cb((dict: Dict) => ({
          ...dict,
          [langCode]: translations.default
        }))
      }
    }
  } catch (error) {
    console.error(`Failed to load translations for ${langCode}:`, error)
    return null
  }
}
