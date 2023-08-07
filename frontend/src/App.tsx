import { RouterProvider } from 'react-router-dom'
import { MantineProvider } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'

import { router } from './router'
const DARK_MODE = 'dark'
const LIGHT_MODE = 'light'

export const App = () => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')

  const theme = prefersDarkMode ? DARK_MODE : LIGHT_MODE

  return (
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
      <RouterProvider router={router} />
    </MantineProvider>
  )
}
