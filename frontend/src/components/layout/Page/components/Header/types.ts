import { HttpError } from '#/api/request'

export type UseDataArgs = {
  onError: (e: HttpError) => void
}
