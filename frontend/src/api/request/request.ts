import { getFreshToken } from '#/api/auth/utils/getFreshToken'

import { HttpError } from './errors'
import { Handler, Handlers, Result } from './types'

export const request = async (request: Request, handlers: Handlers) => {
  const token = getFreshToken()
  request.headers.set('Authorization', `Bearer ${token}`)

  const response = await fetch(request)

  const handler = handlers[response.status]
  if (handler) {
    return handler(response) as Promise<Result<Handler>>
  }

  if (!response.ok) {
    throw new HttpError(response)
  }
  throw TypeError(`fetch: handler for code ${response.status} is not specified`)
}
