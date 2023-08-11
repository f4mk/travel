declare type Nullable<T> = T | null
declare type ValueOf<T> = T[keyof T]
declare type AnyFunction = (...args: any[]) => unknown
