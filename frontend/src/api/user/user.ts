import {
  useMutation,
  UseMutationOptions,
  useQuery,
  UseQueryOptions,
} from '@tanstack/react-query'

import { createRequest, HttpError } from '#/api/request'
import { useGetLocale } from '#/hooks'

import {
  CreateUserRequest,
  CreateUserResponse,
  DeleteUserError,
  DeleteUserRequest,
  DeleteUserResponse,
  GetMeResponse,
  GetUserRequest,
  GetUserResponse,
  UpdateUserError,
  UpdateUserRequest,
  UpdateUserResponse,
  VerifyUserError,
  VerifyUserRequest,
  VerifyUserResponse,
} from './types'

export const useCreateUser = (
  options?: UseMutationOptions<CreateUserResponse, HttpError, CreateUserRequest>
) => {
  const url = '/api/users'
  const lang = useGetLocale()
  return useMutation(async (body: CreateUserRequest) => {
    const newReq = new Request(url, {
      method: 'POST',
      body: JSON.stringify(body),
      headers: { 'Accept-Language': lang },
    })
    const res = await fetch(newReq)
    if (!res.ok) {
      throw new HttpError(res)
    }

    return res.json() as Promise<CreateUserResponse>
  }, options)
}

export const useDeleteUser = (
  options?: UseMutationOptions<
    DeleteUserResponse,
    DeleteUserError,
    DeleteUserRequest
  >
) => {
  const url = '/api/users'
  const lang = useGetLocale()
  return useMutation(
    createRequest({
      url,
      method: 'DELETE',
      lang,
      handleErrorCodes: [400, 401, 403, 404, 500],
    }),
    options
  )
}

export const useUpdateUser = (
  options?: UseMutationOptions<
    UpdateUserResponse,
    UpdateUserError,
    UpdateUserRequest
  >
) => {
  const url = '/api/users'
  const lang = useGetLocale()
  return useMutation(
    createRequest({
      url,
      method: 'PUT',
      lang,
      handleErrorCodes: [400, 401, 403, 404, 500],
    }),
    options
  )
}

export const useGetUser = (
  userId: GetUserRequest,
  options?: UseQueryOptions<GetUserResponse>
) => {
  const url = `/api/users/${userId}`
  const lang = useGetLocale()
  return useQuery({
    queryKey: [url, lang],
    queryFn: createRequest({
      url,
      lang,
      handleErrorCodes: [401, 404, 500],
    }),
    ...options,
  })
}

export const useGetMe = (
  options?: UseQueryOptions<GetMeResponse, HttpError>
) => {
  const url = '/api/users/me'
  const lang = useGetLocale()
  return useQuery({
    queryKey: [url, lang],
    queryFn: createRequest({
      url,
      lang,
      handleErrorCodes: [401, 404, 500],
    }),
    ...options,
  })
}

export const useVerifyUser = (
  options?: UseMutationOptions<
    VerifyUserResponse,
    VerifyUserError,
    VerifyUserRequest
  >
) => {
  const url = '/api/users/verify'
  const lang = useGetLocale()
  return useMutation(
    createRequest({
      url,
      method: 'POST',
      lang,
      handleErrorCodes: [400, 409, 500],
    }),
    options
  )
}
