import {
  useMutation,
  UseMutationOptions,
  useQuery,
  UseQueryOptions,
} from '@tanstack/react-query'

import { createRequest, HttpError } from '#/api/request'
import { useGetLocale } from '#/hooks'

import {
  CreateUserError,
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
} from './types'

export const useCreateUser = (
  options?: UseMutationOptions<
    CreateUserResponse,
    CreateUserError,
    CreateUserRequest
  >
) => {
  const url = '/api/user'
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

export const useDeleteUser = (
  options?: UseMutationOptions<
    DeleteUserResponse,
    DeleteUserError,
    DeleteUserRequest
  >
) => {
  const url = '/api/user'
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
  const url = '/api/user'
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
  const url = `/api/user/${userId}`
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
  const url = '/api/user/me'
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
