import { PropsWithChildren } from 'react'
import { MantineNumberSize } from '@mantine/core'

export type Props = PropsWithChildren & {
  opened: boolean
  onClose: () => void
  size?: MantineNumberSize
}
