import { getFreshToken } from '#/api/auth/utils/getFreshToken'

import { RequestArgs } from './types'
export const createRequest = <R, T = unknown>(args: RequestArgs<T>) => {
  return () => request<T, R>(args)
}
export const request = async <T, R>({
  url,
  method = 'GET',
  body
}: RequestArgs<T>) => {
  const token = getFreshToken()

  const headers = {
    Authorization: `Bearer ${token}`,
    'Content-Type': 'application/json'
  }

  const config = {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined
  }

  const response = await fetch(url, config)
  if (!response.ok) {
    throw new Error(`HTTP error: ${response.status}`)
  }
  return response.json() as R
}
