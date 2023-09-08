import { ErrorObject } from './types'

export class HttpErrorWithPayload extends Error {
  payload: ErrorObject
  response: Response

  constructor(response: Response, payload: ErrorObject) {
    super(payload.error)
    Object.setPrototypeOf(this, HttpErrorWithPayload.prototype)
    this.name = 'HttpErrorWithPayload'
    this.payload = payload
    this.response = response
  }
}

export class HttpError extends Error {
  response: Response

  constructor(response: Response) {
    super(response.statusText)
    Object.setPrototypeOf(this, HttpError.prototype)
    this.name = 'HttpError'
    this.response = response
  }
}
