import { HelmetProvider } from 'react-helmet-async'
import { IntlProvider } from 'react-intl'
import { RouterProvider } from 'react-router-dom'
import { MantineProvider } from '@mantine/core'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

import { ModalProvider } from '#/components/layout/ModalProvider'
import { LocaleProvider, useLocale, useTheme } from '#/hooks'

import { createRouter } from './router'

const queryClient = new QueryClient()

export const App = () => {
  const { locale, t } = useLocale(navigator.language)
  const theme = useTheme()
  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <MantineProvider theme={theme} withGlobalStyles withNormalizeCSS>
          <IntlProvider locale={locale} messages={t}>
            <LocaleProvider value={locale}>
              {/* TODO: find a better way to provide Modal as a child to router
                yet keep it declared in <App/>
              */}
              <RouterProvider router={createRouter(ModalProvider)} />
            </LocaleProvider>
          </IntlProvider>
        </MantineProvider>
      </QueryClientProvider>
    </HelmetProvider>
  )
}
