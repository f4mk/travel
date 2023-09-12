import { useMutation, UseMutationOptions } from '@tanstack/react-query'

import {
  createRequest,
  getTokenFromHeader,
  HttpError,
  setToken,
} from '#/api/request'
import { useGetLocale } from '#/hooks'

import {
  LoginError,
  LoginRequest,
  LoginResponse,
  LogoutAllError,
  LogoutAllRequest,
  LogoutAllResponse,
  LogoutError,
  LogoutRequest,
  LogoutResponse,
  PasswordChangeError,
  PasswordChangeRequest,
  PasswordChangeResponse,
  PasswordResetError,
  PasswordResetRequest,
  PasswordResetResponse,
  PasswordResetSubmitError,
  PasswordResetSubmitRequest,
  PasswordResetSubmitResponse,
} from './types'

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
      headers: { 'Accept-Language': lang },
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

export const useLogout = (
  options?: UseMutationOptions<LogoutResponse, LogoutError, LogoutRequest>
) => {
  const url = '/api/auth/logout'
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 401, 404, 500],
    }),
    options
  )
}

export const useLogoutAll = (
  options?: UseMutationOptions<
    LogoutAllResponse,
    LogoutAllError,
    LogoutAllRequest
  >
) => {
  const url = '/api/auth/logout/all'
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 401, 404, 500],
    }),
    options
  )
}

export const usePasswordReset = (
  options?: UseMutationOptions<
    PasswordResetResponse,
    PasswordResetError,
    PasswordResetRequest
  >
) => {
  const url = '/api/auth/password/reset'
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 500],
    }),
    options
  )
}

export const usePasswordResetSubmit = (
  options?: UseMutationOptions<
    PasswordResetSubmitResponse,
    PasswordResetSubmitError,
    PasswordResetSubmitRequest
  >
) => {
  const url = '/api/auth/password/reset/submit'
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 403, 500],
    }),
    options
  )
}

export const usePasswordChange = (
  options?: UseMutationOptions<
    PasswordChangeResponse,
    PasswordChangeError,
    PasswordChangeRequest
  >
) => {
  const url = '/api/auth/password/change'
  const lang = useGetLocale()

  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 401, 404, 500],
    }),
    options
  )
}
