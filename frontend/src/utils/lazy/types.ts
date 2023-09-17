import { ComponentType, LazyExoticComponent } from 'react'

export type Ns = {
  [K in string]: unknown
}

export type LazyNs<T extends Ns> = {
  [K in keyof T as T[K] extends ComponentType<any>
    ? K
    : never]: T[K] extends ComponentType<any>
    ? LazyExoticComponent<T[K]>
    : never
}
