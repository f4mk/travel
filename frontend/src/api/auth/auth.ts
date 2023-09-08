import { useMutation, UseMutationOptions } from '@tanstack/react-query'

import { createRequest } from '#/api/request'
import { useGetLocale } from '#/hooks'

import { LoginError, LoginRequest, LoginResponse } from './types'

export const useLogin = (
  options?: UseMutationOptions<LoginResponse, LoginError, LoginRequest>
) => {
  const url = '/api/auth/login'
  // TODO: check if it returs locale in proper format
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      lang,
      method: 'POST'
    }),
    options
  )
}
