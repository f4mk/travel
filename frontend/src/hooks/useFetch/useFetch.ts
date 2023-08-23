// TODO:

export type Params<F extends (...args: any[]) => unknown> = F extends (
  ...args: infer A
) => unknown
  ? A
  : never

const fetcher = (...args: Params<typeof fetch>) =>
  fetch(...args).then((res) => res.json())

export const useFetch = (url: string) => {
  return fetcher(url)
}
