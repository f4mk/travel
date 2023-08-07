import { ComponentType, ReactNode } from 'react'

export type Props = {
  component: () => Promise<{ default: ComponentType<any> }>
  fallback: ReactNode
}
