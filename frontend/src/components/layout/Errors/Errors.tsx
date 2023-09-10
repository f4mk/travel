import { memo } from 'react'
import { useQueryErrorResetBoundary } from '@tanstack/react-query'

import { HttpErrorWithPayload } from '#/api/request/errors'

import { ErrorBoundary } from './ErrorBoundary'
import { Fallback } from './Fallback'
import { Props } from './types'

export const Errors = memo(({ fallback, resetList, children }: Props) => {
  const { reset } = useQueryErrorResetBoundary()
  return (
    <ErrorBoundary
      onReset={reset}
      resetList={resetList || []}
      fallback={
        fallback ??
        (({ error }) => {
          if (error instanceof HttpErrorWithPayload) {
            return <Fallback message={error.payload.error} />
          }
          return <Fallback message={error.message} />
        })
      }
    >
      {children}
    </ErrorBoundary>
  )
})
