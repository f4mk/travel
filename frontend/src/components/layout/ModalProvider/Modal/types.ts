import { PropsWithChildren } from 'react'

export type Props = PropsWithChildren & {
  opened: boolean
  onClose: () => void
}
