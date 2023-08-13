import { IntlProvider } from 'react-intl'
import { RouterProvider } from 'react-router-dom'
import { MantineProvider } from '@mantine/core'

import { useLocale } from '#/hooks'
import { useTheme } from '#/hooks'

import { router } from './router'

export const App = () => {
  const { locale, t } = useLocale(navigator.language)
  return (
    <MantineProvider theme={useTheme()} withGlobalStyles withNormalizeCSS>
      <IntlProvider locale={locale} messages={t}>
        <RouterProvider router={router} />
      </IntlProvider>
    </MantineProvider>
  )
}
