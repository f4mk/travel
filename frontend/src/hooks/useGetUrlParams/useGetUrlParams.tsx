import { Params, ParamsMap } from './types'

export const useGetUrlParams = (...params: Params) => {
  const query = window.location.search
  const urlParams = new URLSearchParams(query)

  return params.reduce((acc, param) => {
    const value = urlParams.get(param)
    if (value) {
      acc[param] = value
    }
    return acc
  }, {} as ParamsMap<Params>)
}
