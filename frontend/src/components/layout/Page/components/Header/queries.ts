import { useMemo } from 'react'

import { useGetMe } from '#/api/user'

import { UseDataArgs } from './types'

export const useData = ({ onError }: UseDataArgs) => {
  const { data } = useGetMe({ onError, suspense: true })

  return useMemo(() => {
    const { name, email } = data || {}

    return {
      name,
      email,
    }
  }, [data])
}
