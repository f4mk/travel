import { components } from '#/types/external/auth'

export type RequestArgs<Res, Req> = {
  url: string | ((data: Req) => string)
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: Req
  lang: string
  headers?: Record<string, string>
  handleErrorCodes?: number[] | CustomDataHandler<Res>[]
  handleSuccessCodes?: number[] | CustomDataHandler<Res>[]
}

type CustomDataHandler<Res> = {
  code: number
  handler: (res: Response) => Promise<Res>
}

export type Handler = (_: Response) => Promise<unknown>
export type ErrorObject = components['schemas']['ErrorResponse']
export type Result<H> = H extends {
  [K in number]?: (_: Response) => Promise<infer R>
}
  ? R
  : never
export type Handlers = {
  [K in number]?: Handler
}
