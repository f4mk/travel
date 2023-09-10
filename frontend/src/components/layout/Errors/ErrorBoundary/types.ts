import { ComponentType, PropsWithChildren } from 'react'

type FallbackProps = {
  error: Error
  reset?: () => void
}
export type Props = PropsWithChildren & {
  fallback: ComponentType<FallbackProps>
  resetList: string[]
  onReset?: () => void
}
export type State = {
  error?: Error
  resetList: string[]
}
