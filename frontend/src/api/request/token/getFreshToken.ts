import { HttpError } from '#/api/request/errors'

import { getTokenFromHeader } from './getTokenFromHeader'
import { isTokenValid } from './isTokenValid'
import { getToken, setToken } from './tokenUtils'

export const getFreshToken = async () => {
  const token = getToken()
  if (!token || !isTokenValid(token)) {
    const url = '/api/auth/refresh'
    const defaultLang = 'en-US'
    const newReq = new Request(url, {
      method: 'POST',
      body: JSON.stringify({}),
      headers: { 'Accept-Language': defaultLang }
    })
    const res = await fetch(newReq)
    if (!res.ok) {
      throw new HttpError(res)
    }
    const authHeader = res.headers.get('Authorization')
    const newToken = getTokenFromHeader(authHeader)
    setToken(newToken)
    return newToken
  }

  return token
}
