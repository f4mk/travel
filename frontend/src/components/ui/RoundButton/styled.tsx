import { forwardRef } from 'react'
import styled from '@emotion/styled'
import { UnstyledButton as UnstyledButtonUI } from '@mantine/core'

import { ButtonProps } from './types'

const UnstyledButtonWrapper = forwardRef<
  HTMLButtonElement,
  React.ComponentProps<typeof UnstyledButtonUI>
>((props, ref) => <UnstyledButtonUI {...props} ref={ref} />)

export const UnstyledButton = styled(UnstyledButtonWrapper)<ButtonProps>`
  max-width: 40px;
  max-height: 40px;
  border-radius: 50%;
`
