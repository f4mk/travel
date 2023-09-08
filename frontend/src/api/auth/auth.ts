import { useQuery } from '@tanstack/react-query'

import { createRequest } from '#/api/request/request'

import { LoginRequest } from './types'

// TODO: add options
export const useLogin = ({ email, password }: LoginRequest) => {
  const url = '/api/auth/login'
  return useQuery(
    [url, email, password],
    createRequest({
      url,
      method: 'POST',
      body: { email, password }
    })
  )
}
