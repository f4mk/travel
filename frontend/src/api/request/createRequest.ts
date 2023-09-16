import { request } from './request'
import { Handler, RequestArgs } from './types'
import {
  createUrlOrThrow,
  handleErrorWithPayload,
  handleSuccessWithPayload,
} from './utils'
export const createRequest =
  <Res, Req = void>({
    url,
    method = 'GET',
    lang,
    handleErrorCodes = [],
    handleSuccessCodes = [],
    headers = {},
    body,
  }: RequestArgs<Res, Req>) =>
  (data?: Req) => {
    const payload = body
    const uri = typeof url === 'function' ? createUrlOrThrow(url, data) : url
    const newReq = new Request(uri, {
      method,
      headers: { 'Accept-Language': lang, ...headers },
      body: payload && JSON.stringify(payload),
    })

    const handlers: Record<number, Handler | undefined> = {}
    if (handleSuccessCodes.length === 0) {
      method === 'POST'
        ? (handlers[201] = handleSuccessWithPayload)
        : (handlers[200] = handleSuccessWithPayload)
    } else {
      handleSuccessCodes.forEach((handle) => {
        typeof handle === 'number'
          ? (handlers[handle] = handleSuccessWithPayload)
          : (handlers[handle.code] = handle.handler)
      })
    }

    handleErrorCodes.forEach((handle) => {
      typeof handle === 'number'
        ? (handlers[handle] = handleErrorWithPayload)
        : (handlers[handle.code] = handle.handler)
    })

    return request(newReq, handlers) as Promise<Res>
  }
