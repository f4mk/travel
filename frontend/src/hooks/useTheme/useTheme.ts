import { Theme } from '@emotion/react'
import { useMediaQuery } from '@mantine/hooks'

import { DARK_MODE, LIGHT_MODE } from './consts'

export const useTheme = (): Partial<Theme> => {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')

  const theme = prefersDarkMode ? DARK_MODE : LIGHT_MODE
  return {
    colorScheme: theme,
    globalStyles: (_) => ({
      '*, *::before, *::after': {
        boxSizing: 'border-box',
      },
      html: {
        height: '100%',
      },
      body: {
        height: '100%',
        fontFamily: 'Roboto',
      },
      '#root': {
        height: '100%',
      },
    }),
    defaultRadius: 'md',
  }
}
