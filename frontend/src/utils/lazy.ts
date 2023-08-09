import {
  type ComponentType,
  lazy as reactLazy,
  type LazyExoticComponent
} from 'react'

type Ns = {
  [K in string]: unknown
}

type LazyNs<T extends Ns> = {
  [K in keyof T as T[K] extends ComponentType<any>
    ? K
    : never]: T[K] extends ComponentType<any>
    ? LazyExoticComponent<T[K]>
    : never
}

export const lazy = <T extends Ns>(fn: () => Promise<T>): LazyNs<T> => {
  let promise: Promise<T> | undefined
  return new Proxy(
    {},
    {
      get(_, key) {
        return reactLazy(async () => {
          promise = promise || fn()
          return {
            default: Reflect.get(await promise, key) as any
          }
        })
      }
    }
  ) as LazyNs<T>
}
