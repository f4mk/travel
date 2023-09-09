import { useMutation, UseMutationOptions } from '@tanstack/react-query'

import { getTokenFromHeader, HttpError, setToken } from '#/api/request'
import { useGetLocale } from '#/hooks'

import { LoginError, LoginRequest, LoginResponse } from './types'

export const useLogin = (
  options?: UseMutationOptions<LoginResponse, LoginError, LoginRequest>
) => {
  const url = '/api/auth/login'
  // TODO: check if this provides a proper locale value
  const lang = useGetLocale()

  return useMutation(async (body: LoginRequest) => {
    const newReq = new Request(url, {
      method: 'POST',
      body: JSON.stringify(body),
      headers: { 'Accept-Language': lang }
    })
    const res = await fetch(newReq)
    if (!res.ok) {
      throw new HttpError(res)
    }
    const authHeader = res.headers.get('Authorization')
    setToken(getTokenFromHeader(authHeader))

    return res.json() as Promise<LoginResponse>
  }, options)
}
