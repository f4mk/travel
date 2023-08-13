export const loadTranslation = async (
  langCode: string,
  cb: (translations: any) => void
) => {
  try {
    let translations
    switch (langCode) {
      default:
        translations = await import('#/translations/en.json')
        return cb(translations)
    }
    return translations.default
  } catch (error) {
    console.error(`Failed to load translations for ${langCode}:`, error)
    return null
  }
}
