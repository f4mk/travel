import { HttpErrorWithPayload } from './errors'
import { ErrorObject } from './types'

export const createUrlOrThrow = <Fn extends (data: Req) => string, Req>(
  fn: Fn,
  data: Req | undefined
) => {
  if (!data) {
    throw new Error('query function does not provide data')
  }
  return fn(data)
}

export const handleSuccessWithPayload = (res: Response) => res.json()
export const handleErrorWithPayload = async (res: Response) => {
  const payload: ErrorObject = await res.json()
  throw new HttpErrorWithPayload(res, payload)
}
