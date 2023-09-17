import { lazy as reactLazy } from 'react'

import { LazyNs, Ns } from './types'

export const lazy = <T extends Ns>(fn: () => Promise<T>): LazyNs<T> => {
  let promise: Promise<T> | undefined
  return new Proxy(
    {},
    {
      get(_, key) {
        return reactLazy(async () => {
          promise = promise || fn()
          return {
            default: Reflect.get(await promise, key) as any,
          }
        })
      },
    }
  ) as LazyNs<T>
}
