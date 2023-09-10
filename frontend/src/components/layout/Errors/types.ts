import { ComponentProps, PropsWithChildren } from 'react'

import { ErrorBoundary } from './ErrorBoundary'

export type Props = PropsWithChildren & {
  fallback?: ComponentProps<typeof ErrorBoundary>['fallback']
  resetList?: string[]
}
