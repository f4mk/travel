export type Params = string[]
export type ParamsMap<T extends string[]> = {
  [K in T[number]]: string
}
