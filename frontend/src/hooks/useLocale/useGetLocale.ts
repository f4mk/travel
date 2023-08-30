import { createContext, useContext } from 'react'

const context = createContext(navigator.language)
export const LocaleProvider = context.Provider
export const useGetLocale = () => {
  return useContext(context)
}
