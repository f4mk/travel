import { HelmetProvider } from 'react-helmet-async'
import { IntlProvider } from 'react-intl'
import { RouterProvider } from 'react-router-dom'
import { MantineProvider } from '@mantine/core'

import { ModalProvider } from '#/components/layout/ModalProvider'
import { LocaleProvider, useLocale, useTheme } from '#/hooks'

import { router } from './router'

export const App = () => {
  const { locale, t } = useLocale(navigator.language)

  return (
    <HelmetProvider>
      <MantineProvider theme={useTheme()} withGlobalStyles withNormalizeCSS>
        <IntlProvider locale={locale} messages={t}>
          <LocaleProvider value={locale}>
            <ModalProvider>
              <RouterProvider router={router} />
            </ModalProvider>
          </LocaleProvider>
        </IntlProvider>
      </MantineProvider>
    </HelmetProvider>
  )
}
