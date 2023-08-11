import { FC, useMemo } from 'react'
import { MantineProvider } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'

import { MainLayout } from '#/components/layout/MainLayout'
import { PageContextProvider } from '#/context/usePageContext'

import { AppProps } from './types'

import './styles.css'

const DARK_MODE = 'dark'
const LIGHT_MODE = 'light'

export const App: FC<AppProps> = ({ children, pageContext }) => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)', true)

  const theme = prefersDarkMode ? DARK_MODE : LIGHT_MODE

  const context = useMemo(() => {
    return pageContext
  }, [pageContext])

  return (
    <PageContextProvider pageContext={context}>
      <MantineProvider
        theme={{
          colorScheme: theme,
          globalStyles: (_) => ({
            '*, *::before, *::after': {
              boxSizing: 'border-box'
            },

            body: {
              fontFamily: 'Roboto'
            }
          })
        }}
        withGlobalStyles
        withNormalizeCSS
      >
        <MainLayout>{children}</MainLayout>
      </MantineProvider>
    </PageContextProvider>
  )
}
